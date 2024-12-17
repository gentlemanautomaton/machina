package swtpmconfigfile

import "strings"

// Option is an option in a software TPM configuration file.
//
// For an option to be valid, it must have a name.
type Option struct {
	Name  string
	Value string
}

// String returns a string representation of the option in the form
// NAME=VALUE.
//
// It returns an empty string if the option lacks a name.
func (opt Option) String() string {
	if opt.Name == "" {
		return ""
	}

	return opt.Name + "=" + opt.Value
}

// Options holds a set of configuration options for a software TPM
// configuration file.
type Options []Option

// Add adds an option with the given name and value.
//
// If name is empty the option is not added.
func (opts *Options) Add(name, value string) {
	if name == "" {
		return
	}
	*opts = append(*opts, Option{Name: name, Value: value})
}

// String returns a multiline string that defines each of the contained
// options on a separate line.
func (opts Options) String() string {
	var b strings.Builder
	for _, option := range opts {
		b.WriteString(option.String())
		b.WriteString("\n")
	}
	return b.String()
}
