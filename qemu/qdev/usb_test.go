package qdev_test

import (
	"fmt"

	"github.com/gentlemanautomaton/machina/qemu/qdev"
	"github.com/gentlemanautomaton/machina/qemu/qhost"
	"github.com/gentlemanautomaton/machina/qemu/qhost/blockdev"
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

func ExampleUSBAttachedSCSI() {
	var (
		resources qhost.Resources
		topo      qdev.Topology
	)

	// Grab a reference to the node graph for block devices.
	graph := resources.BlockDevs()

	// Prepare the disk's file protocol block device
	file, err := blockdev.File{
		Name:     "uas-disk",
		Path:     "uas-disk.iso",
		ReadOnly: true,
	}.Connect(graph)
	if err != nil {
		panic(err)
	}

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

	// Add a USB Attached SCSI Controller to the USB Controller
	scsi, err := usb.AddSCSI()
	if err != nil {
		panic(err)
	}

	// Add a SCSI disk to the SCSI controller
	if _, err := scsi.AddDisk(file); err != nil {
		panic(err)
	}

	// Print the configuration
	options := append(resources.Options(), topo.Options()...)
	for _, option := range options {
		fmt.Printf("%s\n", option)
	}

	// Output:
	// -blockdev driver=file,node-name=uas-disk,read-only=on,filename=uas-disk.iso
	// -device ioh3420,id=pcie.1.0,chassis=0,bus=pcie.0,addr=1.0,multifunction=on
	// -device qemu-xhci,id=usb,bus=pcie.1.0,p2=4,p3=4
	// -device usb-uas,id=uas,bus=usb.0,port=1
	// -device scsi-hd,id=uas.0.0,bus=uas.0,channel=0,scsi-id=0,lun=0,drive=uas-disk
}
