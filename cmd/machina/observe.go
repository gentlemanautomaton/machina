package main

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/gentlemanautomaton/machina"
)

// ObserveCmd listens for events for the given virtual machines.
type ObserveCmd struct {
	Machines []machina.MachineName `kong:"arg,predictor=machines,help='Virtual machines to observe.'"`
}

// Run executes the observation command.
func (cmd ObserveCmd) Run(ctx context.Context) error {
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
			sockets := attrs.CommandSocketPaths(vms[i].MachineInfo)
			if !attrs.Enabled || len(sockets) == 0 {
				fmt.Printf("Cannot observe %s: no QMP socket available\n", name)
				return
			}

			client, err := connectToQMP(sockets)
			if err != nil {
				fmt.Printf("Failed to observe %s: %v\n", name, err)
				return
			}
			defer client.Close()

			version, caps := client.ServerInfo()
			fmt.Printf("Connected to %s (%s, capabilities: %s)\n", name, version, caps)
			listener := client.Listen()

			for {
				event, err := listener.Receive(ctx)
				if err != nil {
					if err == context.DeadlineExceeded || err == context.Canceled || err == io.EOF {
						fmt.Printf("%s: QMP socket closed.\n", name)
						return
					}
					fmt.Printf("%s: %v\n", name, err)
					return
				}
				if data := event.Data.Bytes(); len(data) > 0 {
					fmt.Printf("%s: %s: %s\n", name, event.Event, string(data))
				} else {
					fmt.Printf("%s: %s\n", name, event.Event)
				}
			}
		}(i)
	}

	wg.Wait()

	return nil
}
