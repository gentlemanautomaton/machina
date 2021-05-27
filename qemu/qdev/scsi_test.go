package qdev_test

import (
	"fmt"

	"github.com/gentlemanautomaton/machina/qemu/qdev"
	"github.com/gentlemanautomaton/machina/qemu/qhost"
	"github.com/gentlemanautomaton/machina/qemu/qhost/blockdev"
)

func ExampleSCSI() {
	var (
		host qhost.Resources
		topo qdev.Topology
	)

	// Allocate an I/O thread for the SCSI Controller to use.
	ioThread, err := host.AddIOThread()
	if err != nil {
		panic(err)
	}

	// Prepare a block graph the describes the I/O processing for the
	// disk image on the QEMU host.
	graph := host.BlockDevs()

	// We'll refer to the final node in the graph by its node name.
	name := blockdev.NodeName("testdrive")

	// Add a file protocol node to the graph that reads and writes to our
	// disk image.
	file, err := blockdev.File{
		Name: name.Child("file"),
		Path: blockdev.FilePath("/tmp/test-drive.raw"),
	}.Connect(graph)
	if err != nil {
		panic(err)
	}

	// Add a raw format node to the graph that interprets the file data
	// as a raw image.
	drive, err := blockdev.Raw{Name: name}.Connect(file)
	if err != nil {
		panic(err)
	}

	// Add a PCI Express Root Port that we'll connect the SCSI Controller to
	root, err := topo.AddRoot()
	if err != nil {
		panic(err)
	}

	// Add a SCSI Controller with the I/O Thread that we prepared earlier
	scsi, err := root.AddVirtioSCSI(ioThread)
	if err != nil {
		panic(err)
	}

	// Attach a SCSI Disk to the controller that's backed by the drive we
	// prepared earlier
	if _, err := scsi.AddDisk(drive); err != nil {
		panic(err)
	}

	// Print the configuration
	options := topo.Options()
	for _, option := range options {
		fmt.Printf("%s\n", option)
	}

	// Output:
	// -device ioh3420,id=pcie.1.0,chassis=0,bus=pcie.0,addr=0.0,multifunction=on
	// -device virtio-scsi-pci,id=scsi.0,bus=pcie.1.0,iothread=iothread.0,num_queues=4
	// -device scsi-hd,id=scsi.0.0,bus=scsi.0,channel=0,scsi-id=0,lun=0,drive=testdrive
}
