package main

import "context"

// StopCmd attempts to stop the systemd units for the given virtual
// machines.
//
// TODO: Perform a more graceful shutdown via QMP.
type StopCmd struct {
	Machines []string `kong:"arg,predictor=machines,help='Virtual machines to stop.'"`
}

// Run executes the machine stop command.
func (cmd StopCmd) Run(ctx context.Context) error {
	units, err := LoadMachineUnits(cmd.Machines...)
	if err != nil {
		return err
	}
	return systemctl(ctx, "stop", units)
}
