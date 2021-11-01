package qguest

import "github.com/gentlemanautomaton/machina/qemu"

// Settings describes all of the non-device properties of a QEMU guest.
type Settings struct {
	Identity  Identity
	Firmware  Firmware
	Clock     Clock
	Processor Processor
	Memory    Memory
	QMP       QMP
	Spice     Spice
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
	{
		params := qemu.Parameters{
			{Name: "type", Value: "q35"},
			//{Name: "accel", Value: "kvm"},
		}
		params = append(params, s.Spice.MachineParameters()...)
		params = append(params, s.Firmware.MachineParameters()...)
		opts.Add("machine", params...)
	}
	opts = append(opts, s.Processor.Options()...)
	opts = append(opts, s.Memory.Options()...)
	opts = append(opts, s.Clock.Options()...)
	opts = append(opts, s.QMP.Options()...)
	opts = append(opts, s.Spice.Options()...)
	opts = append(opts, s.Globals.Options()...)
	opts.Add("boot", qemu.Parameters{
		{Name: "menu", Value: "on"},
		{Name: "reboot-timeout", Value: "5000"},
	}...)

	return opts
}
