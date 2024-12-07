package qemugen

import (
	"errors"
	"fmt"

	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/qemu/qhost/chardev"
)

// https://wiki.qemu.org/Features/GuestAgent

func applyQEMUAgent(qga machina.QEMUAgent, vars machina.Vars, t Target) error {
	if !qga.Enabled {
		return nil
	}

	port, err := qga.EffectivePort(vars)
	if err != nil {
		return fmt.Errorf("failed to determine QEMU Agent port: %w", err)
	}

	if port == 0 {
		return errors.New("missing QEMU guest agent port")
	}

	// Grab a reference to the device registry for host characters devices
	registry := t.VM.Resources.CharDevs()

	// Prepare a communication channel for the host and guest
	socket, err := chardev.TCPSocket{
		ID:      chardev.ID("guestagent"),
		Host:    chardev.SocketHost("127.0.0.1"),
		Port:    chardev.SocketPort(port),
		Server:  true,
		NoWait:  true,
		NoDelay: true,
	}.AddTo(registry)
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
