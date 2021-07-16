package qemugen

import (
	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/qemu/qvm"
)

func applyIdentity(m machina.Machine, vm *qvm.Definition) error {
	// Apply identity values
	vm.Settings.Identity.Name = string(m.Name)
	vm.Settings.Identity.ID = machina.UUID(m.ID)
	// TODO: Consider generating a random ID if one was not specified
	return nil
}
