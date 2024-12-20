package qdev

import "github.com/gentlemanautomaton/machina/qemu/qhost/tpmdev"

// TPMTIS is a Trusted Platform Module device that is accessible to the guest
// through a mapped memory region. It implements the TPM Interface
// Specification.
type TPMTIS struct {
	device tpmdev.ID
}

// Driver returns the driver for the TPM device, tpm-tis.
func (tpm TPMTIS) Driver() Driver {
	return "tpm-tis"
}

// Properties returns the properties of the TPM device.
func (tpm TPMTIS) Properties() Properties {
	props := Properties{
		{Name: string(tpm.Driver())},
		{Name: "tpmdev", Value: string(tpm.device)},
	}
	return props
}

// TPMCRB is a Trusted Platform Module device that is accessible to the guest
// through a mapped memory region. It implements the TPM Command Response
// Buffer specification.
type TPMCRB struct {
	device tpmdev.ID
}

// Driver returns the driver for the TPM device, tpm-crb.
func (tpm TPMCRB) Driver() Driver {
	return "tpm-crb"
}

// Properties returns the properties of the TPM device.
func (tpm TPMCRB) Properties() Properties {
	props := Properties{
		{Name: string(tpm.Driver())},
		{Name: "tpmdev", Value: string(tpm.device)},
	}
	return props
}
