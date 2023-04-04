package qmpcmd

// QueryCPU is a QMP command that returns information about the CPU in a
// virtual machine.
type QueryCPU struct {
	Response string
}

// Command returns the QMP command name.
func (action QueryCPU) Command() string {
	return "query-cpus-fast"
}

// CommandArgs returns a nil JSON byte slice.
func (action QueryCPU) CommandArgs() ([]byte, error) {
	return nil, nil
}

// CommandResponse unmarshals the JSON-encoded response to a QMP command.
func (action *QueryCPU) CommandResponse(response []byte) error {
	action.Response = string(response)
	return nil
}
