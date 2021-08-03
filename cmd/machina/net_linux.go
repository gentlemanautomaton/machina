// +build linux

// https://github.com/gentlemanautomaton/vmcore-systemd/blob/master/bin/kvm-ifup.sh
// https://pkg.go.dev/github.com/vishvananda/netlink
// https://www.linux-kvm.org/page/Networking
// https://wiki.qemu.org/Documentation/Networking

package main

import (
	"fmt"

	"github.com/gentlemanautomaton/machina"
	"github.com/vishvananda/netlink"
)

func enableConnection(machine machina.MachineName, conn machina.Connection, sys machina.System) error {
	network, ok := sys.Network[conn.Network]
	if !ok {
		return fmt.Errorf("invalid network name \"%s\"", conn.Network)
	}

	// Find the link on the local system
	link, err := netlink.LinkByName(machina.MakeLinkName(machine, conn))
	if err != nil {
		return fmt.Errorf("link not found: %v", err)
	}

	// Find the bridge on the local system
	bridgeLink, err := netlink.LinkByName(network.Device)
	if err != nil {
		return fmt.Errorf("bridge not found: %v", err)
	}
	bridge, ok := bridgeLink.(*netlink.Bridge)
	if !ok {
		return fmt.Errorf("network link \"%s\" is not a bridge", network.Device)
	}

	// Turn up the link
	if err := netlink.LinkSetUp(link); err != nil {
		return fmt.Errorf("failed to turn the link up: %v", err)
	}

	// Add the link to the bridge
	if err := netlink.LinkSetMaster(link, bridge); err != nil {
		return fmt.Errorf("failed to add the link to the bridge: %v", err)
	}

	return nil
}

func disableConnection(machine machina.MachineName, conn machina.Connection, sys machina.System) error {
	// Find the link on the local system
	link, err := netlink.LinkByName(machina.MakeLinkName(machine, conn))
	if err != nil {
		return fmt.Errorf("link not found: %v", err)
	}

	// Remove the link from the bridge
	if err := netlink.LinkSetNoMaster(link); err != nil {
		return fmt.Errorf("failed to remove the link from the bridge: %v", err)
	}

	// Turn down the link
	if err := netlink.LinkSetDown(link); err != nil {
		return fmt.Errorf("failed to turn the link down: %v", err)
	}

	return nil
}
