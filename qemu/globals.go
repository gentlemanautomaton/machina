package qemu

// Global describes a global configuration property for a QEMU driver.
type Global struct {
	Driver   string
	Property string
	Value    string
}

// Valid returns true if the global has a driver and property. It does not
// evaluate the semantic meaning and correctness of the global.
func (g Global) Valid() bool {
	return g.Driver != "" && g.Property != ""
}

// Option returns the global property as a QEMU option.
func (g Global) Option() Option {
	return Option{
		Type: "global",
		Parameters: Parameters{
			{Name: "driver", Value: g.Driver},
			{Name: "property", Value: g.Property},
			{Name: "value", Value: g.Value},
		},
	}
}

// Globals holds a set of global configuration properties for QEMU drivers.
type Globals []Global

// Add adds a global configuration property with the given driver, property
// and value.
//
// If driver or property are empty the global is not added.
func (globals *Globals) Add(driver, property, value string) {
	if driver == "" || property == "" {
		return
	}
	*globals = append(*globals, Global{Driver: driver, Property: property, Value: value})
}

// Options returns a set of QEMU virtual machine options that implement the
// global driver properties.
func (globals Globals) Options() Options {
	if len(globals) == 0 {
		return nil
	}
	opts := make(Options, 0, len(globals))
	for _, global := range globals {
		if global.Valid() {
			opts = append(opts, global.Option())
		}
	}
	return opts
}
