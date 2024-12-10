package systemdgen

import (
	"regexp"
	"strings"
)

// Option is a command option that can be executed in a systemd unit with
// properly quoted arguments.
type Option interface {
	Valid() bool
	OptionPrefix() string
	OptionType() string
	OptionParameters() string
}

// https://www.freedesktop.org/software/systemd/man/systemd.syntax.html#Quoting

// QuoteOptions returns a multiline string for invocation of a command with
// the given options.
//
// The returned string will be properly quoted so that it is suitable for use
// in a systemd exec command line.
func QuoteOptions[T Option, Options ~[]T](opts Options) string {
	var b strings.Builder
	for i, option := range opts {
		last := i == len(opts)-1
		b.WriteString(QuoteOption(option))
		if !last {
			b.WriteString(" \\\n")
		}
	}
	return b.String()
}

// QuoteOption returns a string representation of the option.
//
// It returns an empty string if the option is invalid.
//
// The returned string will be properly quoted so that it is suitable for use
// in a systemd exec command line.
func QuoteOption[T Option](opt T) string {
	if !opt.Valid() {
		return ""
	}

	optionParameters := opt.OptionParameters()
	switch optionParameters {
	case "":
		return QuoteArg(opt.OptionPrefix() + opt.OptionType())
	default:
		return QuoteArg(opt.OptionPrefix()+opt.OptionType()) + " " + QuoteArg(optionParameters)
	}
}

var quoteNeeded = regexp.MustCompile(`[^\w@%+=:,./-]`)

var escapeQuotedChars = strings.NewReplacer(`\`, `\\`, `"`, `\"`, `'`, `\'`)

// QuoteArg returns the argument in a form suitable for use in a systemd
// exec command line. The returned string is only quoted if necessary.
func QuoteArg(arg string) string {
	// Try using no quotes at all
	if !quoteNeeded.MatchString(arg) {
		return arg
	}

	// Quote the arg and escape its contents
	return `"` + escapeQuotedChars.Replace(arg) + `"`
}
