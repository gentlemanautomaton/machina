package qemugen

import (
	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/qemu/qdev"
	"github.com/gentlemanautomaton/machina/qemu/qguest"
	"github.com/gentlemanautomaton/machina/qemu/qvm"
)

// Build prepares a QEMU virtual machine definition for the given machina
// machine and system configuration.
func Build(m machina.Machine, sys machina.System) (qvm.Definition, error) {
	def, err := machina.Build(m, sys)
	if err != nil {
		return qvm.Definition{}, err
	}

	vm := qvm.Definition{}
	controllers := qdev.NewControllerMap(&vm.Topology)
	target := Target{VM: &vm, Controllers: controllers, BootOrder: new(qdev.BootOrder)}

	if err := applyDefaults(&vm); err != nil {
		return qvm.Definition{}, err
	}
	if err := applyIdentity(m, &vm); err != nil {
		return qvm.Definition{}, err
	}
	if err := applyFirmware(m.Info(), def, sys.Storage, target); err != nil {
		return qvm.Definition{}, err
	}
	if err := applyAttributes(m.Info(), def.Vars, def.Attributes, target); err != nil {
		return qvm.Definition{}, err
	}
	if err := applyVolumes(m.Info(), def.Vars, def.Volumes, sys.Storage, target); err != nil {
		return qvm.Definition{}, err
	}
	if err := applyConnections(m.Name, def.Connections, sys.Network, target); err != nil {
		return qvm.Definition{}, err
	}
	if err := applyDevices(def.Devices, sys.MediatedDevices, target); err != nil {
		return qvm.Definition{}, err
	}

	return vm, nil
}

func applyDefaults(vm *qvm.Definition) error {
	// Clock settings
	vm.Settings.Clock.Base = qguest.ClockBaseUTC
	vm.Settings.Clock.Isolation = qguest.ClockIsolationHost
	vm.Settings.Clock.DriftFix = qguest.ClockDriftFixSlew

	// Paravirtualized panic support
	if _, err := vm.Topology.AddPanic(); err != nil {
		return err
	}

	return nil
}

func applyAttributes(machine machina.MachineInfo, vars machina.Vars, attrs machina.Attributes, target Target) error {
	// Apply CPU attributes
	if sockets := attrs.CPU.Sockets; sockets > 0 {
		target.VM.Settings.Processor.Sockets = sockets
	}
	if cores := attrs.CPU.Cores; cores > 0 {
		target.VM.Settings.Processor.Cores = cores
	}
	if threads := attrs.CPU.Threads; threads > 0 {
		target.VM.Settings.Processor.Threads = threads
	}
	if attrs.Enlightenments.Enabled {
		target.VM.Settings.Processor.HyperV = true
	}

	// Apply memory attributes
	if ram := attrs.Memory.RAM; ram > 0 {
		target.VM.Settings.Memory.Allocation = qguest.MB(ram)
	}

	// Apply QEMU Machine Protocol attributes
	if err := applyQMP(machine, attrs.QMP, target); err != nil {
		return err
	}

	// Apply guest agent attributes
	if err := applyQEMUAgent(attrs.Agent.QEMU, vars, target); err != nil {
		return err
	}

	// Apply spice protocol attributes
	if err := applySpice(attrs.Spice, vars, target); err != nil {
		return err
	}

	return nil
}
