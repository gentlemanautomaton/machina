package main

import (
	"os"

	"github.com/gentlemanautomaton/machina"
)

// ConfDir returns the path where machina configuration should be installed
// on the local system.
func ConfDir() string {
	if fi, err := os.Stat(machina.LinuxConfDir); err != nil || !fi.IsDir() {
		return "."
	}
	return machina.LinuxConfDir
}

// MachineDir returns the path where machina machine configuration should be
// installed on the local system.
func MachineDir() string {
	if fi, err := os.Stat(machina.LinuxMachineDir); err != nil || !fi.IsDir() {
		return "machine.conf.d"
	}
	return machina.LinuxMachineDir
}
