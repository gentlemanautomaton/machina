package qvm

import (
	"github.com/gentlemanautomaton/machina/qemu"
	"github.com/gentlemanautomaton/machina/qemu/qdev"
	"github.com/gentlemanautomaton/machina/qemu/qguest"
	"github.com/gentlemanautomaton/machina/qemu/qhost"
)

// Definition describes the configuration of a QEMU virtual machine.
type Definition struct {
	Settings  qguest.Settings
	Resources qhost.Resources
	Topology  qdev.Topology
}

// Options returns a set of QEMU configuration options for the QEMU
// virtual machine definition.
func (def *Definition) Options() qemu.Options {
	var opts qemu.Options

	// Guest Configuration
	opts = append(opts, def.Settings.Options()...)

	// Host Configuration
	opts = append(opts, def.Resources.Options()...)

	// Devices
	opts = append(opts, def.Topology.Options()...)

	return opts
}
