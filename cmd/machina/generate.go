package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/gentlemanautomaton/machina/qemu/qvm"
	"github.com/gentlemanautomaton/machina/qemugen"
)

// GenerateCmd generates a QEMU invocation for one or more virtual machines.
type GenerateCmd struct {
	Inline   bool     `kong:"optional,help='Place generated arguments on a single line.'"`
	Machines []string `kong:"arg,optional,help='Virtual machines to generate.'"`
}

// Run executes the machine config generation command.
func (cmd GenerateCmd) Run(ctx context.Context) error {
	sys, err := LoadSystem()
	if err != nil {
		return fmt.Errorf("failed to load system configuration: %v", err)
	}

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
		vms = append(vms, vm)
	}

	for i, vm := range vms {
		if i > 0 {
			fmt.Println()
		}
		options := vm.Options()
		/*
			if err != nil {
				return fmt.Errorf("failed to prepare configuration options for \"%s\": %v", vm.Name, err)
			}
		*/
		if cmd.Inline {
			fmt.Printf("%s\n", strings.Join(options.Args(), " "))
		} else {
			for i, option := range options {
				last := i == len(options)-1
				if last {
					fmt.Printf("%s\n", option)
				} else {
					fmt.Printf("%s \\\n", option)
				}
			}
		}
	}

	return nil
}
