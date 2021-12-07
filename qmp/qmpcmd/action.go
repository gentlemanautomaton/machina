package qmpcmd

// Action is a QMP command with no arguments and no expected response.
type Action string

// Command returns the action as the command name.
func (action Action) Command() string {
	return string(action)
}

// CommandArgs returns a nil JSON byte slice.
func (action Action) CommandArgs() ([]byte, error) {
	return nil, nil
}

// CommandResponse unmarshals a JSON-encoded response to a QMP command.
//
// No response is expected for actions, so this function does nothing.
func (action Action) CommandResponse([]byte) error {
	return nil
}
