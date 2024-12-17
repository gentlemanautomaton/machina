package swtpmcert

import (
	"github.com/gentlemanautomaton/machina/commandoption"
)

// Option is an option for a software TPM emulator.
//
// For an option to be valid, it must have a type. Parameters are optional.
type Option commandoption.Data

// String returns a string representation of the option.
//
// It returns an empty string if the option lacks a type.
func (opt Option) String() string {
	return commandoption.String(opt)
}

// Prefix returns the option prefix used by swtpm_cert, which is a double
// dash character "--".
func (opt Option) Prefix() string {
	return "--"
}

// Options holds a set of configuration options for a software TPM emulator.
type Options = commandoption.Options[Option]

// Parameter describes a parameter for a swtpm option.
type Parameter = commandoption.Parameter

// Parameters hold a set of parameters for a swtpm option.
type Parameters = commandoption.Parameters
