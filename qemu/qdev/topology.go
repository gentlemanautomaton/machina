package qdev

import (
	"errors"
	"fmt"

	"github.com/gentlemanautomaton/machina/qemu"
	"github.com/gentlemanautomaton/machina/qemu/qhost/blockdev"
	"github.com/gentlemanautomaton/machina/qemu/qhost/tpmdev"
)

// https://github.com/qemu/qemu/blob/master/docs/qdev-device-use.txt

var (
	// ErrRootComplexFull is returned when the addition of a new PCI Express
	// Root would exceed [MaxRootDevices].
	ErrRootComplexFull = errors.New("the PCI Express root complex is full and cannot accommodate more devices")
)

// Topology describes the PCI Express device topology for a virtual machine.
// It holds a set of PCI Express Root ports present in the PCI Express Root
// Complex.
//
// To add PCI Express devices, add PCI Express Roots then add devices
// to those roots.
//
// TODO: Consider representing the PCI Express Root Complex with its own
// struct.
type Topology struct {
	devices      []Device
	sata         int
	buses        BusMap
	allocated    AllocationSet
	reservations ReservationMap
	nextAddr     Addr
	nextIndex    int
}

// Reserve will attempt to reserve a PCI Express device address slot for
// specific device roles or types. Each slot can provide up to 8 devices via
// multifunction addressing. It can be called multiple times to reserve more
// than one slot for a role.
//
// A non-nil error will be returned if the slot is already allocated,
// assigned to a different role or exceeds [MaxRootDevices].
func (t *Topology) Reserve(role Role, slot int) error {
	t.init()

	if role == "" {
		return errors.New("a PCI Express device address reservation was requested for an empty role")
	}
	if slot < 0 || slot >= MaxRootDevices {
		return fmt.Errorf("a PCI Express device address reservation for the \"%s\" role was requested for slot %d, which is an invalid slot number", role, slot)
	}
	if t.allocated[slot] > 0 {
		return fmt.Errorf("a PCI Express device address reservation for the \"%s\" role cannot be created for slot %d because the slot is already in use", role, slot)
	}
	if existing := t.reservations.Role(slot); existing != "" {
		if role == existing {
			return nil // Ignore duplicate reservations.
		}
		return fmt.Errorf("a PCI Express device address reservation for the \"%s\" role cannot be created for slot %d because the slot is already reserved for the \"%s\" role", role, slot, existing)
	}

	t.reservations.Add(role, slot)

	return nil
}

// AddRoot adds a new PCI Express Root Port device to the PCI Express Root
// Complex.
//
// An error is returned if the addition would cause the root complex to exceed
// [MaxRootDevices].
//
// TODO: Consider allowing the caller to supply a preferred bus address or
// hotplug setting.
func (t *Topology) AddRoot() (*Root, error) {
	return t.AddRootForRole("")
}

// AddRoot adds a new PCI Express Root Port device to the PCI Express Root
// Complex. If role is non-empty and a reservation exists for that role, it
// will attempt to use a reserved address for the root.
//
// If a role is not supplied, or the number of device addresses reserved for
// the role have been exhausted, it will use the next available non-reserved
// PCI Express address that is available.
//
// An error is returned if the addition would cause the root complex to exceed
// [MaxRootDevices].
//
// TODO: Consider allowing the caller to supply a preferred hotplug setting.
func (t *Topology) AddRootForRole(role Role) (*Root, error) {
	t.init()

	addr, _, err := t.allocate(role)
	if err != nil {
		return nil, err
	}

	root := &Root{
		id:      ID(fmt.Sprintf("pcie.%d.%d", addr.Slot, addr.Function)),
		chassis: addr.Chassis(),
		addr:    addr,
		buses:   t.buses,
	}
	t.devices = append(t.devices, root)

	return root, nil
}

// AddVirtioVGA connects a PCI Express Virtio VGA display device to the
// PCI Express Root Complex.
func (t *Topology) AddVirtioVGA() (VGA, error) {
	t.init()

	addr, _, err := t.allocate(RoleDisplayDevice)
	if err != nil {
		return VGA{}, err
	}

	vga := VGA{
		id:   t.buses.Allocate("vga"),
		addr: addr,
	}
	t.devices = append(t.devices, vga)

	return vga, nil
}

// AddQXL connects a PCI Express QXL display device to the PCI Express Root
// Complex.
func (t *Topology) AddQXL() (QXL, error) {
	t.init()

	addr, _, err := t.allocate(RoleDisplayDevice)
	if err != nil {
		return QXL{}, err
	}

	secondary := t.buses.Count("qxl") > 0
	qxl := QXL{
		id:        t.buses.Allocate("qxl"),
		addr:      addr,
		secondary: secondary,
	}
	t.devices = append(t.devices, qxl)

	return qxl, nil
}

