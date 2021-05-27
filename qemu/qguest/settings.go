package qguest

import "github.com/gentlemanautomaton/machina/qemu"

// Settings describes all of the non-device properties of a QEMU guest.
type Settings struct {
	Identity  Identity
	Processor Processor
	Memory    Memory
}

// Options returns a set of QEMU virtual machine options for implementing
// the specification.
func (s Settings) Options() qemu.Options {
	var opts qemu.Options

	opts = append(opts, s.Identity.Options()...)
	opts.Add("enable-kvm")
	opts.Add("machine", qemu.Parameters{
		{Name: "type", Value: "q35"},
		{Name: "accel", Value: "kvm"},
	}...)
	opts = append(opts, s.Processor.Options()...)
	opts = append(opts, s.Memory.Options()...)

	opts.Add("nodefaults")
	opts.Add("nographic")

	return opts
}
