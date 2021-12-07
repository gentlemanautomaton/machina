package qmpmsg

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
)

// msgIDRegexp matches machina message identifiers encoded as strings in JSON.
var msgIDRegexp = regexp.MustCompile(`^machina\.([0-9]+)\.([0-9]+)`)

// ID uniquely identifies a machina QMP request by combining a random 64-bit
// client identifier with a message sequence number.
type ID struct {
	Client   uint64
	Sequence uint64
}

// IsZero returns true if the ID holds a zero value.
func (id ID) IsZero() bool {
	return id.Client == 0 && id.Sequence == 0
}

// String returns a string representation of the ID.
func (id ID) String() string {
	return fmt.Sprintf("machina.%d.%d", id.Client, id.Sequence)
}

// MarshalText marshals the ID as text.
func (id ID) MarshalText() (text []byte, err error) {
	return []byte(id.String()), nil
}

// UnmarshalJSON attempts to interpret the given JSON data as a machina
// message ID encoded in a JSON string.
//
// The function always returns a nil error, even if the data does not contain
// a machina message ID. This is done in case qemu broadcasts responses
// meant for other clients that use a different message ID format.
func (id *ID) UnmarshalJSON(data []byte) error {
	// Interpret the data as a string
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return nil
	}

	// Verify the message format and capture the client and sequence
	parts := msgIDRegexp.FindStringSubmatch(s)
	if len(parts) != 3 {
		return nil
	}

	// Parse the client as an integer
	client, clientErr := strconv.ParseUint(parts[1], 10, 64)
	if clientErr != nil {
		return nil
	}

	// Parse the sequence as an integer
	seq, seqErr := strconv.ParseUint(parts[2], 10, 64)
	if seqErr != nil {
		return nil
	}

	id.Client = client
	id.Sequence = seq

	return nil
}
