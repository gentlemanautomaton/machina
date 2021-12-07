package qmpcmd

import "encoding/json"

// Capabilities is a QMP command that declares QMP capabilities requested by
// the client.
type Capabilities struct {
	// Enable holds the set of QMP capabilities requested by a client.
	Enable []string `json:"enable,omitempty"`
}

// Command returns the command name "qmp_capabilities".
func (caps Capabilities) Command() string {
	return "qmp_capabilities"
}

// CommandArgs returns the QMP capabilities command arguments marshaled as
// a JSON byte slice.
func (caps Capabilities) CommandArgs() ([]byte, error) {
	return json.Marshal(caps)
}

// CommandResponse unmarshals a JSON-encoded response to a QMP capabilities
// command.
//
// No response is expected, so this function does nothing.
func (caps Capabilities) CommandResponse([]byte) error {
	return nil
}
