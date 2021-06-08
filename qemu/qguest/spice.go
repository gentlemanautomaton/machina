package qguest

import (
	"strconv"

	"github.com/gentlemanautomaton/machina/qemu"
)

// Spice configures the spice protocol for external access to the virtual
// machine.
type Spice struct {
	Enabled             bool
	Port                int
	Addr                string
	DisableTicketing    bool
	DisableCopyPaste    bool
	DisableFileTransfer bool
}

// Parameters returns the parameters used for configuring spice.
func (s Spice) Parameters() qemu.Parameters {
	var params qemu.Parameters

	if s.Port > 0 {
		params.Add("port", strconv.Itoa(s.Port))
	}
	if s.Addr != "" {
		params.Add("addr", s.Addr)
	}
	if s.DisableTicketing {
		params.AddValue("disable-ticketing")
	}
	if s.DisableCopyPaste {
		params.AddValue("disable-copy-paste")
	}
	if s.DisableFileTransfer {
		params.AddValue("disable-agent-file-xfer")
	}

	return params
}

// Options returns a set of QEMU virtual machine options for enabling spice.
func (s Spice) Options() qemu.Options {
	if !s.Enabled {
		return nil
	}

	return qemu.Options{{Type: "spice", Parameters: s.Parameters()}}
}
