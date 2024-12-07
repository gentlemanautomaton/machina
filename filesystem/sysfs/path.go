package sysfs

import "strings"

// Path is a sysfs path to a device on the host system.
type Path string

// Valid returns true if the path starts with "/sys/".
func (p Path) Valid() bool {
	return strings.HasPrefix(string(p), "/sys/")
}
