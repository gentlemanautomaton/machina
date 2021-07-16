package main

import (
	"context"
	"fmt"

	"github.com/gentlemanautomaton/machina"
)

// ConnectCmd enables the connections to one or more virtual machines.
type ConnectCmd struct {
	Machines []string `kong:"arg,optional,help='Virtual machines to bridge.'"`
}

// Run executes the connect command.
func (cmd ConnectCmd) Run(ctx context.Context) error {
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
			if err := enableConnection(vm.Name, conn, sys); err != nil {
				if firstError == nil {
					firstError = err
				}
				fmt.Printf("%s: Failed: %v\n", machina.LinkName(vm.Name, conn), err)
			} else {
				fmt.Printf("%s: Enabled\n", machina.LinkName(vm.Name, conn))
			}
		}
	}

	return firstError
}
