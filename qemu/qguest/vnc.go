package qguest

import (
	"strconv"

	"github.com/gentlemanautomaton/machina/qemu"
)

// VNC configures the VNC protocol for external access to the virtual
// machine.
type VNC struct {
	Enabled          bool
	UnixSocketPath   string
	TCPAddr          string
	TCPDisplayNumber int
	WebSocketAddr    string
	WebSocketPort    int
}

// Parameters returns the parameters used for configuring spice.
func (vnc VNC) Parameters() qemu.Parameters {
	var params qemu.Parameters

	if vnc.UnixSocketPath != "" {
		params.AddValue("unix:" + vnc.UnixSocketPath)
	} else {
		display := vnc.TCPDisplayNumber
		if display < 0 {
			display = 0
		}
		if vnc.TCPAddr != "" {
			params.AddValue(vnc.TCPAddr + ":" + strconv.Itoa(display))
		} else {
			params.AddValue(":" + strconv.Itoa(display))
		}
	}

	if vnc.WebSocketPort > 0 {
		if vnc.WebSocketAddr != "" {
			params.Add("websocket", vnc.WebSocketAddr+":"+strconv.Itoa(vnc.WebSocketPort))
		} else {
			params.Add("websocket", strconv.Itoa(vnc.WebSocketPort))
		}
	}

	return params
}

// Options returns a set of QEMU virtual machine options for enabling spice.
func (vnc VNC) Options() qemu.Options {
	if !vnc.Enabled {
		return nil
	}

	return qemu.Options{{Type: "vnc", Parameters: vnc.Parameters()}}
}
