package qmpmsg

import "encoding/json"

// Type identifies a QMP message type.
type Type int

// Message types.
const (
	UnknownMessage = 0
	ReponseMessage = 1
	EventMessage   = 2
)

// detector determines which fields are present in a QMP response message.
type detector struct {
	Event placeholder `json:"event"`
	Data  placeholder `json:"return"`
	Error placeholder `json:"error"`
}

// placeholder identifies whether a JSON field is present.
type placeholder bool

func (p placeholder) Present() bool {
	return bool(p)
}

func (p *placeholder) UnmarshalJSON(b []byte) error {
	*p = true
	return nil
}

// DetectType attempts to determine the type of message contained in msg.
//
// If the message is malformed or the type cannot be determined, it returns
// UnknownMessage.
func DetectType(msg []byte) Type {
	var d detector
	if err := json.Unmarshal(msg, &d); err != nil {
		return UnknownMessage
	}
	switch {
	case d.Event.Present():
		return EventMessage
	case d.Data.Present(), d.Error.Present():
		return ReponseMessage
	default:
		return UnknownMessage
	}
}
