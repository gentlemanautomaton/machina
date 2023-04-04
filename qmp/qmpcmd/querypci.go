package qmpcmd

// QueryCPU is a QMP command that returns information about the PCI bus in a
// virtual machine.
type QueryPCI struct {
	Response string
}

// Command returns the QMP command name.
func (action QueryPCI) Command() string {
	return "query-pci"
}

// CommandArgs returns a nil JSON byte slice.
func (action QueryPCI) CommandArgs() ([]byte, error) {
	return nil, nil
}

// CommandResponse unmarshals the JSON-encoded response to a QMP command.
func (action *QueryPCI) CommandResponse(response []byte) error {
	action.Response = string(response)
	return nil
}
