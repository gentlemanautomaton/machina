package main

import (
	"context"
)

// DisableCmd disables systemd units for one or more virtual machines.
type DisableCmd struct {
	Machines []string `kong:"arg,optional,help='Virtual machines to disable.'"`
}

// Run executes the machine disablement command.
func (cmd DisableCmd) Run(ctx context.Context) error {
	units, err := LoadMachineUnits(cmd.Machines...)
	if err != nil {
		return err
	}
	return systemctl(ctx, "disable", units)
}
