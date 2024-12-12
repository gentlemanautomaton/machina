package main

import (
	"context"

	"github.com/gentlemanautomaton/machina"
)

// EnableCmd enables systemd units for one or more virtual machines.
type EnableCmd struct {
	Machines []machina.MachineName `kong:"arg,predictor=machines,help='Virtual machines to enable.'"`
}

// Run executes the machine enablement command.
func (cmd EnableCmd) Run(ctx context.Context) error {
	units, err := LoadMachineUnits(cmd.Machines...)
	if err != nil {
		return err
	}
	return systemctl(ctx, "enable", units)
}
