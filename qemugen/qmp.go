package qemugen

import (
	"fmt"

	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/qemu/qguest"
	"github.com/gentlemanautomaton/machina/qemu/qhost/chardev"
)

// https://qemu-project.gitlab.io/qemu/system/security.html#monitor-console-qmp-and-hmp

func applyQMP(machine machina.MachineInfo, proto machina.QMP, t Target) error {
	if !proto.Enabled {
		return nil
	}

	// Grab a reference to the device registry for host characters devices
	registry := t.VM.Resources.CharDevs()

	// Determine the file system path to the socket
	socketPaths := proto.AllSocketPaths(machine)
	if len(socketPaths) == 0 {
		return fmt.Errorf("no socket paths were provided for QMP")
	}

	// Enable the QMP server
	t.VM.Settings.QMP = qguest.QMP{
		Enabled: true,
		Mode:    "control",
	}

	for i, socketPath := range socketPaths {
		// Prepare a unix domain socket for QMP
		socket, err := chardev.UnixSocket{
			ID:     chardev.ID(fmt.Sprintf("qmp.%d", i)),
			Path:   chardev.SocketPath(socketPath),
			Server: true,
			NoWait: true,
		}.AddTo(registry)
		if err != nil {
			return err
		}

		// Add the socket to the QMP server
		t.VM.Settings.QMP.Devices = append(t.VM.Settings.QMP.Devices, socket.ID())
	}

	return nil
}
