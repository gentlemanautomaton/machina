package qguest

import (
	"github.com/gentlemanautomaton/machina/qemu"
	"github.com/gentlemanautomaton/machina/qemu/qhost/blockdev"
)

// https://listman.redhat.com/archives/libvir-list/2019-January/msg01292.html
// https://bugzilla.redhat.com/show_bug.cgi?id=1686552
// https://github.com/qemu/qemu/search?q=pflash0&type=commits
// https://github.com/qemu/qemu/commit/ebc29e1beab02646702c8cb9a1d29b68f72ad503

// Firmware holds firmware configuration for a QEMU virtual machine.
type Firmware struct {
	Code blockdev.NodeName
	Vars blockdev.NodeName
}

// MachineParameters returns a set of machine parameters for the firmware.
func (f Firmware) MachineParameters() qemu.Parameters {
	var params qemu.Parameters
	if f.Code != "" {
		params.Add("pflash0", string(f.Code))
		if f.Vars != "" {
			params.Add("pflash1", string(f.Vars))
		}
	}
	return params
}
