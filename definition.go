package machina

// Definition holds the definition of a machine tag.
type Definition struct {
	Vars        Vars         `json:"vars,omitempty"`
	Attributes  Attributes   `json:"attrs,omitempty"`
	Volumes     []Volume     `json:"volumes,omitempty"`
	Connections []Connection `json:"connections,omitempty"`
	Devices     []Device     `json:"devices,omitempty"`
}

// Config adds the attributes configuration to the summary.
func (d *Definition) Config(info MachineInfo, out Summary) {
	if len(d.Vars) > 0 {
		out.Add("Vars:")
		out.Descend()
		for key, value := range d.Vars {
			out.Add("%s: %s", key, value)
		}
		out.Ascend()
	}

	d.Attributes.Config(info, d.Vars, out)

	if len(d.Volumes) > 0 {
		out.Add("Volumes:")
		out.Descend()
		for i := range d.Volumes {
			d.Volumes[i].Config(out)
		}
		out.Ascend()
	}

	if len(d.Connections) > 0 {
		out.Add("Connections:")
		out.Descend()
		for i := range d.Connections {
			d.Connections[i].Config(out)
		}
		out.Ascend()
	}

	if len(d.Devices) > 0 {
		out.Add("Devices:")
		out.Descend()
		for i := range d.Devices {
			d.Devices[i].Config(out)
		}
		out.Ascend()
	}
}

// MergeDefinitions merges a set of definitions in order. If more than one
// volume exists with the same name, only the first is included.
func MergeDefinitions(defs ...Definition) Definition {
	var (
		vars  []Vars
		attrs []Attributes
		vols  []Volume
		conns []Connection
		devs  []Device
	)

	for i := range defs {
		vars = append(vars, defs[i].Vars)
		attrs = append(attrs, defs[i].Attributes)
		vols = append(vols, defs[i].Volumes...)
		conns = append(conns, defs[i].Connections...)
		devs = append(devs, defs[i].Devices...)
	}

	return Definition{
		Vars:        MergeVars(vars...),
		Attributes:  MergeAttributes(attrs...),
		Volumes:     MergeVolumes(vols...),
		Connections: MergeConnections(conns...),
		Devices:     MergeDevices(devs...),
	}
}
