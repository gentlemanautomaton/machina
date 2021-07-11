package qdev

import (
	"errors"

	"github.com/gentlemanautomaton/machina/qemu/qhost/blockdev"
)

// https://bugzilla.redhat.com/show_bug.cgi?id=1368300

const (
	// MaxSATADevices is the maximum number of SATA devices supported by the
	// ICH9's SATA controller on q35 machines.
	MaxSATADevices = 6
)

var (
	// ErrSATAFull is returned when the addition of a new SATA device
	// would exceed MaxSATADevices.
	ErrSATAFull = errors.New("the SATA controller is full and cannot accommodate more devices")
)

// SATACD is a SATA CD-ROM device.
type SATACD struct {
	id       ID
	bus      ID
	blockdev blockdev.NodeName
}

// Driver returns the driver for the SATA CD device, ide-cd.
//
// The Q35 machine foolishly uses "ide" as the identifier for its AHCI bus
// that's built into its ICH9 controller. As such, ide-cd devices are actually
// interpreted as SATA CD devices on this architecture.
//
// See: https://bugzilla.redhat.com/show_bug.cgi?id=1368300
func (cd SATACD) Driver() Driver {
	return "ide-cd"
}

// Properties returns the properties of the SATA CD device.
func (cd SATACD) Properties() Properties {
	return Properties{
		{Name: string(cd.Driver())},
		{Name: "id", Value: string(cd.id)},
		{Name: "bus", Value: string(cd.bus)},
		{Name: "drive", Value: string(cd.blockdev)},
	}
}
