package devfs

import "strings"

// Path is a devfs path to a device on the host system.
type Path string

// Valid returns true if the path starts with "/dev/".
func (p Path) Valid() bool {
	return strings.HasPrefix(string(p), "/dev/")
}