// AddPanic connects a paravirtualized panic device to the PCI Express Root
// Complex as an integrated PCI device.
func (t *Topology) AddPanic() (PVPanic, error) {
	t.init()

	addr, _, err := t.allocate("")
	if err != nil {
		return PVPanic{}, err
	}

	p := PVPanic{
		id:   t.buses.Allocate("panic"),
		addr: addr,
	}
	t.devices = append(t.devices, p)

	return p, nil
}

// AddTPMTIS connects a Trusted Platform Module device to the machine via memory mapping.
// The device will implement the TPM Interface Specification.
func (t *Topology) AddTPMTIS(device tpmdev.ID) (TPMTIS, error) {
	t.init()

	tpm := TPMTIS{
		device: device,
	}
	t.devices = append(t.devices, tpm)
	return tpm, nil
}

// AddTPMCRB connects a Trusted Platform Module device to the machine via memory mapping.
// The device will implement the TPM Command Response Buffer specification.
func (t *Topology) AddTPMCRB(device tpmdev.ID) (TPMCRB, error) {
	t.init()

	tpm := TPMCRB{
		device: device,
	}
	t.devices = append(t.devices, tpm)
	return tpm, nil
}

// AddCDROM connects a SATA CD-ROM device to the AHCI bus built into the
// q35 machine's ICH9 controller.
func (t *Topology) AddCDROM(bdev blockdev.Node) (SATACD, error) {
	t.init()

	if t.sata+1 > MaxSATADevices {
		return SATACD{}, ErrSATAFull
	}

	// On the q35 machine the built-in AHCI bus is named ide.1
	// https://bugzilla.redhat.com/show_bug.cgi?id=1368300
	cd := SATACD{
		id:       t.buses.Allocate("sata"),
		bus:      "ide.1",
		blockdev: bdev.Name(),
	}
	t.devices = append(t.devices, cd)

	t.sata++

	return cd, nil
}

// Devices returns all of the PCI Express Roots within the PCI Express Root
// Complex.
func (t *Topology) Devices() []Device {
	// Return a copy of the slice.
	devices := make([]Device, 0, len(t.devices))
	devices = append(devices, t.devices...)
	return devices
}

// Options returns a set of QEMU virtual machine options for creating the
// topology.
func (t *Topology) Options() qemu.Options {
	var opts qemu.Options

	Walk(t.Devices(), func(depth int, device Device) {
		props := device.Properties()
		if len(props) == 0 {
			return
		}

		opts = append(opts, qemu.Option{
			Type:       "device",
			Parameters: props,
		})
	})

	return opts
}

// init ensures that the topology data structures have been initialized.
func (t *Topology) init() {
	if t.devices == nil {
		t.devices = make([]Device, 0, MaxRootDevices*MaxDeviceFunctions)
		t.buses = make(BusMap)
		t.reservations = make(ReservationMap)

		// The first PCI Express device address is dedicated to the memory
		// controller host and platform. Make sure that we don't hand out
		// a device address that could conflict with it.
		t.allocated[0] = MaxDeviceFunctions
		t.nextAddr = Addr{Slot: 1, Function: 0}
	}
}

// allocate returns an available PCI Express bus address and device index.
// If role is non-empty, it allocates an address from the reservation if one
// is available. If a role is not provided, or no reservations are available,
// it returns the next available non-reservered address.
//
// It returns [ErrRootComplexFull] if all PCI Express devices within the
// topology have been allocated or reserved.
func (t *Topology) allocate(role Role) (addr Addr, index int, err error) {
	// If a role has been requested, look for an available device address
	// within a reservation first.
	if role != "" {
		slots := t.reservations[role]
		for _, slot := range slots {
			if count := t.allocated[slot]; count < MaxDeviceFunctions {
				// An available device address was found within the
				// reservation.
				addr, index = Addr{Slot: slot, Function: count}, t.nextIndex
				t.allocated.Add(addr)
				t.nextIndex++
				return
			}
		}
	}

	// Take the next address and index.
	addr, index = t.nextAddr, t.nextIndex

	// If the next address matches a reserved slot, search for a slot that
	// is available.
	if t.reservations.Role(addr.Slot) != "" {
		addr.Function = 0
		for addr.Slot < MaxRootDevices {
			addr.Slot++
			if t.reservations.Role(addr.Slot) == "" {
				break
			}
		}
	}

	// Make sure we don't exceed the root complex's device limit.
	if addr.Slot > MaxRootDevices {
		return Addr{}, 0, ErrRootComplexFull
	}

	// Increment the next address.
	t.nextAddr = addr.Next()
	t.nextIndex++

	return
}
