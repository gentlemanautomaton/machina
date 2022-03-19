package machina

import (
	"fmt"
	"os"
	"strconv"
)

// PatternMapper is a function that can map variables to values.
type PatternMapper func(string) string

// StringPattern is a pattern that can undergo variable expansion to produce
// string values.
type StringPattern string

// Expand returns the expanded string for the given mapper.
func (pattern StringPattern) Expand(mapper PatternMapper) string {
	return os.Expand(string(pattern), mapper)
}

// IntPattern is a pattern that can undergo variable expansion to produce
// integer values.
type IntPattern string

// Expand returns the expanded integer for the given mapper. If the expanded
// value cannot be converted to an integer, an error is returned.
func (pattern IntPattern) Expand(mapper PatternMapper) (int, error) {
	s := os.Expand(string(pattern), mapper)
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("integer pattern \"%s\": %w", pattern, err)
	}
	return value, nil
}

// PortPattern is a pattern that can undergo variable expansion to produce
// network port numbers.
type PortPattern string

// Expand returns the expanded integer for the given mapper. If the expanded
// value cannot be converted to an integer, an error is returned.
func (pattern PortPattern) Expand(mapper PatternMapper) (int, error) {
	const min, max = 1, 65535

	s := os.Expand(string(pattern), mapper)

	value, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("port pattern \"%s\": %w", pattern, err)
	}
	if value < min || value > max {
		return 0, fmt.Errorf("port pattern \"%s\": invalid resulting port number %d", pattern, value)
	}

	return value, nil
}
