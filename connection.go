package machina

import "fmt"

// ConnectionName is the name of a network connection on a machine.
type ConnectionName string

// Connection describes a network connection.
type Connection struct {
	Name    ConnectionName `json:"name"`
	Network NetworkName    `json:"network"`
	IP      string         `json:"ip"`
	MAC     string         `json:"mac"`
}

// String returns a string representation of the network connection
// configuration.
func (c Connection) String() string {
	return fmt.Sprintf("%s: %s (ip: %s, mac: %s)", c.Name, c.Network, c.IP, c.MAC)
}

// Config adds the connection configuration to the summary.
func (c *Connection) Config(out Summary) {
	out.Add("%s", c)
}

// MergeConnections merges a set of connections in order. If more than one
// connection exists with the same name, only the first is included.
func MergeConnections(conns ...Connection) []Connection {
	lookup := make(map[ConnectionName]bool)
	out := make([]Connection, 0, len(conns))
	for _, conn := range conns {
		if seen := lookup[conn.Name]; seen {
			continue
		}
		lookup[conn.Name] = true
		out = append(out, conn)
	}
	return out
}

// LinkName returns the network interface name for a connection.
func LinkName(machine MachineName, conn Connection) string {
	return fmt.Sprintf("%s-%s", machine, conn.Name)
}
