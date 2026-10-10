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

	// Reserve PCI Express bus addresses for display controllers, network
	// devices and SCSI Controllers. This helps to avoid bus address changes
	// that can cause guests perceive unwanted hardware changes.
	vm.Topology.Reserve(qdev.RoleDisplayDevice, 1)
	vm.Topology.Reserve(qdev.RoleVFIODevice, 2)
	vm.Topology.Reserve(qdev.RoleSerialController, 3)
	vm.Topology.Reserve(qdev.RoleUSBController, 4)
	vm.Topology.Reserve(qdev.RoleUSBController, 5)
	vm.Topology.Reserve(qdev.RoleBlockStorageDevice, qdev.MaxRootDevices-8)
	vm.Topology.Reserve(qdev.RoleBlockStorageDevice, qdev.MaxRootDevices-7)
	vm.Topology.Reserve(qdev.RoleBlockStorageDevice, qdev.MaxRootDevices-6)
	vm.Topology.Reserve(qdev.RoleBlockStorageDevice, qdev.MaxRootDevices-5)
	vm.Topology.Reserve(qdev.RoleSCSIController, qdev.MaxRootDevices-4)
	vm.Topology.Reserve(qdev.RoleSCSIController, qdev.MaxRootDevices-3)
	vm.Topology.Reserve(qdev.RoleNetworkDevice, qdev.MaxRootDevices-2)
	vm.Topology.Reserve(qdev.RoleNetworkDevice, qdev.MaxRootDevices-1)

	if err := applyDefaults(&vm); err != nil {
		return qvm.Definition{}, err
	}
	if err := applyIdentity(m, &vm); err != nil {
		return qvm.Definition{}, err
	}
	if err := applyFirmware(m.Info(), def, sys.Storage, target); err != nil {
		return qvm.Definition{}, err
	}
	if err := applyEarlyAttributes(m.Info(), def.Vars, def.Attributes, sys.Processor, target); err != nil {
		return qvm.Definition{}, err
	}
	if err := applyDisplay(def.Attributes, target); err != nil {
		return qvm.Definition{}, err
	}
	if err := applyDevices(def.Devices, sys.MediatedDevices, target); err != nil {
		return qvm.Definition{}, err
	}
	if err := applyLateAttributes(m.Info(), def.Vars, def.Attributes, sys.Processor, target); err != nil {
		return qvm.Definition{}, err
	}
	if err := applyVolumes(m.Info(), def.Vars, def.Volumes, sys.Storage, target); err != nil {
		return qvm.Definition{}, err
	}
	if err := applyConnections(m.Name, def.Connections, sys.Network, target); err != nil {
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

func applyEarlyAttributes(machine machina.MachineInfo, vars machina.Vars, attrs machina.Attributes, processors machina.ProcessorMap, target Target) error {
	// Apply CPU attributes.
	if err := applyCPU(attrs, processors, target); err != nil {
		return err
	}

	// Apply memory attributes.
	if ram := attrs.Memory.RAM; ram > 0 {
		target.VM.Settings.Memory.Allocation = qguest.MB(ram)
	}

	// Apply Trusted Platform Module attributes.
	if err := applyTPM(machine, attrs.TPM, target); err != nil {
		return err
	}

	// Apply QEMU Machine Protocol attributes.
	if err := applyQMP(machine, attrs.QMP, target); err != nil {
		return err
	}

	// Apply guest agent attributes.
	if err := applyQEMUAgent(attrs.Agent.QEMU, vars, target); err != nil {
		return err
	}

	return nil
}

func applyLateAttributes(machine machina.MachineInfo, vars machina.Vars, attrs machina.Attributes, processors machina.ProcessorMap, target Target) error {
	// Apply remote protocol USB attributes.
	if err := applyRemoteUSB(attrs, vars, target); err != nil {
		return err
	}

	// Apply VNC protocol attributes.
	if err := applyVNC(attrs.VNC, vars, target); err != nil {
		return err
	}

	// Apply spice protocol attributes.
	if err := applySpice(attrs.Spice, vars, target); err != nil {
		return err
	}

	return nil
}
