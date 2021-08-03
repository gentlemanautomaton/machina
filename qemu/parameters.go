package qemu

import "strings"

// Parameter describes a parameter for a QEMU virtual machine option.
type Parameter struct {
	Name  string
	Value string
}

// String returns a string representation of the parameter.
//
// If both name and value are specified it returns a string in the form
// [name]=[value]. Otherwise, it returns [name] or [value], whichever is
// non-empty.
func (param Parameter) String() string {
	if param.Name == "" {
		return param.Value
	}
	if param.Value == "" {
		return param.Name
	}
	return param.Name + "=" + param.Value
}

// Add adds a named parameter with the give name and value.
//
// If name and value are both empty the property is not added.
func (params *Parameters) Add(name, value string) {
	if name == "" && value == "" {
		return
	}
	*params = append(*params, Parameter{Name: name, Value: value})
}

// Parameters hold a set of parameters for a QEMU virtual machine option.
type Parameters []Parameter

// AddValue adds a property with the given value but no name.
//
// If value is empty the property is not added.
func (params *Parameters) AddValue(value string) {
	if value == "" {
		return
	}
	*params = append(*params, Parameter{Value: value})
}

// String returns a string representation of the properties in the form
// expected by QEMU options.
func (params Parameters) String() string {
	if len(params) == 0 {
		return ""
	}

	list := make([]string, 0, len(params))
	for _, param := range params {
		if s := param.String(); s != "" {
			list = append(list, s)
		}
	}

	return strings.Join(list, ",")
}
