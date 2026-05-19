package string_tools

import (
	"strings"
)

func SplitString(s string, sep string, values ...*string) {
	words := strings.SplitN(s, sep, len(values))
	for i := range len(words) {
		*values[i] = words[i]
	}
}
