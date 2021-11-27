package qguest

import (
	"github.com/gentlemanautomaton/machina/qemu"
	"github.com/gentlemanautomaton/machina/qemu/qhost/chardev"
)

// QMP holds QEMU Machine Protocol configuration for a QEMU virtual machine.
type QMP struct {
	Enabled bool
	Devices []chardev.ID
	Mode    string
	Pretty  bool
}

// Parameters returns the parameters used for configuring QMP.
func (q QMP) Parameters(device chardev.ID) qemu.Parameters {
	var params qemu.Parameters

	params.Add("chardev", string(device))
	if q.Mode != "" {
		params.Add("mode", q.Mode)
	}
	if q.Pretty {
		params.AddValue("pretty")
	}

	return params
}

// Options returns a set of QEMU virtual machine options for specifying
// its QMP configuration.
func (q QMP) Options() qemu.Options {
	if !q.Enabled {
		return nil
	}

	var opts qemu.Options
	for _, device := range q.Devices {
		opts.Add("mon", q.Parameters(device)...)
	}

	return opts
}
