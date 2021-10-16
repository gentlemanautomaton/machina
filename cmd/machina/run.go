package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/qemugen"
)

// RunCmd starts a QEMU virtual machine.
type RunCmd struct {
	Machine string `kong:"arg,predictor=machines,optional,help='Virtual machine to run.'"`
}

// Run executes the run command.
func (cmd RunCmd) Run(ctx context.Context) (err error) {
	sys, err := LoadSystem()
	if err != nil {
		return fmt.Errorf("failed to load system configuration: %v", err)
	}

	machine, err := LoadMachine(cmd.Machine)
	if err != nil {
		return fmt.Errorf("failed to load machine configuration for \"%s\": %v", cmd.Machine, err)
	}

	composed, err := machina.Build(machine, sys)
	if err != nil {
		return fmt.Errorf("failed to build configuration for \"%s\": %v", cmd.Machine, err)
	}

	vm, err := qemugen.Build(machine, sys)
	if err != nil {
		return fmt.Errorf("failed to build configuration for \"%s\": %v", cmd.Machine, err)
	}

	if err := prepareMachine(machine.Info(), composed, sys); err != nil {
		return err
	}

	defer func(err *error) {
		for _, device := range composed.Devices {
			if devErr := teardownDevice(device, sys); devErr != nil {
				if *err != nil {
					*err = devErr
				}
			}
		}
	}(&err)

	args := vm.Options().Args()
	kvm := exec.CommandContext(ctx, "qemu-system-x86_64", args...)
	kvm.Stdout = os.Stdout
	kvm.Stderr = os.Stderr

	if err := kvm.Start(); err != nil {
		return fmt.Errorf("failed to start QEMU: %v", err)
	}

	return kvm.Wait()
}
