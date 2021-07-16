package machina

// NetworkName identifies a network on the local system by a well-known name.
//
// TODO: Document restrictions on network names and add validity checks.
// These names are used in various places when generating QEMU arguments.
// Spaces in particular could lead to argument parsing badness.
type NetworkName string

// NetworkMap maps network names to networks on the local system.
type NetworkMap map[NetworkName]Network

// Network defines a network that a machine can be connected to.
type Network struct {
	Device string `json:"device"`
	Up     string `json:"up"`
	Down   string `json:"down"`
}

// String returns a string representation of the network configuration.
func (n Network) String() string {
	return n.Device
}
