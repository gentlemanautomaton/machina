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
	"github.com/gentlemanautomaton/machina/swtpm"
	"github.com/gentlemanautomaton/machina/swtpmgen"
	"github.com/gentlemanautomaton/machina/systemdgen"
	"github.com/gentlemanautomaton/systemdconf"
)

// GenerateCmd generates a QEMU invocation for one or more virtual machines.
type GenerateCmd struct {
	Preview  bool                  `kong:"optional,help='Print the generated systemd configurtion but do no apply it.'"`
	Machines []machina.MachineName `kong:"arg,predictor=machines,help='Specify virtual machines to generate.'"`
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
	var tpms []swtpm.Settings
	for _, name := range cmd.Machines {
		machine, err := LoadMachine(name)
		if err != nil {
			return fmt.Errorf("failed to load machine configuration for \"%s\": %v", name, err)
		}
		vm, err := qemugen.Build(machine, sys)
		if err != nil {
			return fmt.Errorf("failed to build QEMU configuration for \"%s\": %v", name, err)
		}
		tpm, err := swtpmgen.BuildSettings(machine, sys)
		if err != nil {
			return fmt.Errorf("failed to build software TPM configuration for \"%s\": %v", name, err)
		}
		machines = append(machines, machine.Info())
		vms = append(vms, vm)
		tpms = append(tpms, tpm)
	}

	var qemuUnits []string
	var tpmUnits []string
	for i := range vms {
		var (
			qemuOptions        = vms[i].Options()
			tpmEmulatorOptions = tpms[i].Emulator.Options()
			tpmSetupOptions    = tpms[i].Setup.Options()
		)

		var qemuBindToUnits []string
		if tpms[i].Emulator.Enabled {
			qemuBindToUnits = append(qemuBindToUnits, systemdgen.UnitNameForTPM(machines[i].Name)+".service")
		}

		var qemuBuf bytes.Buffer
		if _, err := systemdconf.WriteSections(&qemuBuf, systemdgen.BuildQEMU(machines[i], qemuOptions, qemuBindToUnits...)...); err != nil {
			return fmt.Errorf("failed to prepare QEMU configuration for %s: %v", machines[i].Name, err)
		}
		qemuUnits = append(qemuUnits, qemuBuf.String())

		var tpmBuf bytes.Buffer
		if _, err := systemdconf.WriteSections(&tpmBuf, systemdgen.BuildTPM(machines[i], tpmEmulatorOptions, tpmSetupOptions)...); err != nil {
			return fmt.Errorf("failed to prepare software TPM configuration for %s: %v", machines[i].Name, err)
		}
		tpmUnits = append(tpmUnits, tpmBuf.String())
	}

	for i := range vms {
		qemuUnitName := systemdgen.UnitNameForQEMU(machines[i].Name)
		tpmUnitName := systemdgen.UnitNameForTPM(machines[i].Name)

		if cmd.Preview {
			if i > 0 {
				fmt.Println()
			}
			fmt.Printf("----%s----\n", qemuUnitName)
			fmt.Print(qemuUnits[i])
			if tpms[i].Emulator.Enabled {
				fmt.Println()
				fmt.Printf("----%s----\n", tpmUnitName)
				fmt.Print(tpmUnits[i])
			}
			continue
		}

		if err := writeUnit(qemuUnitName, qemuUnits[i]); err != nil {
			return err
		}

		if tpms[i].Emulator.Enabled {
			if err := writeUnit(tpmUnitName, tpmUnits[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

func writeUnit(unit, config string) error {
	// Check for the presence of the systemd directory on the local machine
	if fi, err := os.Stat(machina.LinuxUnitDir); err != nil || !fi.IsDir() {
		return errors.New("generation is only supported on systems that store systemd units in /etc/systemd/system")
	}

	unitPath := filepath.Join(machina.LinuxUnitDir, unit+".service")

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
