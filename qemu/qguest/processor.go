package qguest

import (
	"strconv"
	"strings"

	"github.com/gentlemanautomaton/machina/qemu"
)

// Processor describes the processor configuration of a QEMU guest.
type Processor struct {
	Brand          string
	Sockets        int
	Cores          int
	ThreadsPerCore int
	HyperV         bool
}

// CPU returns the parameters of the CPU configuration. It includes any
// processor enlightenments.
func (p Processor) CPU() qemu.Parameters {
	params := qemu.Parameters{{Name: "host"}}

	// Simultaneous multi-threading on AMD processors must be enabled
	// explicitly.
	if p.ThreadsPerCore > 1 && strings.EqualFold(p.Brand, "AMD") {
		params = append(params, qemu.Parameter{Name: "topoext", Value: "on"})
	}

	// Simulate a Hyper-V hypervisor if requested.
	if p.HyperV {
		params = append(params, qemu.Parameters{
			{Name: "hv-relaxed"},
			{Name: "hv-vapic"},
			{Name: "hv-spinlocks", Value: "0x1fff"},
			{Name: "hv-vpindex"},
			{Name: "hv-runtime"},
			{Name: "hv-time"},
			{Name: "hv-synic"},
			{Name: "hv-stimer"},
			{Name: "hv-tlbflush"},
			{Name: "hv-ipi"},
			{Name: "hv-frequencies"},
			{Name: "hv-reenlightenment"},
			{Name: "hv-stimer-direct"},
			{Name: "hv-emsr-bitmap"},
			{Name: "hv-xmm-input"},
			{Name: "hv-tlbflush-ext"},
			{Name: "hv-tlbflush-direct"},
		}...)

		// Special Hyper-V enablements for AMD processors.
		if strings.EqualFold(p.Brand, "AMD") {
			params = append(params, qemu.Parameters{
				{Name: "hv-avic", Value: "on"},
			}...)
		}

		// Special Hyper-V enablements for Intel processors.
		if strings.EqualFold(p.Brand, "Intel") {
			params = append(params, qemu.Parameters{
				{Name: "hv-evmcs"},
			}...)
		}

		// TODO: Support hv-no-nonarch-coresharing once CPU pinning works.
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
	if p.ThreadsPerCore > 0 {
		params.Add("threads", strconv.Itoa(p.ThreadsPerCore))
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
