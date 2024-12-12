package main

import (
	"context"

	"github.com/gentlemanautomaton/machina"
)

// StartCmd attempts to start the systemd units for the given virtual
// machines.
type StartCmd struct {
	Machines []machina.MachineName `kong:"arg,predictor=machines,help='Virtual machines to start.'"`
}

// Run executes the machine start command.
func (cmd StartCmd) Run(ctx context.Context) error {
	units, err := LoadMachineUnits(cmd.Machines...)
	if err != nil {
		return err
	}
	return systemctl(ctx, "start", units)
}
