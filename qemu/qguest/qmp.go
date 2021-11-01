package qguest

import (
	"github.com/gentlemanautomaton/machina/qemu"
	"github.com/gentlemanautomaton/machina/qemu/qhost/chardev"
)

// QMP holds QEMU Machine Protocol configuration for a QEMU virtual machine.
type QMP struct {
	Enabled bool
	Device  chardev.ID
	Mode    string
	Pretty  bool
}

// Parameters returns the parameters used for configuring QMP.
func (q QMP) Parameters() qemu.Parameters {
	var params qemu.Parameters

	params.Add("chardev", string(q.Device))
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

	return qemu.Options{{Type: "mon", Parameters: q.Parameters()}}
}
