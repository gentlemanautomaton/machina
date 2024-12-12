package main

import (
	"context"

	"github.com/gentlemanautomaton/machina"
)

// DisableCmd disables systemd units for one or more virtual machines.
type DisableCmd struct {
	Machines []machina.MachineName `kong:"arg,predictor=machines,help='Virtual machines to disable.'"`
}

// Run executes the machine disablement command.
func (cmd DisableCmd) Run(ctx context.Context) error {
	units, err := LoadMachineUnits(cmd.Machines...)
	if err != nil {
		return err
	}
	return systemctl(ctx, "disable", units)
}
