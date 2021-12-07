package qmpmsg

import (
	"encoding/json"
	"time"
)

// Event holds an asynchronous event.
type Event struct {
	Event     string    `json:"event"`
	Data      Data      `json:"data,omitempty"`
	Timestamp Timestamp `json:"timestamp"`
}

// Timestamp is a time value that can be umarshaled from a QMP timestamp
// value encoded in a JSON object.
type Timestamp time.Time

type timestamp struct {
	Seconds      int64 `json:"seconds"`
	Microseconds int64 `json:"microseconds"`
}

// UnmarshalJSON unmarshals a QMP timestamp from the given JSON object data.
func (m *Timestamp) UnmarshalJSON(b []byte) error {
	var ts timestamp
	if err := json.Unmarshal(b, &ts); err != nil {
		return err
	}
	*m = Timestamp(time.Unix(ts.Seconds, ts.Microseconds*100))
	return nil
}
