package qdev_test

import (
	"fmt"

	"github.com/gentlemanautomaton/machina/qemu/qdev"
)

func ExampleQXL() {
	var topo qdev.Topology

	const heads = 2

	for i := 0; i < heads; i++ {
		// Add QXL display devices directly to the root complex
		if _, err := topo.AddQXL(); err != nil {
			panic(err)
		}
	}

	// Print the configuration
	options := topo.Options()
	for _, option := range options {
		fmt.Printf("%s\n", option)
	}

	// Output:
	// -device qxl-vga,id=qxl.0,bus=pcie.0,addr=1.0,multifunction=on
	// -device qxl,id=qxl.1,bus=pcie.0,addr=1.1
}
