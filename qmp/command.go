package qmp

import (
	"encoding/json"
	"fmt"

	"github.com/gentlemanautomaton/machina/qmp/qmpmsg"
)

// Command is a QMP command that can be sent by a client.
type Command interface {
	// Command returns the name of the command to be issued.
	Command() string

	// CommandArgs returns the arguments of the command, marshaeled as JSON.
	// If no arguments are provided, a nil byte slice will be returned.
	CommandArgs() ([]byte, error)

	// CommandResponse unmarshals the response to the command, if present.
	CommandResponse([]byte) error
}

// command is used for marshaling QMP commands without arguments.
type command struct {
	Execute string    `json:"execute"`
	ID      qmpmsg.ID `json:"id,omitempty"`
}

// commandWithArgs is used for marshaling QMP commands with arguments.
type commandWithArgs struct {
	Execute string          `json:"execute"`
	Args    json.RawMessage `json:"arguments,omitempty"`
	ID      qmpmsg.ID       `json:"id,omitempty"`
}

// response is used for unmarshaling QMP command responses.
type response struct {
	Data  responseData `json:"return,omitempty"`
	Error qmpmsg.Error `json:"error,omitempty"`
	ID    qmpmsg.ID    `json:"id,omitempty"`
}

// responseData captures QMP command response data and hands it a
// CommandResponse handler provided by the Command interface. This allows
// commands to
type responseData struct {
	Command Command
}

func (r responseData) UnmarshalJSON(b []byte) error {
	if r.Command == nil {
		return nil
	}
	return r.Command.CommandResponse(b)
}

func marshalAndValidateArgs(cmd Command) (args []byte, err error) {
	args, err = cmd.CommandArgs()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments for %s: %w", cmd.Command(), err)
	}
	if len(args) > 0 && !json.Valid(args) {
		return nil, fmt.Errorf("failed to marshal arguments for %s: invalid command argument data", cmd.Command())
	}
	return
}
