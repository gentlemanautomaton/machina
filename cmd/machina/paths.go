package main

import (
	"os"
)

const (
	linuxBinDir            = "/usr/bin"
	linuxConfDir           = "/etc/machina"
	linuxMachineDir        = "/etc/machina/machine.conf.d"
	linuxUnitDir           = "/etc/systemd/system"
	linuxBashCompletionDir = "/usr/share/bash-completion/completions"
)

// ConfDir returns the path where machina configuration should be installed
// on the local system.
func ConfDir() string {
	if fi, err := os.Stat(linuxConfDir); err != nil || !fi.IsDir() {
		return "."
	}
	return linuxConfDir
}

// MachineDir returns the path where machina machine configuration should be
// installed on the local system.
func MachineDir() string {
	if fi, err := os.Stat(linuxMachineDir); err != nil || !fi.IsDir() {
		return "machine.conf.d"
	}
	return linuxMachineDir
}
