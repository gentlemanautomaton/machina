package main

import (
	"context"
)

// StatusCmd prints the systemd unit status for one or more virtual machines.
type StatusCmd struct {
	Machines []string `kong:"arg,help='Virtual machines to enable.'"`
}

// Run executes the machine status command.
func (cmd StatusCmd) Run(ctx context.Context) error {
	units, err := LoadMachineUnits(cmd.Machines...)
	if err != nil {
		return err
	}
	return systemctl(ctx, "status", units)
}
