package qemu

import "strings"

// TODO: Someday consider outputting configuration in the readconfig format,
// whenver that gets finalized. See:
// https://lists.gnu.org/archive/html/qemu-devel/2020-11/msg02934.html

// Option is an option for a QEMU virtual machine.
//
// For an option to be valid, it must have a type. Parameters are optional.
type Option struct {
	Type       string
	Parameters Parameters
}

// String returns a string representation of the option.
//
// It returns an empty string if the option is invalid.
func (opt Option) String() string {
	if opt.Type == "" {
		return ""
	}

	switch params := opt.Parameters.String(); params {
	case "":
		return "-" + opt.Type
	default:
		return "-" + opt.Type + " " + params
	}
}

// Valid returns true if the option has a type. It does not evaluate
// the semantic meaning and correctness of the option.
func (opt Option) Valid() bool {
	return opt.Type != ""
}

// Options holds a set of configuration options for a QEMU virtual machine.
type Options []Option

// Add adds an option with the given type and parameters.
//
// If type is empty the property is not added.
func (opts *Options) Add(typ string, params ...Parameter) {
	if typ == "" {
		return
	}
	*opts = append(*opts, Option{Type: typ, Parameters: params})
}

// Args returns the command line arguments for invocation of a QEMU virtual
// machine with the given options.
func (opts Options) Args() []string {
	if len(opts) == 0 {
		return nil
	}
	args := make([]string, 0, len(opts)*2)
	for _, opt := range opts {
		if opt.Valid() {
			args = append(args, "-"+opt.Type)
			if params := opt.Parameters.String(); params != "" {
				args = append(args, params)
			}
		}
	}
	return args
}

// String returns a multiline string for invocation of a QEMU virtual
// machine with the given options.
func (opts Options) String() string {
	var b strings.Builder
	for i, option := range opts {
		last := i == len(opts)-1
		b.WriteString(option.String())
		if !last {
			b.WriteString(" \\\n")
		}
	}
	return b.String()
}
