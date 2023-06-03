package summary

import (
	"fmt"
	"strings"
)

// Builder is used to build multiline summaries.
type Builder struct {
	strings.Builder
	indent  int
	started bool
}

// Descend increases the current indent level.
func (s *Builder) Descend() {
	s.indent++
}

// Ascend descreases the current indent level. It panics if the indent level
// becomes negative.
func (s *Builder) Ascend() {
	s.indent--
	if s.indent < 0 {
		panic("mismatched ascent/descent calls in summary.Builder")
	}
}

// StartLine starts a new line and writes the current indent to the summary.
func (s *Builder) StartLine() {
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
func (s *Builder) Printf(format string, a ...interface{}) {
	s.WriteString(fmt.Sprintf(format, a...))
}

// Add starts a new line and prints the given string as its value.
func (s *Builder) Add(format string, a ...interface{}) {
	s.StartLine()
	s.Printf(format, a...)
}
