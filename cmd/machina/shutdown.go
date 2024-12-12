package main

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/qmp"
	"github.com/gentlemanautomaton/machina/qmp/qmpcmd"
)

// ShutdownCmd sends a shutdown command to the given virtual machines.
type ShutdownCmd struct {
	Machines []machina.MachineName `kong:"arg,predictor=machines,help='Virtual machines to shutdown gracefully.'"`
	System   bool                  `kong:"system,help='Use QMP sockets reserved for systemd.'"`
	Wait     bool                  `kong:"wait,help='Wait for the virtual machines to stop before returning.'"`
	Timeout  time.Duration         `kong:"timeout,help='Requests forceful termination if the timeout is exceeded before a virtual machine is stopped. Implies wait.'"`
}

// Run executes the graceful shutdown command.
func (cmd ShutdownCmd) Run(ctx context.Context) error {
	vms, _, err := LoadAndComposeMachines(cmd.Machines...)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(vms))

	for i := range vms {
		go func(i int) {
			defer wg.Done()

			name := vms[i].Name
			attrs := vms[i].Attributes.QMP

			// Collect the set of potential sockets we can use.
			var sockets []string
			if cmd.System {
				sockets = attrs.SystemSocketPaths(vms[i].MachineInfo)
			} else {
				sockets = attrs.CommandSocketPaths(vms[i].MachineInfo)
			}
			if !attrs.Enabled || len(sockets) == 0 {
				fmt.Printf("%s: Failed to issue shutdown command: no QMP socket available.\n", name)
				return
			}

			// Attempt to open a socket for communication with the virtual
			// machine via QMP.
			client, err := connectToQMP(sockets)
			if err != nil {
				fmt.Printf("%s: Failed to issue shutdown command: %v\n", name, err)
				return
			}
			defer client.Close()

			// If we're going to be listening for events, prepare the event
			// listener now so that we can receive notifications for the
			// commands we're about to issue.
			var listener *qmp.Listener
			if cmd.Wait || cmd.Timeout > 0 {
				listener = client.Listen()
			}

			// Send a graceful shutdown command.
			if err := client.Execute(ctx, qmpcmd.SystemPowerdown); err != nil {
				fmt.Printf("%s: Failed to issue shutdown command: %v\n", name, err)
				return
			}

			// Report success.
			fmt.Printf("%s: Issued shutdown command.\n", name)

			// Exit if a blocking operation has not been requested.
			if !cmd.Wait && cmd.Timeout == 0 {
				return
			}

			// If a timeout has been provided, derive a new context with
			// that timeout. We'll use this while receiving events to stop
			// listening at the right time.
			listenerCtx := ctx
			if cmd.Timeout > 0 {
				var cancel context.CancelFunc
				listenerCtx, cancel = context.WithTimeout(ctx, cmd.Timeout)
				defer cancel()
			}

			// Send repeated shutdown requests every five seconds.
			// This may be needed to convince some guests to shut down.
			const powerdownInterval = time.Second * 5
			var stopIssuingPowerdowns func()
			{
				// Prepare a context that we'll use to stop a ticker later on.
				tickerCtx, cancel := context.WithCancel(ctx)

				// Prepare a channel that will close when our ticker stops.
				done := make(chan struct{})

				// The stopIssuingPowerdowns function can be used to stop the
				// ticker. It's safe to call it more than once.
				stopIssuingPowerdowns = func() {
					cancel()
					<-done
				}

				// Always stop the ticker when we're done handling this VM.
				defer stopIssuingPowerdowns()

				// Start a goroutine that will send commands every 5 seconds
				// until tickerCtx is cancelled.
				go func() {
					ticker := time.NewTicker(powerdownInterval)
					defer close(done)
					defer ticker.Stop()

					for {
						select {
						case <-ticker.C:
							client.Execute(ctx, qmpcmd.SystemPowerdown)
						case <-tickerCtx.Done():
							return
						}
					}
				}()
			}

			// Process events while we wait for something to happen.
			for {
				event, err := listener.Receive(listenerCtx)

				switch err {
				case io.EOF:
					// The QMP socket closed, which indicates that shutdown
					// is complete.

					// Stop sending system powerdown commands.
					stopIssuingPowerdowns()

					// Report completion.
					fmt.Printf("%s: QMP socket closed. Shutdown complete.\n", name)
					return
				case context.DeadlineExceeded, context.Canceled:
					if cmd.Timeout == 0 || err == ctx.Err() {
						// Our timeout didn't expire, but we received a
						// cancellation request via the command's ctx.
						return
					}

					// Our timeout expired. Stop sending system powerdown
					// commands and send a quit command isntead.
					stopIssuingPowerdowns()

					fmt.Printf("%s: Timeout expired. Sending QUIT message.\n", name)
					if err := client.Execute(ctx, qmpcmd.Quit); err != nil {
						fmt.Printf("%s: Failed to send QUIT: %v\n", name, err)
					}

					// Switch to wait mode.
					listenerCtx = ctx
				case nil:
					if data := event.Data.Bytes(); len(data) > 0 {
						fmt.Printf("%s: %s: %s\n", name, event.Event, string(data))
					} else {
						fmt.Printf("%s: %s\n", name, event.Event)
					}
				default:
					fmt.Printf("%s: %v\n", name, err)
				}
			}
		}(i)
	}

	wg.Wait()

	return nil
}
