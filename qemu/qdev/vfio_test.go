package qdev_test

import (
	"fmt"

	"github.com/gentlemanautomaton/machina/qemu/qdev"
)

func ExampleVFIO() {
	var topo qdev.Topology

	// Add a PCI Express Root Port that we'll connect the device to
	root, err := topo.AddRoot()
	if err != nil {
		panic(err)
	}

	// Add a VFIO PCI passthrough device to the root port
	if _, err := root.AddVFIO("/sys/bus/mdev/devices/79db70f1-92a5-4879-95f7-6325b66c1ff9"); err != nil {
		panic(err)
	}

	// Print the configuration
	options := topo.Options()
	for _, option := range options {
		fmt.Printf("%s\n", option)
	}

	// Output:
	// -device ioh3420,id=pcie.1.0,chassis=0,bus=pcie.0,addr=1.0,multifunction=on
	// -device vfio-pci,id=vfio.0,bus=pcie.1.0,sysfsdev=/sys/bus/mdev/devices/79db70f1-92a5-4879-95f7-6325b66c1ff9
}
