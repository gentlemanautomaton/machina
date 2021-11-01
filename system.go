package machina

// System holds configuration for the virtual machine host system.
type System struct {
	// Storage defines storage pools available on the host system.
	Storage StorageMap `json:"storage,omitempty"`

	// Network defines network pools available on the host system.
	Network NetworkMap `json:"network,omitempty"`

	// MediatedDevices is a list of mediated devices available on the host
	// system.
	MediatedDevices MediatedDeviceMap `json:"mdev,omitempty"`

	// Tag defines tags available on the host system.
	Tag TagMap `json:"tag,omitempty"`
}

// Summary returns a multiline string summarizing the system configuration.
func (sys System) Summary() string {
	var out summarizer
	out.Descend()

	if len(sys.Storage) > 0 {
		out.Add("Storage:")
		out.Descend()
		for name, store := range sys.Storage {
			out.Add("%s: %s", name, store.Path)
		}
		out.Ascend()
	}

	if len(sys.Network) > 0 {
		out.Add("Networks:")
		out.Descend()
		for name, network := range sys.Network {
			out.Add("%s: %s", name, network)
		}
		out.Ascend()
	}

	if len(sys.MediatedDevices) > 0 {
		out.Add("Mediated Devices:")
		out.Descend()
		for name, device := range sys.MediatedDevices {
			out.Add("%s:", name)
			out.Descend()
			device.Config(&out)
			out.Ascend()
		}
		out.Ascend()
	}

	if len(sys.Tag) > 0 {
		out.Add("Tags:")
		out.Descend()
		for tag, def := range sys.Tag {
			out.Add("%s:", tag)
			out.Descend()
			def.Config(MachineInfo{}, &out)
			out.Ascend()
		}
		out.Ascend()
	}

	return out.String()
}
