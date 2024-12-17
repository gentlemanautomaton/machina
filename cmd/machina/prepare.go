package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/filesystem/mdevfs"
)

// PrepareCmd prepares the host environment for a machine.
type PrepareCmd struct {
	QEMU PrepareQEMUCmd `kong:"cmd,help='Prepares the host environment for a qemu/kvm process to start.'"`
}

// PrepareQEMUCmd prepares the host environment for a machine's qemu/kvm
// process to run.
type PrepareQEMUCmd struct {
	Machines []machina.MachineName `kong:"arg,predictor=machines,help='Virtual machines to prepare.'"`
}

// Run executes the machine preparation command.
func (cmd PrepareQEMUCmd) Run(ctx context.Context) error {
	machines, sys, err := LoadAndComposeMachines(cmd.Machines...)
	if err != nil {
		return err
	}

	for i := range machines {
		if err := prepareMachine(machines[i].MachineInfo, machines[i].Definition, sys); err != nil {
			return err
		}
	}

	return nil
}

func prepareMachine(info machina.MachineInfo, definition machina.Definition, sys machina.System) error {
	for _, device := range definition.Devices {
		if err := prepareDevice(device, sys); err != nil {
			return err
		}
	}
	for _, conn := range definition.Connections {
		if err := prepareConnection(conn, sys); err != nil {
			return err
		}
	}
	return nil
}

func prepareDevice(device machina.Device, sys machina.System) error {
	// Mediated devices require a device identifier
	if device.ID.IsZero() {
		return nil
	}

	// Look through the mediated devices in machina's system
	// configuration for one that supplies the desired device class
	mdevs := sys.MediatedDevices.WithClass(device.Class)
	if len(mdevs) == 0 {
		return fmt.Errorf("failed to locate device class %s for device %s", device.Class, device.Name)
	}

	// Check whether mediated devices are supported
	if !mdevfs.Supported() {
		return errors.New("mediated devices are not supported on the local system")
	}

	// Check whether the mediated device already exists
	{
		device := mdevfs.NewMediatedDevice(device.ID)
		if exists, err := device.Exists(); err != nil {
			return err
		} else if exists {
			return nil
		}
	}

	// Loop through all of the mediated devices that were matched
	for _, mdev := range mdevs {
		// Translate the device class to the supported type name
		// supplied by machina's system configuration
		tname := mdev.Types[device.Class]
		if tname == "" {
			continue
		}

		// Prepare access to the physical device through sysfs
		parent := mdevfs.NewPhysicalDevice(mdev.Address)
		if exists, err := parent.Exists(); err != nil {
			return fmt.Errorf("failed to open mediated device file system for %s: %v", mdev.Address, err)
		} else if !exists {
			continue
		}

		// Enumerate the supported types in sysfs
		types, err := parent.Types()
		if err != nil {
			return fmt.Errorf("failed to enumrate supported types for mediated device %s: %v", mdev.Address, err)
		}

		// Search the enumerated types for one with the requestd type
		// name
		typ, found := types.FindName(tname)
		if !found {
			continue
		}

		// Create the mediated device
		return typ.Create(device.ID)
	}

	return nil
}

func prepareConnection(conn machina.Connection, sys machina.System) error {
	return nil
}
