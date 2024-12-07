package tpmdev

import (
	"github.com/gentlemanautomaton/machina/qemu"
)

// Backend identifies a QEMU TPM device backend.
type Backend string

// ID uniquely identifies a TPM device in QEMU's TPM device layer.
type ID string

// Property describes a QEMU TPM device property.
type Property = qemu.Parameter

// Properties hold a set of QEMU TPM device properties.
type Properties = qemu.Parameters

// Device is a Trusted Platform Module device on the QEMU host.
type Device interface {
	Backend() Backend
	ID() ID
	Properties() Properties
}
