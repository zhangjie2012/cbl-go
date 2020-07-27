package cbl

import (
	"strings"
)

// StringSplit string split by `sep`
// the standard library `Split` when `source` is a empty string, will return a `[""]` slice.
// In fact, most of time, we want return a length equal 0 slice, `[]`
func StringSplit(source string, sep string) []string {
	source2 := strings.TrimSpace(source)
	if source2 == "" {
		return make([]string, 0)
	}
	return strings.Split(source2, sep)
}
