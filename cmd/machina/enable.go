package main

import (
	"context"
)

// EnableCmd enables systemd units for one or more virtual machines.
type EnableCmd struct {
	Machines []string `kong:"arg,help='Virtual machines to enable.'"`
}

// Run executes the machine enablement command.
func (cmd EnableCmd) Run(ctx context.Context) error {
	units, err := LoadMachineUnits(cmd.Machines...)
	if err != nil {
		return err
	}
	return systemctl(ctx, "enable", units)
}
