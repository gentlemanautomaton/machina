package qmpmsg

// Error holds a QMP error returned by a command.
type Error struct {
	Class       string `json:"class"`
	Description string `json:"desc"`
}

// IsZero returns true if the error is empty.
func (e Error) IsZero() bool {
	return e.Class == "" && e.Description == ""
}

// Error returns the QMP error as a string, so that it can be used as a Go
// error type.
func (e Error) Error() string {
	switch {
	case e.Class != "" && e.Description != "":
		return e.Class + " error: " + e.Description
	case e.Class != "":
		return e.Class
	case e.Description != "":
		return e.Description
	default:
		return "unspecified qmp error"
	}
}
