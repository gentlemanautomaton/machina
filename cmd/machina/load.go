package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/systemdgen"
)

// ComposedMachine holds virtual machine information and that has been loaded
// and coalesced with the local system configuration.
type ComposedMachine struct {
	machina.MachineInfo
	machina.Definition
}

// EnumMachines attempts to load the set of machina machine names that are
// present on the local system.
func EnumMachines() (names []machina.MachineName, err error) {
	root := os.DirFS(MachineDir())
	matches, err := fs.Glob(root, "*.conf.json")
	if err != nil {
		return nil, err
	}
	for _, name := range matches {
		names = append(names, machina.MachineName(strings.TrimSuffix(name, ".conf.json")))
	}
	return names, nil
}

// LoadMachineConnections loads a set of connections for the given names.
//
// If any of the names are for machines, all of the connections for those
// machines will be returned.
func LoadMachineConnections(names ...string) (conns []machina.MachineConnection, sys machina.System, err error) {
	sys, err = LoadSystem()
	if err != nil {
		return nil, machina.System{}, fmt.Errorf("failed to load system configuration: %v", err)
	}

	definitions := make(map[machina.MachineName]machina.Definition)
	seen := make(map[string]bool)
	for _, name := range names {
		// Parse the name into machine and optional connection name parts
		machineName, connName := parseMachineConnection(name)

		// If we haven't loaded the definition for this machine yet, load it
		// now
		definition, loaded := definitions[machineName]
		if !loaded {
			machine, err := LoadMachine(machineName)
			if err != nil {
				return nil, sys, fmt.Errorf("failed to load machine configuration for \"%s\": %v", name, err)
			}
			definition, err = machina.Build(machine, sys)
			if err != nil {
				return nil, sys, fmt.Errorf("failed to build configuration for \"%s\": %v", name, err)
			}
			definitions[machineName] = definition
		}

		// If a connection name hasn't been provided, add all of the machine's connections
		if connName == "" {
			for _, conn := range definition.Connections {
				link := machina.MakeLinkName(machineName, conn)
				if seen[link] {
					continue // Already added
				}
				conns = append(conns, machina.MachineConnection{
					Machine:    machineName,
					Connection: conn,
				})
				seen[link] = true
			}
			continue
		}

		// Fine the given connetion name in the list
		found := false
		for _, conn := range definition.Connections {
			if conn.Name != connName {
				continue
			}
			link := machina.MakeLinkName(machineName, conn)
			if seen[link] {
				continue // Already added
			}
			conns = append(conns, machina.MachineConnection{
				Machine:    machineName,
				Connection: conn,
			})
			seen[link] = true
			found = true
			break
		}

		if !found {
			return nil, sys, fmt.Errorf("failed to locate connection \"%s\" within \"%s\"", connName, machineName)
		}
	}

	return conns, sys, nil
}

func parseMachineConnection(name string) (machine machina.MachineName, conn machina.ConnectionName) {
	parts := strings.SplitN(name, ".", 2)
	switch len(parts) {
	case 0:
		return "", ""
	case 1:
		return machina.MachineName(parts[0]), ""
	default:
		return machina.MachineName(parts[0]), machina.ConnectionName(parts[1])
	}
}

// LoadAndComposeMachines attempts to load the machina definitions for the given
// machine names. It also returns the system configuration.
func LoadAndComposeMachines(names ...machina.MachineName) (vms []ComposedMachine, sys machina.System, err error) {
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

// LoadMachineUnits attempts to load the machina definitions for the given
// machine names and return the unit names for each of them. An error
// is returned if the definition for any machine could not be loaded.
//
// FIXME: Ensure that the returned unit names are valid.
func LoadMachineUnits(names ...machina.MachineName) (units []string, err error) {
	machines, err := LoadMachines(names...)
	if err != nil {
		return nil, err
	}

	for i := range machines {
		units = append(units, systemdgen.UnitNameForQEMU(machines[i].Name))
	}

	return units, nil
}

// LoadMachines attempts to load the machina definitions for the given
// machine names. The machines are returned as-is, without application
// of tags or other system-wide configuration.
func LoadMachines(names ...machina.MachineName) (machines []machina.Machine, err error) {
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
func LoadMachine(name machina.MachineName) (m machina.Machine, err error) {
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
