package main

import (
	"context"
	"fmt"

	"github.com/gentlemanautomaton/machina"
)

// ConnectCmd enables the connections to one or more virtual machines.
type ConnectCmd struct {
	MachinesOrConnections []string `kong:"arg,help='Machines or individual connections to bridge. Use [machine] or [machine].[conn].'"`
}

// Run executes the connect command.
func (cmd ConnectCmd) Run(ctx context.Context) error {
	mconns, sys, err := LoadMachineConnections(cmd.MachinesOrConnections...)
	if err != nil {
		return err
	}

	var firstError error
	for _, mconn := range mconns {
		link := machina.MakeLinkName(mconn.Machine, mconn.Connection)
		if err := enableConnection(mconn.Machine, mconn.Connection, sys); err != nil {
			if firstError == nil {
				firstError = err
			}
			fmt.Printf("%s: failed: %v\n", link, err)
		} else {
			fmt.Printf("%s: enabled\n", link)
		}
	}
	return firstError
}
