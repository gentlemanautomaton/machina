package systemdgen

import (
	"regexp"
	"strings"

	"github.com/gentlemanautomaton/machina/commandoption"
)

// https://www.freedesktop.org/software/systemd/man/systemd.syntax.html#Quoting

// QuoteOptions returns a multiline string for invocation of a command with
// the given options.
//
// The returned string will be properly quoted so that it is suitable for use
// in a systemd exec command line.
func QuoteOptions[T commandoption.Option, Options commandoption.Options[T]](opts Options) string {
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
func QuoteOption[T commandoption.Option](opt T) string {
	data := commandoption.Data(opt)
	if data.Type == "" {
		return ""
	}

	switch params := data.Parameters.String(); params {
	case "":
		return QuoteArg(opt.Prefix() + data.Type)
	default:
		return QuoteArg(opt.Prefix()+data.Type) + " " + QuoteArg(params)
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
