package qdev_test

import (
	"fmt"

	"github.com/gentlemanautomaton/machina/qemu/qdev"
)

func ExampleQXL() {
	var topo qdev.Topology

	const heads = 2

	for i := 0; i < heads; i++ {
		// Add a PCI Express Root Port that we'll connect the QXL Display to
		root, err := topo.AddRoot()
		if err != nil {
			panic(err)
		}

		// Add two QXL Display devices
		if _, err := root.AddQXL(); err != nil {
			panic(err)
		}
	}

	// Print the configuration
	options := topo.Options()
	for _, option := range options {
		fmt.Printf("%s\n", option)
	}

	// Output:
	// -device ioh3420,id=pcie.1.0,chassis=0,bus=pcie.0,addr=0.0,multifunction=on
	// -device qxl-vga,id=qxl.0,bus=pcie.1.0
	// -device ioh3420,id=pcie.1.1,chassis=0,bus=pcie.0,addr=0.1
	// -device qxl,id=qxl.1,bus=pcie.1.1
}
