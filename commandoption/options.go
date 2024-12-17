package commandoption

import "strings"

// Data defines the data that is common to all command options.
//
// Data is a type alias, which allows it to be used directly when defining
// options that will match the Option interface:
//
//	type Option commandoption.Data
//
// For an option to be valid, it must have a type. Parameters are optional.
type Data = struct {
	Type       string
	Parameters Parameters
}

// Option is a type constraint that matches an option for a command.
//
// For an option to be valid, it must have a type. Parameters are optional.
type Option interface {
	~Data
	Prefix() string
	String() string
}

// String returns a string representation of the option.
//
// It returns an empty string if the option lacks a type.
func String[T Option](opt T) string {
	data := Data(opt)
	if data.Type == "" {
		return ""
	}

	switch params := data.Parameters.String(); params {
	case "":
		return opt.Prefix() + data.Type
	default:
		return opt.Prefix() + data.Type + " " + params
	}
}

// Options holds a set of command options.
type Options[T Option] []T

// Add adds an option with the given type and parameters.
//
// If type is empty the option is not added.
func (opts *Options[T]) Add(typ string, params ...Parameter) {
	if typ == "" {
		return
	}
	*opts = append(*opts, T{Type: typ, Parameters: params})
}

// Args returns the command line arguments for the command options.
func (opts Options[T]) Args() []string {
	if len(opts) == 0 {
		return nil
	}
	args := make([]string, 0, len(opts)*2)
	for _, opt := range opts {
		data := Data(opt)
		if data.Type != "" {
			args = append(args, opt.Prefix()+data.Type)
			if params := data.Parameters.String(); params != "" {
				args = append(args, params)
			}
		}
	}
	return args
}

// String returns a string representation of the command options.
func (opts Options[T]) String() string {
	return strings.Join(opts.Args(), " ")
}

// Multiline returns a multiline string for invocation of a command
// with the given options.
func (opts Options[T]) Multiline() string {
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
