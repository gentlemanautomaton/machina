package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/qemu/qvm"
	"github.com/gentlemanautomaton/machina/qemugen"
	"github.com/gentlemanautomaton/machina/systemdgen"
	"github.com/gentlemanautomaton/systemdconf"
)

// GenerateCmd generates a QEMU invocation for one or more virtual machines.
type GenerateCmd struct {
	Preview  bool     `kong:"optional,help='Print the generated systemd configurtion but do no apply it.'"`
	Machines []string `kong:"arg,help='Specify virtual machines to generate.'"`
}

// Run executes the machine config generation command.
//
// FIXME: Perform additional validation before generating the config
func (cmd GenerateCmd) Run(ctx context.Context) error {
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

	var units []string
	for i := range vms {
		options := vms[i].Options()

		var buf bytes.Buffer
		if _, err := systemdconf.WriteSections(&buf, systemdgen.Build(machines[i], options)...); err != nil {
			return fmt.Errorf("failed to prepare configuration for %s: %v", machines[i].Name, err)
		}
		units = append(units, buf.String())
	}

	for i := range vms {
		unit := fmt.Sprintf("machina-%s", machines[i].Name)

		if cmd.Preview {
			if i > 0 {
				fmt.Println()
			}
			fmt.Printf("----%s----\n", machines[i].Name)
			fmt.Print(units[i])
			continue
		}

		if err := writeUnit(unit, units[i]); err != nil {
			return err
		}
	}

	return nil
}

func writeUnit(unit, config string) error {
	// Check for the presence of the systemd directory on the local machine
	if fi, err := os.Stat(linuxUnitDir); err != nil || !fi.IsDir() {
		return errors.New("generation is only supported on systems that store systemd units in /etc/systemd/system")
	}

	unitPath := filepath.Join(linuxUnitDir, unit+".service")

	file, err := os.Create(unitPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, strings.NewReader(config)); err != nil {
		return err
	}

	return nil
}
