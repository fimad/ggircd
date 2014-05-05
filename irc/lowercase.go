package irc

import (
	"bytes"
	"strings"
)

// lowercase takes a string and returns the IRC lower case version. The IRC spec
// defines {}|^ to be the lower case equivalents of []\~.
func lowercase(in string) string {
	var out bytes.Buffer
	for i := 0; i < len(in); i++ {
		switch in[i] {
		case '[':
			out.WriteRune('{')
		case ']':
			out.WriteRune('}')
		case '\\':
			out.WriteRune('|')
		case '~':
			out.WriteRune('^')
		default:
			out.WriteString(strings.ToLower(string(in[i])))
		}
	}
	return out.String()
}
