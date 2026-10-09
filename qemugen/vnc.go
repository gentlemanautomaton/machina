package qemugen

import (
	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/qemu/qguest"
)

func applyVNC(vnc machina.VNC, vars machina.Vars, t Target) error {
	if !vnc.Enabled() {
		return nil
	}

	// Enable the VNC protocol.
	t.VM.Settings.VNC = qguest.VNC{
		Enabled: true,
	}

	if vnc.UnixSocket.Enabled {
		t.VM.Settings.VNC.UnixSocketPath = t.VM.Settings.VNC.UnixSocketPath
	}

	if vnc.TCP.Enabled {
		port, err := vnc.TCP.EffectivePort(vars)
		if err != nil {
			return err
		}
		// QEMU assumes that that VNC is always served on TCP port number 5900
		// plus some display number. Here we subtract 5900 from the desired
		// port number to convert it into a display number.
		t.VM.Settings.VNC.TCPDisplayNumber = port - 5900
		t.VM.Settings.VNC.TCPAddr = "127.0.0.1"
	}

	if vnc.WebSocket.Enabled {
		port, err := vnc.WebSocket.EffectivePort(vars)
		if err != nil {
			return err
		}
		t.VM.Settings.VNC.WebSocketPort = port
		t.VM.Settings.VNC.WebSocketAddr = "127.0.0.1"
	}

	return nil
}
