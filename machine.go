package machina

// Machine describes an individual virtual machine in machina. It contains
// identity information, tags, and a definition.
//
// If tags are present, the machine's definition can be merged with its tag
// definitions through use of the Build function.
//
// The Machine structure is intended to be marshaled to and from JSON. It
// defines the format of files in the machina.conf.d directory.
type Machine struct {
	Name        MachineName `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	ID          MachineID   `json:"id,omitempty"`
	Tags        []Tag       `json:"tags,omitempty"`
	Definition
}

// Info returns copy of the machine's entity information by itself.
func (m Machine) Info() MachineInfo {
	return MachineInfo{
		ID:          m.ID,
		Description: m.Description,
		Name:        m.Name,
	}
}

// Summary returns a multiline string summarizing the machine configuration.
func (m Machine) Summary() string {
	var out summarizer
	out.Descend()

	if m.Name != "" {
		out.Add("Name: %s", m.Name)
	}

	if m.Description != "" {
		out.Add("Description: %s", m.Description)
	}

	if !m.ID.IsZero() {
		out.Add("ID: %s", m.ID)
	}

	if len(m.Tags) > 0 {
		out.StartLine()
		for i, tag := range m.Tags {
			if i == 0 {
				out.Printf("Tags: %s", tag)
			} else {
				out.Printf(",%s", tag)
			}
		}
	}

	m.Definition.Config(m.Info(), &out)

	return out.String()
}

// MachineInfo holds identifying information for a machine.
type MachineInfo struct {
	ID          MachineID
	Description string
	Name        MachineName
}

// Seed returns an identity generation seed for the machine info.
func (info MachineInfo) Seed() Seed {
	return Seed{info: info}
}

// Vars returns a set of identifying machine variables. These can be used
// as variables for expansion.
func (info MachineInfo) Vars() Vars {
	return Vars{
		"machine-name":        string(info.Name),
		"machine-description": info.Description,
		"machine-id":          info.ID.String(),
	}
}

// MachineName is the name of a machina virtual machine.
//
// TODO: Document restrictions on machine names and add validity checks.
// These names are used in various places when generating QEMU arguments.
// Spaces in particular could lead to argument parsing badness.
type MachineName string

// MachineID is a universally unique identifer for a machine.
type MachineID UUID

// IsZero returns true if the machine ID holds a zero value.
func (m MachineID) IsZero() bool {
	return m == MachineID{}
}

// String returns a string representation of the machine ID.
func (m MachineID) String() string {
	return UUID(m).String()
}

// MarshalText implements the encoding.TextMarshaler interface.
func (m MachineID) MarshalText() (text []byte, err error) {
	return UUID(m).MarshalText()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (m *MachineID) UnmarshalText(text []byte) error {
	return (*UUID)(m).UnmarshalText(text)
}
