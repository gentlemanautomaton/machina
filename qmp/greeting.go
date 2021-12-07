package qmp

import (
	"encoding/json"
	"fmt"
)

// Greeting is sent by the server when a connection is first established.
type Greeting struct {
	QMP GreetingPayload `json:"QMP"`
}

// UnmarshalQMP unmarshals the given QMP greeting, encoded in JSON.
func (g *Greeting) UnmarshalQMP(msg []byte) error {
	if err := json.Unmarshal(msg, g); err != nil {
		return fmt.Errorf("failed to parse QMP greeting: %w", err)
	}
	return nil
}

// GreetingPayload is the payload of a greeting message.
type GreetingPayload struct {
	Version      VersionInfo  `json:"version"`
	Capabilities Capabilities `json:"capabilities"`
}

// VersionInfo identifies the version of the server that's running.
type VersionInfo struct {
	Qemu    VersionTriple `json:"qemu"`
	Package string        `json:"package"`
}

// String returns a string representation of the version.
func (v VersionInfo) String() string {
	return fmt.Sprintf("qemu %s / %s", v.Qemu, v.Package)
}

// VersionTriple specifies the major, minor and macro parts of a version
// identifier.
type VersionTriple struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Micro int `json:"micro"`
}

// String returns a string representation of the version triple.
func (v VersionTriple) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Micro)
}

// Capability identifies a capability supported by the server.
type Capability string

// Capabilities is a set of supported QMP capabilities.
type Capabilities []string
