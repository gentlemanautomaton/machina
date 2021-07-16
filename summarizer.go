package machina

import (
	"fmt"
	"strings"
)

// Summary is an interface for summarizing configuration.
type Summary interface {
	Add(format string, a ...interface{})
	Descend()
	Ascend()
}

// summarizer is used to build multiline summaries.
type summarizer struct {
	strings.Builder
	indent  int
	started bool
}

// Descend increases the current indent level.
func (s *summarizer) Descend() {
	s.indent++
}

// Ascend descreases the current indent level. It panics if the indent level
// becomes negative.
func (s *summarizer) Ascend() {
	s.indent--
	if s.indent < 0 {
		panic("mismatched ascent/descent calls in summarizer")
	}
}

// StartLine starts a new line and writes the current indent to the summary.
func (s *summarizer) StartLine() {
	if s.started {
		s.WriteRune('\n')
	} else {
		s.started = true
	}
	for i := 0; i < s.indent; i++ {
		s.WriteString("  ")
	}
}

// Printf prints the given string to the summarizer in the same manner as
// fmt.Printf.
func (s *summarizer) Printf(format string, a ...interface{}) {
	s.WriteString(fmt.Sprintf(format, a...))
}

// Add starts a new line and prints the given string as its value.
func (s *summarizer) Add(format string, a ...interface{}) {
	s.StartLine()
	s.Printf(format, a...)
}
