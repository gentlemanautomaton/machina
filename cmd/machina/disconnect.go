package main

import (
	"context"
	"fmt"

	"github.com/gentlemanautomaton/machina"
)

// DisconnectCmd enables the connections to one or more virtual machines.
type DisconnectCmd struct {
	Machines []string `kong:"arg,help='Virtual machines to bridge.'"`
}

// Run executes the disconnect command.
func (cmd DisconnectCmd) Run(ctx context.Context) error {
	vms, sys, err := LoadAndComposeMachines(cmd.Machines...)
	if err != nil {
		return err
	}

	var firstError error
	for i, vm := range vms {
		if i > 0 {
			fmt.Println()
		}

		for _, conn := range vm.Connections {
			if err := disableConnection(vm.Name, conn, sys); err != nil {
				if firstError == nil {
					firstError = err
				}
				fmt.Printf("%s: Failed: %v\n", machina.MakeLinkName(vm.Name, conn), err)
			} else {
				fmt.Printf("%s: Disabled\n", machina.MakeLinkName(vm.Name, conn))
			}
		}
	}

	return firstError
}
