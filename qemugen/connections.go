package qemugen

import (
	"fmt"

	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/qemu/qhost"
)

func applyConnections(machine machina.MachineName, conns []machina.Connection, networks machina.NetworkMap, t Target) error {
	if len(conns) == 0 {
		return nil
	}

	// Add a netdev and device for each connection.
	for _, conn := range conns {
		// Find the connection's network
		network, ok := networks[conn.Network]
		if !ok {
			return fmt.Errorf("connection %s uses an unspecified machina network: %s", conn.Name, conn.Network)
		}

		// Determine the link name
		link := machina.MakeLinkName(machine, conn)

		// If up/down scripts were provided, use those
		up, down := qhost.NoScript, qhost.NoScript
		if network.Up != "" {
			up = qhost.Script(network.Up)
		}
		if network.Down != "" {
			down = qhost.Script(network.Down)
		}

		// Add the host's netdev resource for this connection
		tap, err := t.VM.Resources.AddNetworkTap(link, up, down)
		if err != nil {
			return err
		}

		// Add a PCI Express Root device that we'll connect a Network Controller
		// to.
		root, err := t.VM.Topology.AddRoot()
		if err != nil {
			return err
		}

		// Add a Virtio Network Controller.
		if _, err := root.AddVirtioNetwork(conn.MAC, tap); err != nil {
			return err
		}
	}

	return nil
}
