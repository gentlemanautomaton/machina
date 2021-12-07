package qmp

// Message holds a QMP message.
type MessageHandler struct {
	Handler func(msg []byte) error
}

// UnmarshalJSON copies the given JSON data to the respones.
func (m *MessageHandler) UnmarshalJSON(b []byte) error {
	return m.Handler(b)
}
