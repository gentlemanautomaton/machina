package qdev

import "github.com/gentlemanautomaton/machina/qemu/qhost/tpmdev"

// TPM is a Trusted Platform Module device that is accessible to the guest
// through a mapped memory region.
type TPM struct {
	device tpmdev.ID
}

// Driver returns the driver for the TPM device, tpm-tis.
func (tpm TPM) Driver() Driver {
	return "tpm-tis"
}

// Properties returns the properties of the TPM device.
func (tpm TPM) Properties() Properties {
	props := Properties{
		{Name: string(tpm.Driver())},
		{Name: "tpmdev", Value: string(tpm.device)},
	}
	return props
}
