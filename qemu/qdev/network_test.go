package qdev_test

import (
	"fmt"

	"github.com/gentlemanautomaton/machina/qemu/qdev"
	"github.com/gentlemanautomaton/machina/qemu/qhost"
)

func ExampleNetwork() {
	var (
		host qhost.Resources
		topo qdev.Topology
	)

	// Prepare a network tap that connects to the kvmbr0 interface
	tap, err := host.AddNetworkTap("kvmbr0", "", "")
	if err != nil {
		panic(err)
	}

	// Add a PCI Express Root Port that we'll connect the Network Controller to
	root, err := topo.AddRoot()
	if err != nil {
		panic(err)
	}

	// Add the Network Controller and connect it to the host network tap
	if _, err := root.AddVirtioNetwork("00:00:00:00:00:00", tap); err != nil {
		panic(err)
	}

	// Print the configuration
	options := topo.Options()
	for _, option := range options {
		fmt.Printf("%s\n", option)
	}

	// Output:
	// -device ioh3420,id=pcie.1.0,chassis=0,bus=pcie.0,addr=0.0,multifunction=on
	// -device virtio-net-pci,bus=pcie.1.0,mac=00:00:00:00:00:00,netdev=net.0
}
