package qemugen

import (
	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/qemu/qhost/chardev"
	"github.com/gentlemanautomaton/machina/qemu/qhost/tpmdev"
)

// https://qemu-project.gitlab.io/qemu/specs/tpm.html

func applyTPM(machine machina.MachineInfo, proto machina.TPM, t Target) error {
	if !proto.Enabled {
		return nil
	}

	// Grab a reference to the device registry for host character devices.
	charDevs := t.VM.Resources.CharDevs()

	// Connect to the swtpm process that serves as a TPM emulator via a unix
	// socket at a well-known location.
	socket, err := chardev.UnixSocket{
		ID:   chardev.ID("tpm.0.socket"),
		Path: chardev.SocketPath(proto.SocketPath(machine)),
	}.Add(charDevs)
	if err != nil {
		return err
	}

	// Grab a reference to the device registry for host TPM devices.
	tpmDevs := t.VM.Resources.TPMDevs()
	tpm, err := tpmdev.Emulated{
		ID:     tpmdev.ID("tpm.0"),
		Device: socket.ID(),
	}.Add(tpmDevs)
	if err != nil {
		return err
	}

	// Add a memory-mapped TPM TIS device to the machine.
	if _, err := t.VM.Topology.AddTPM(tpm.ID()); err != nil {
		return err
	}

	return nil
}
