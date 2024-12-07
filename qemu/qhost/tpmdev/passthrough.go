package tpmdev

import (
	"errors"
	"fmt"

	"github.com/gentlemanautomaton/machina/filesystem/devfs"
	"github.com/gentlemanautomaton/machina/filesystem/sysfs"
)

// Passthrough describes a Trusted Platfrom Module device on the QEMU host
// that will be passed into and directly used by a guest.
type Passthrough struct {
	ID         ID
	Path       devfs.Path
	CancelPath sysfs.Path
}

// Add creates a new passthrough TPM device with the given options and adds
// it to the TPM device registry.
//
// The returned TPM device is immutable and can safely be copied by value.
//
// An error is returned if the device cannot be added to the device registry
// or the TPM configuration is invalid.
func (p Passthrough) Add(m Registry) (PassthroughDevice, error) {
	if err := p.validate(); err != nil {
		return PassthroughDevice{}, err
	}
	if m == nil {
		return PassthroughDevice{}, fmt.Errorf("a nil TPM device registry was provided when creating the \"%s\" emulated TPM", p.ID)
	}
	dev := PassthroughDevice{
		opts: p,
	}
	if err := m.Add(dev); err != nil {
		return PassthroughDevice{}, fmt.Errorf("failed to add the \"%s\" emulated TPM to the TPM device registry: %v", p.ID, err)
	}
	return dev, nil
}

func (p Passthrough) validate() error {
	if p.ID == "" {
		return errors.New("the emulated TPM has an empty TPM device ID")
	}
	return nil
}

// PassthroughDevice is a Trusted Platfrom Module device on the QEMU host
// that will be passed into and directly used by a guest.
type PassthroughDevice struct {
	opts Passthrough
}

// Backend returns the name of the TPM device backend, passthrough.
func (p PassthroughDevice) Backend() Backend {
	return "passthrough"
}

// ID returns an ID that uniquely identifies the TPM device on the host.
func (p PassthroughDevice) ID() ID {
	return p.opts.ID
}

// Properties returns the TPM device properties of the passthrough TPM.
func (p PassthroughDevice) Properties() Properties {
	props := Properties{
		{Name: string(p.Backend())},
		{Name: "id", Value: string(p.opts.ID)},
	}
	if p.opts.Path != "" {
		props = append(props, Property{Name: "path", Value: string(p.opts.Path)})
	}
	if p.opts.CancelPath != "" {
		props = append(props, Property{Name: "cancel-path", Value: string(p.opts.CancelPath)})
	}
	return props
}
