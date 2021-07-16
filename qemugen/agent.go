package qemugen

import (
	"errors"

	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/qemu/qhost/chardev"
)

// https://wiki.qemu.org/Features/GuestAgent

func applyQEMUAgent(qga machina.QEMUAgent, t Target) error {
	if !qga.Enabled {
		return nil
	}

	if qga.Port == 0 {
		return errors.New("missing QEMU guest agent port")
	}

	// Grab a reference to the device registry for host characters devices
	registry := t.VM.Resources.CharDevs()

	// Prepare a communication channel for the host and guest
	socket, err := chardev.TCPSocket{
		ID:      chardev.ID("guestagent"),
		Host:    chardev.SocketHost("127.0.0.1"),
		Port:    chardev.SocketPort(qga.Port),
		Server:  true,
		NoWait:  true,
		NoDelay: true,
	}.Add(registry)
	if err != nil {
		return err
	}

	// Add a Virtio Serial Controller
	serial, err := t.Controllers.Serial()
	if err != nil {
		return err
	}

	// Add a serial port that's connected to the vdagent channel
	if _, err := serial.AddPort(socket.ID(), "org.qemu.guest_agent.0"); err != nil {
		return err
	}

	return nil
}
