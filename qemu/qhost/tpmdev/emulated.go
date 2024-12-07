package tpmdev

import (
	"errors"
	"fmt"

	"github.com/gentlemanautomaton/machina/qemu/qhost/chardev"
)

// Emulated describes a Trusted Platform Module device that's bound to an
// emulator via a character device on the QEMU host.
type Emulated struct {
	ID     ID
	Device chardev.ID
}

// AddTo creates a new emulated TPM device with the given options and adds
// it to the TPM device registry.
//
// The returned TPM device is immutable and can safely be copied by value.
//
// An error is returned if the device cannot be added to the device registry
// or the TPM configuration is invalid.
func (e Emulated) AddTo(m Registry) (EmulatedDevice, error) {
	if err := e.validate(); err != nil {
		return EmulatedDevice{}, err
	}
	if m == nil {
		return EmulatedDevice{}, fmt.Errorf("a nil TPM device registry was provided when creating the \"%s\" emulated TPM", e.ID)
	}
	dev := EmulatedDevice{
		opts: e,
	}
	if err := m.Add(dev); err != nil {
		return EmulatedDevice{}, fmt.Errorf("failed to add the \"%s\" emulated TPM to the TPM device registry: %v", e.ID, err)
	}
	return dev, nil
}

func (e Emulated) validate() error {
	if e.ID == "" {
		return errors.New("the emulated TPM has an empty TPM device ID")
	}
	if e.Device == "" {
		return errors.New("the emulated TPM has an empty character device")
	}
	return nil
}

// EmulatedDevice is a Trusted Platform Module device that's bound to an
// emulator via a character device on the QEMU host.
type EmulatedDevice struct {
	opts Emulated
}

// Backend returns the name of the TPM device backend, emulated.
func (e EmulatedDevice) Backend() Backend {
	return "emulated"
}

// ID returns an ID that uniquely identifies the TPM device on the host.
func (e EmulatedDevice) ID() ID {
	return e.opts.ID
}

// Properties returns the TPM device properties of the emulated TPM.
func (e EmulatedDevice) Properties() Properties {
	props := Properties{
		{Name: string(e.Backend())},
		{Name: "id", Value: string(e.opts.ID)},
		{Name: "chardev", Value: string(e.opts.Device)},
	}
	return props
}
