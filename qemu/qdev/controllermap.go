package qdev

import (
	"github.com/gentlemanautomaton/machina/qemu/qhost"
)

// ControllerMap keeps track of which device controllers have been added to
// a QEMU virtual machine definition.
type ControllerMap struct {
	topo   *Topology
	serial *Serial
	scsi   map[qhost.IOThread]*SCSI
	usb    *USB
	uas    *USBAttachedSCSI
}

// NewControllerMap returns a controller map that will keep track of
// device controllers assigned to the given QEMU virtual machine definition.
//
// TODO: Consider walking the existing topology to find controllers.
func NewControllerMap(topo *Topology) *ControllerMap {
	return &ControllerMap{topo: topo}
}

// Topology returns the virtual machine topology that the map is bound to.
func (m *ControllerMap) Topology() *Topology {
	return m.topo
}

// Serial returns a serial controller device for the virtual machine.
func (m *ControllerMap) Serial() (*Serial, error) {
	if m.serial != nil {
		return m.serial, nil
	}

	// Add a PCI Express Root device that we'll connect the Serial
	// Controller to
	root, err := m.topo.AddRoot()
	if err != nil {
		return nil, err
	}

	// Add the Virtio Serial Controller
	serial, err := root.AddVirtioSerial()
	if err != nil {
		return nil, err
	}

	m.serial = serial

	return serial, nil
}

// SCSI returns a SCSI controller device for the virtual machine.
func (m *ControllerMap) SCSI(iothread qhost.IOThread) (*SCSI, error) {
	if m.scsi != nil {
		if controller, ok := m.scsi[iothread]; ok {
			return controller, nil
		}
	} else {
		m.scsi = make(map[qhost.IOThread]*SCSI, 4)
	}

	// Add a PCI Express Root device that we'll connect the SCSI Controller to
	root, err := m.topo.AddRoot()
	if err != nil {
		return nil, err
	}

	// Add the Virtio SCSI Controller with the given I/O thread
	scsi, err := root.AddVirtioSCSI(iothread)
	if err != nil {
		return nil, err
	}

	m.scsi[iothread] = scsi

	return scsi, nil
}

// USB returns a USB controller device for the virtual machine.
func (m *ControllerMap) USB() (*USB, error) {
	if m.usb != nil {
		return m.usb, nil
	}

	// Add a PCI Express Root device that we'll connect the USB
	// Controller to
	root, err := m.topo.AddRoot()
	if err != nil {
		return nil, err
	}

	// Add the Virtio USB Controller
	usb, err := root.AddUSB()
	if err != nil {
		return nil, err
	}

	m.usb = usb

	return usb, nil
}

// USBAttachedSCSI returns a USB Attached SCSI controller device for the
// virtual machine.
func (m *ControllerMap) USBAttachedSCSI() (*USBAttachedSCSI, error) {
	if m.uas != nil {
		return m.uas, nil
	}

	// Add a USB controller
	usb, err := m.USB()
	if err != nil {
		return nil, err
	}

	// Add the USB Attached SCSI controller
	uas, err := usb.AddSCSI()
	if err != nil {
		return nil, err
	}

	m.uas = uas

	return uas, nil
}
