package qemu_test

import (
	"fmt"

	"github.com/gentlemanautomaton/machina/qemu/qguest"
	"github.com/gentlemanautomaton/machina/qemu/qhost"
	"github.com/gentlemanautomaton/machina/qemu/qhost/blockdev"
	"github.com/gentlemanautomaton/machina/qemu/qvm"
	"github.com/google/uuid"
)

func Example() {
	// Prepare a QEMU virtual machine definition
	var vm qvm.Definition

	// Specify some basic settings for the guest
	vm.Settings = qguest.Settings{
		Identity: qguest.Identity{
			Name: "test-machine",
			ID:   uuid.MustParse("{00000000-0000-0000-0000-000000000001}"),
		},
		Processor: qguest.Processor{
			Sockets: 1,
			Cores:   24,
		},
		Memory: qguest.Memory{
			Allocation: qguest.GB(2),
		},
	}

	// Add a network controller with a network tap on kvmbr0
	{
		tap, err := vm.Resources.AddNetworkTap(
			"kvmbr0",
			qhost.Script("/etc/qemu/if-up.sh"),
			qhost.Script("/etc/qemu/if-down.sh"),
		)
		if err != nil {
			panic(err)
		}

		root, err := vm.Topology.AddRoot()
		if err != nil {
			panic(err)
		}

		if _, err := root.AddVirtioNetwork("00:00:00:00:00:00", tap); err != nil {
			panic(err)
		}
	}

	// Add an OS storage volume as a Virtio SCSI disk backed by a raw format
	// disk image in the home directory
	{
		graph := vm.Resources.BlockDevs()
		name := blockdev.NodeName("test-os")

		// Add a file node to the block device graph
		file, err := blockdev.File{
			Name: name.Child("file"), // Derive a unique node name
			Path: blockdev.FilePath("~/test-os.raw"),
		}.Connect(graph)
		if err != nil {
			panic(err)
		}

		// Add a raw format node to the block device graph that connects
		// to the file node we just created
		disk, err := blockdev.Raw{Name: name}.Connect(file)
		if err != nil {
			panic(err)
		}

		// Add an iothread that will be used by the SCSI Controller
		iothread, err := vm.Resources.AddIOThread()
		if err != nil {
			panic(err)
		}

		// Add a root port for the SCSI controller
		root, err := vm.Topology.AddRoot()
		if err != nil {
			panic(err)
		}

		// Connect the SCSI controller to the root port
		scsi, err := root.AddVirtioSCSI(iothread)
		if err != nil {
			panic(err)
		}

		// Add a SCSI disk to the controller that is backed by the format node
		// we prepared earlier
		if _, err := scsi.AddDisk(disk); err != nil {
			panic(err)
		}
	}

	// Print each option on its own line
	for _, opt := range vm.Options() {
		fmt.Printf("%s \\\n", opt)
	}

	// Output:
	// -uuid 00000000-0000-0000-0000-000000000001 \
	// -name test-machine \
	// -enable-kvm \
	// -machine type=q35,accel=kvm \
	// -cpu host \
	// -smp sockets=1,cores=24 \
	// -m size=2G \
	// -nodefaults \
	// -nographic \
	// -object iothread,id=iothread.0 \
	// -blockdev driver=file,node-name=test-os-file,filename=~/test-os.raw \
	// -blockdev driver=raw,node-name=test-os,file=test-os-file \
	// -netdev tap,id=net.0,ifname=kvmbr0,script=/etc/qemu/if-up.sh,downscript=/etc/qemu/if-down.sh \
	// -device ioh3420,id=pcie.1.0,chassis=0,bus=pcie.0,addr=0.0,multifunction=on \
	// -device virtio-net-pci,bus=pcie.1.0,mac=00:00:00:00:00:00,netdev=net.0 \
	// -device ioh3420,id=pcie.1.1,chassis=0,bus=pcie.0,addr=0.1 \
	// -device virtio-scsi-pci,id=scsi.0,bus=pcie.1.1,iothread=iothread.0,num_queues=4 \
	// -device scsi-hd,id=scsi.0.0,bus=scsi.0,channel=0,scsi-id=0,lun=0,drive=test-os \
}
