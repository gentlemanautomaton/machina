package main

import (
	"context"
	"fmt"

	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/qemu/qvm"
	"github.com/gentlemanautomaton/machina/qemugen"
)

// ArgsCmd displays the QEMU arguments used to invocake one or more virtual machines.
type ArgsCmd struct {
	Machines []string `kong:"arg,predictor=machines,help='Virtual machines to show QEMU arguments for.'"`
}

// Run executes the machine config generation command.
//
// FIXME: Perform additional validation before generating the config
func (cmd ArgsCmd) Run(ctx context.Context) error {
	sys, err := LoadSystem()
	if err != nil {
		return fmt.Errorf("failed to load system configuration: %v", err)
	}

	var machines []machina.MachineInfo
	var vms []qvm.Definition
	for _, name := range cmd.Machines {
		machine, err := LoadMachine(name)
		if err != nil {
			return fmt.Errorf("failed to load machine configuration for \"%s\": %v", name, err)
		}
		vm, err := qemugen.Build(machine, sys)
		if err != nil {
			return fmt.Errorf("failed to build configuration for \"%s\": %v", name, err)
		}
		machines = append(machines, machine.Info())
		vms = append(vms, vm)
	}

	for i := range vms {
		if i > 0 {
			fmt.Println()
		}
		if len(vms) > 1 {
			fmt.Printf("----%s----\n", machines[i].Name)
		}
		fmt.Println(vms[i].Options().String())
	}

	return nil
}
