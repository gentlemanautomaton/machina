package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/gentlemanautomaton/machina/qmp/qmpcmd"
)

// ShutdownCmd sends a shutdown command to the given virtual machines.
type ShutdownCmd struct {
	Machines []string `kong:"arg,predictor=machines,help='Virtual machines to shutdown gracefully.'"`
	System   bool     `kong:"system,help='Use QMP sockets reserved for systemd.'"`
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
			var sockets []string
			if cmd.System {
				sockets = attrs.SystemSocketPaths(vms[i].MachineInfo)
			} else {
				sockets = attrs.CommandSocketPaths(vms[i].MachineInfo)
			}
			if !attrs.Enabled || len(sockets) == 0 {
				fmt.Printf("Cannot shutdown %s: no QMP socket available\n", name)
				return
			}

			client, err := connectToQMP(sockets)
			if err != nil {
				fmt.Printf("Failed to shutdown %s: %v\n", name, err)
				return
			}
			defer client.Close()

			if err := client.Execute(ctx, qmpcmd.SystemPowerdown); err != nil {
				fmt.Printf("Failed to shutdown %s: %v\n", name, err)
			}
			fmt.Printf("Issued shutdown to %s\n", name)

			// FIXME: Wait for it to exit? But we can't wait for the unit
			// to become inactive, because the unit itself will call this
			// function.
		}(i)
	}

	wg.Wait()

	return nil
}
