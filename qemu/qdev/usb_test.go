package qdev_test

import (
	"fmt"

	"github.com/gentlemanautomaton/machina/qemu/qdev"
)

func ExampleUSB() {
	var topo qdev.Topology

	// Add a PCI Express Root Port that we'll connect the USB Controller to
	root, err := topo.AddRoot()
	if err != nil {
		panic(err)
	}

	// Add a USB Controller
	usb, err := root.AddUSB()
	if err != nil {
		panic(err)
	}

	// Attach a USB Tablet to the controller
	if _, err := usb.AddTablet(); err != nil {
		panic(err)
	}

	// Print the configuration
	options := topo.Options()
	for _, option := range options {
		fmt.Printf("%s\n", option)
	}

	// Output:
	// -device ioh3420,id=pcie.1.0,chassis=0,bus=pcie.0,addr=1.0,multifunction=on
	// -device qemu-xhci,id=usb,bus=pcie.1.0,p2=4,p3=4
	// -device usb-tablet,id=usb.0.1,bus=usb.0,port=1
}
