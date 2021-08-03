package machina

import (
	"strings"
	"unicode"
)

func cleanInterfaceName(name string) string {
	return strings.Map(func(r rune) rune {
		if r == '/' || unicode.IsSpace(r) || !unicode.IsPrint(r) {
			return -1
		}
		return r
	}, name)
}
