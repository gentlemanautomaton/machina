package main

import (
	"context"
	"fmt"

	"github.com/gentlemanautomaton/machina"
)

// DisconnectCmd enables the connections to one or more virtual machines.
type DisconnectCmd struct {
	MachinesOrConnections []string `kong:"arg,help='Machines or individual connections to remove from the bridge. Use [machine] or [machine].[conn]..'"`
}

// Run executes the disconnect command.
func (cmd DisconnectCmd) Run(ctx context.Context) error {
	return disconnect(cmd.MachinesOrConnections)
}

func disconnect(names []string) error {
	mconns, sys, err := LoadMachineConnections(names...)
	if err != nil {
		return err
	}

	var firstError error
	for _, mconn := range mconns {
		link := machina.MakeLinkName(mconn.Machine, mconn.Connection)
		if err := disableConnection(mconn.Machine, mconn.Connection, sys); err != nil {
			if firstError == nil {
				firstError = err
			}
			fmt.Printf("%s: failed: %v\n", link, err)
		} else {
			fmt.Printf("%s: disabled\n", link)
		}
	}

	return firstError
}
