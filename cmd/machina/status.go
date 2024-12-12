package main

import (
	"context"

	"github.com/gentlemanautomaton/machina"
)

// StatusCmd prints the systemd unit status for one or more virtual machines.
type StatusCmd struct {
	Machines []machina.MachineName `kong:"arg,predictor=machines,help='Virtual machines to report the status of.'"`
}

// Run executes the machine status command.
func (cmd StatusCmd) Run(ctx context.Context) error {
	units, err := LoadMachineUnits(cmd.Machines...)
	if err != nil {
		return err
	}
	return systemctl(ctx, "status", units)

	// TODO: Consider direct dbus calls instead of shelling out to systemctl
	//
	// https://github.com/systemd/systemd/blob/dc131951b5f903b698f624a0234560d7a822ff21/src/systemctl/systemctl-show.c#L298
}
