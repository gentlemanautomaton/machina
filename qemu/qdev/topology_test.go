package qdev_test

import (
	"testing"

	"github.com/gentlemanautomaton/machina/qemu/qdev"
)

func TestTopologyAddressAllocation(t *testing.T) {
	var topo qdev.Topology

	addNineRoots := func() {
		for i := range 9 {
			if i%3 == 0 {
				topo.AddRoot()
			} else {
				topo.AddRootForRole("test")
			}
		}
	}

	addNineRootsForRole := func(role qdev.Role) {
		for range 9 {
			topo.AddRootForRole(role)
		}
	}

	topo.Reserve(qdev.RoleSCSIController, 2)
	topo.Reserve(qdev.RoleNetworkDevice, 3)
	topo.Reserve(qdev.RoleNetworkDevice, 4)

	addNineRoots()
	addNineRootsForRole(qdev.RoleNetworkDevice)
	addNineRootsForRole(qdev.RoleSCSIController)
	topo.AddVirtioVGA()
	addNineRootsForRole(qdev.RoleNetworkDevice)
	addNineRoots()

	expected := []string{
		"-device pcie-root-port,id=pcie.1.0,chassis=8,bus=pcie.0,addr=1.0,hotplug=off,multifunction=on",
		"-device pcie-root-port,id=pcie.1.1,chassis=9,bus=pcie.0,addr=1.1,hotplug=off",
		"-device pcie-root-port,id=pcie.1.2,chassis=10,bus=pcie.0,addr=1.2,hotplug=off",
		"-device pcie-root-port,id=pcie.1.3,chassis=11,bus=pcie.0,addr=1.3,hotplug=off",
		"-device pcie-root-port,id=pcie.1.4,chassis=12,bus=pcie.0,addr=1.4,hotplug=off",
		"-device pcie-root-port,id=pcie.1.5,chassis=13,bus=pcie.0,addr=1.5,hotplug=off",
		"-device pcie-root-port,id=pcie.1.6,chassis=14,bus=pcie.0,addr=1.6,hotplug=off",
		"-device pcie-root-port,id=pcie.1.7,chassis=15,bus=pcie.0,addr=1.7,hotplug=off",
		"-device pcie-root-port,id=pcie.5.0,chassis=40,bus=pcie.0,addr=5.0,hotplug=off,multifunction=on",
		"-device pcie-root-port,id=pcie.3.0,chassis=24,bus=pcie.0,addr=3.0,hotplug=off,multifunction=on",
		"-device pcie-root-port,id=pcie.3.1,chassis=25,bus=pcie.0,addr=3.1,hotplug=off",
		"-device pcie-root-port,id=pcie.3.2,chassis=26,bus=pcie.0,addr=3.2,hotplug=off",
		"-device pcie-root-port,id=pcie.3.3,chassis=27,bus=pcie.0,addr=3.3,hotplug=off",
		"-device pcie-root-port,id=pcie.3.4,chassis=28,bus=pcie.0,addr=3.4,hotplug=off",
		"-device pcie-root-port,id=pcie.3.5,chassis=29,bus=pcie.0,addr=3.5,hotplug=off",
		"-device pcie-root-port,id=pcie.3.6,chassis=30,bus=pcie.0,addr=3.6,hotplug=off",
		"-device pcie-root-port,id=pcie.3.7,chassis=31,bus=pcie.0,addr=3.7,hotplug=off",
		"-device pcie-root-port,id=pcie.4.0,chassis=32,bus=pcie.0,addr=4.0,hotplug=off,multifunction=on",
		"-device pcie-root-port,id=pcie.2.0,chassis=16,bus=pcie.0,addr=2.0,hotplug=off,multifunction=on",
		"-device pcie-root-port,id=pcie.2.1,chassis=17,bus=pcie.0,addr=2.1,hotplug=off",
		"-device pcie-root-port,id=pcie.2.2,chassis=18,bus=pcie.0,addr=2.2,hotplug=off",
		"-device pcie-root-port,id=pcie.2.3,chassis=19,bus=pcie.0,addr=2.3,hotplug=off",
		"-device pcie-root-port,id=pcie.2.4,chassis=20,bus=pcie.0,addr=2.4,hotplug=off",
		"-device pcie-root-port,id=pcie.2.5,chassis=21,bus=pcie.0,addr=2.5,hotplug=off",
		"-device pcie-root-port,id=pcie.2.6,chassis=22,bus=pcie.0,addr=2.6,hotplug=off",
		"-device pcie-root-port,id=pcie.2.7,chassis=23,bus=pcie.0,addr=2.7,hotplug=off",
		"-device pcie-root-port,id=pcie.5.1,chassis=41,bus=pcie.0,addr=5.1,hotplug=off",
		"-device virtio-vga,id=vga.0,bus=pcie.0,addr=5.2",
		"-device pcie-root-port,id=pcie.4.1,chassis=33,bus=pcie.0,addr=4.1,hotplug=off",
		"-device pcie-root-port,id=pcie.4.2,chassis=34,bus=pcie.0,addr=4.2,hotplug=off",
		"-device pcie-root-port,id=pcie.4.3,chassis=35,bus=pcie.0,addr=4.3,hotplug=off",
		"-device pcie-root-port,id=pcie.4.4,chassis=36,bus=pcie.0,addr=4.4,hotplug=off",
		"-device pcie-root-port,id=pcie.4.5,chassis=37,bus=pcie.0,addr=4.5,hotplug=off",
		"-device pcie-root-port,id=pcie.4.6,chassis=38,bus=pcie.0,addr=4.6,hotplug=off",
		"-device pcie-root-port,id=pcie.4.7,chassis=39,bus=pcie.0,addr=4.7,hotplug=off",
		"-device pcie-root-port,id=pcie.5.3,chassis=43,bus=pcie.0,addr=5.3,hotplug=off",
		"-device pcie-root-port,id=pcie.5.4,chassis=44,bus=pcie.0,addr=5.4,hotplug=off",
		"-device pcie-root-port,id=pcie.5.5,chassis=45,bus=pcie.0,addr=5.5,hotplug=off",
		"-device pcie-root-port,id=pcie.5.6,chassis=46,bus=pcie.0,addr=5.6,hotplug=off",
		"-device pcie-root-port,id=pcie.5.7,chassis=47,bus=pcie.0,addr=5.7,hotplug=off",
		"-device pcie-root-port,id=pcie.6.0,chassis=48,bus=pcie.0,addr=6.0,hotplug=off,multifunction=on",
		"-device pcie-root-port,id=pcie.6.1,chassis=49,bus=pcie.0,addr=6.1,hotplug=off",
		"-device pcie-root-port,id=pcie.6.2,chassis=50,bus=pcie.0,addr=6.2,hotplug=off",
		"-device pcie-root-port,id=pcie.6.3,chassis=51,bus=pcie.0,addr=6.3,hotplug=off",
		"-device pcie-root-port,id=pcie.6.4,chassis=52,bus=pcie.0,addr=6.4,hotplug=off",
		"-device pcie-root-port,id=pcie.6.5,chassis=53,bus=pcie.0,addr=6.5,hotplug=off",
	}

	options := topo.Options()
	if len(options) != len(expected) {
		t.Fail()
	}
	for i := range options {
		if got, want := options[i].String(), expected[i]; got != want {
			t.Errorf("unexpected topology option %d: \"%s\" (want \"%s\")", i, got, want)
		}
	}
}
