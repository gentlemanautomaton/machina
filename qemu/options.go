package qemu

import "github.com/gentlemanautomaton/machina/commandoption"

// TODO: Someday consider outputting configuration in the readconfig format,
// whenver that gets finalized. See:
// https://lists.gnu.org/archive/html/qemu-devel/2020-11/msg02934.html

// Option is an option for a QEMU virtual machine.
//
// For an option to be valid, it must have a type. Parameters are optional.
type Option commandoption.Data

// String returns a string representation of the option.
//
// It returns an empty string if the option lacks a type.
func (opt Option) String() string {
	return commandoption.String(opt)
}

// Prefix returns the option prefix used by QEMU, which is a single
// dash character "-".
func (opt Option) Prefix() string {
	return "-"
}

// Options holds a set of configuration options for a QEMU virtual machine.
type Options = commandoption.Options[Option]

// Parameter describes a parameter for a QEMU virtual machine option.
type Parameter = commandoption.Parameter

// Parameters hold a set of parameters for a QEMU virtual machine option.
type Parameters = commandoption.Parameters
