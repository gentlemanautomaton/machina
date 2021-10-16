package main

import (
	"context"
	"errors"

	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/mdevfs"
)

// TeardownCmd removes host resources that were previously prepared for a
// virtual machine.
type TeardownCmd struct {
	Machines []string `kong:"arg,predictor=machines,help='Virtual machines to teardown.'"`
}

// Run executes the machine teardown command.
func (cmd TeardownCmd) Run(ctx context.Context) error {
	machines, sys, err := LoadAndComposeMachines(cmd.Machines...)
	if err != nil {
		return err
	}

	for i := range machines {
		teardownMachine(machines[i].MachineInfo, machines[i].Definition, sys)
	}

	return nil
}

func teardownMachine(info machina.MachineInfo, definition machina.Definition, sys machina.System) error {
	// TODO: Clean up as many devices and connections as possible, even if one hits an error
	for _, device := range definition.Devices {
		if err := teardownDevice(device, sys); err != nil {
			return err
		}
	}
	for _, conn := range definition.Connections {
		if err := teardownConnection(conn, sys); err != nil {
			return err
		}
	}
	return nil
}

func teardownDevice(device machina.Device, sys machina.System) error {
	// If a device ID has not been provided in the machina configuration
	// we don't have any way of inspecting its condition.
	if device.ID.IsZero() {
		return nil
	}

	// Check whether mediated devices are supported
	if !mdevfs.Supported() {
		return errors.New("mediated devices are not supported on the local system")
	}

	// Check whether there is a mediated device with the device ID
	dev := mdevfs.NewMediatedDevice(device.ID)
	exists, err := dev.Exists()
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	// If a mediated device exists, remove it
	if err := dev.Remove(); err != nil {
		return err
	}

	return nil
}

func teardownConnection(conn machina.Connection, sys machina.System) error {
	return nil
}
