package machina

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

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

// maxIfaceLength is the maximum length of an interface name in linux.
// It excludes the null terminator.
const maxIfaceLength = 15

// MakeLinkName returns the network interface name for a connection.
func MakeLinkName(machine MachineName, conn Connection) string {
	name := fmt.Sprintf("%s.%s", machine, conn.Name)

	// If the link name can be used as-is, just do that
	if len(name) < maxIfaceLength {
		return cleanInterfaceName(name)
	}

	// If the link name is too long, use a hash of it
	return hashLinkName(name)
}

// hashLinkName returns a 15 character base64-encoded hash of name.
func hashLinkName(name string) string {
	hasher := sha256.New()
	hasher.Write([]byte(name))
	hashed := hasher.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(hashed[0:11])
}
