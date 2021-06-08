package qguest

import "github.com/gentlemanautomaton/machina/qemu"

// Settings describes all of the non-device properties of a QEMU guest.
type Settings struct {
	Identity  Identity
	Processor Processor
	Memory    Memory
	Globals   qemu.Globals
}

// Options returns a set of QEMU virtual machine options for implementing
// the specification.
func (s Settings) Options() qemu.Options {
	var opts qemu.Options

	opts = append(opts, s.Identity.Options()...)
	opts.Add("enable-kvm")
	opts.Add("nodefaults")
	opts.Add("nographic")
	opts.Add("machine", qemu.Parameters{
		{Name: "type", Value: "q35"},
		//{Name: "accel", Value: "kvm"},
		{Name: "vmport", Value: "off"},
	}...)
	opts = append(opts, s.Processor.Options()...)
	opts = append(opts, s.Memory.Options()...)
	opts = append(opts, s.Globals.Options()...)

	return opts
}
