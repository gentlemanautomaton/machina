package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gentlemanautomaton/machina"
)

// ComposedMachine holds virtual machine information and that has been loaded
// and coalesced with the local system configuration.
type ComposedMachine struct {
	machina.MachineInfo
	machina.Definition
}

// EnumMachines attempts to load the set of machina machine names that are
// present on the local system.
func EnumMachines() (names []string, err error) {
	root := os.DirFS(MachineDir())
	matches, err := fs.Glob(root, "*.conf.json")
	if err != nil {
		return nil, err
	}
	for _, name := range matches {
		names = append(names, strings.TrimSuffix(name, ".conf.json"))
	}
	return names, nil
}

// LoadAndComposeMachines attempts to load the machina definitions for the given
// machine names. It also returns the system configuration.
func LoadAndComposeMachines(names ...string) (vms []ComposedMachine, sys machina.System, err error) {
	sys, err = LoadSystem()
	if err != nil {
		return nil, machina.System{}, fmt.Errorf("failed to load system configuration: %v", err)
	}

	for _, name := range names {
		machine, err := LoadMachine(name)
		if err != nil {
			return nil, sys, fmt.Errorf("failed to load machine configuration for \"%s\": %v", name, err)
		}
		definition, err := machina.Build(machine, sys)
		if err != nil {
			return nil, sys, fmt.Errorf("failed to build configuration for \"%s\": %v", name, err)
		}
		vms = append(vms, ComposedMachine{
			MachineInfo: machine.Info(),
			Definition:  definition,
		})
	}

	return vms, sys, nil
}

// LoadMachines attempts to load the machina definitions for the given
// machine names. The machines are returned as-is, without application
// of tags or other system-wide configuration.
func LoadMachines(names ...string) (machines []machina.Machine, err error) {
	for _, name := range names {
		machine, err := LoadMachine(name)
		if err != nil {
			return nil, fmt.Errorf("failed to load machine configuration for \"%s\": %v", name, err)
		}
		machines = append(machines, machine)
	}

	return machines, nil
}

// LoadMachine attempts to load the machine configuration for the given
// machine name.
func LoadMachine(name string) (m machina.Machine, err error) {
	path := filepath.Join(MachineDir(), fmt.Sprintf("%s.conf.json", name))
	f, err := os.Open(path)
	if err != nil {
		return machina.Machine{}, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	if err := dec.Decode(&m); err != nil {
		return machina.Machine{}, err
	}

	if m.Name == "" {
		m.Name = machina.MachineName(strings.TrimSuffix(filepath.Base(f.Name()), ".conf.json"))
	}

	return m, nil
}

// LoadSystem attempts to load the system configuration from a
// "system.conf.json" file.
func LoadSystem() (sys machina.System, err error) {
	f, err := os.Open(filepath.Join(ConfDir(), "machina.conf.json"))
	if err != nil {
		return machina.System{}, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	if err := dec.Decode(&sys); err != nil {
		return machina.System{}, err
	}

	return sys, nil
}
