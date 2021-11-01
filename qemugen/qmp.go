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
	socketPath := proto.MakeSocketPath(machine)
	if socketPath == "" {
		return fmt.Errorf("a socket path was not provided for QMP")
	}

	// Prepare a unix domain socket for QMP
	socket, err := chardev.UnixSocket{
		ID:     chardev.ID("qmp"),
		Path:   chardev.SocketPath(socketPath),
		Server: true,
		NoWait: true,
	}.Add(registry)
	if err != nil {
		return err
	}

	// Enable the QMP server
	t.VM.Settings.QMP = qguest.QMP{
		Enabled: true,
		Device:  socket.ID(),
		Mode:    "control",
	}

	return nil
}
