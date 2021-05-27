package qguest

import (
	"strconv"

	"github.com/gentlemanautomaton/machina/qemu"
)

// Processor describes the processor configuration of a QEMU guest.
type Processor struct {
	Sockets int
	Cores   int
	Threads int
	HyperV  bool
}

// CPU returns the parameters of the CPU configuration. It includes any
// processor entitlements.
func (p Processor) CPU() qemu.Parameters {
	params := qemu.Parameters{{Name: "host"}}
	if p.HyperV {
		params = append(params, qemu.Parameters{
			{Name: "hv_relaxed"},
			{Name: "hv_spinlocks", Value: "0x1fff"},
			{Name: "hv_vapic"},
			{Name: "hv_time"},
		}...)
	}
	return params
}

// SMP returns the parameters for the desired level of simultaneous
// multithreading.
func (p Processor) SMP() qemu.Parameters {
	var params qemu.Parameters
	if p.Sockets > 0 {
		params.Add("sockets", strconv.Itoa(p.Sockets))
	}
	if p.Cores > 0 {
		params.Add("cores", strconv.Itoa(p.Cores))
	}
	if p.Threads > 0 {
		params.Add("threads", strconv.Itoa(p.Threads))
	}
	return params
}

// Options returns a set of QEMU virtual machine options for specifying
// its processor configuration.
func (p Processor) Options() qemu.Options {
	var opts qemu.Options

	if params := p.CPU(); len(params) > 0 {
		opts.Add("cpu", params...)
	}

	if params := p.SMP(); len(params) > 0 {
		opts.Add("smp", params...)
	}

	return opts
}
