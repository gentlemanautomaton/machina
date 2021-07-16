package qemugen

import (
	"github.com/gentlemanautomaton/machina/qemu/qdev"
	"github.com/gentlemanautomaton/machina/qemu/qvm"
)

// Target holds information about the target QEMU virtual machine that is
// being generated.
type Target struct {
	VM          *qvm.Definition
	Controllers *qdev.ControllerMap
}
