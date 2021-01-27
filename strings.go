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

// TruncateString
// Ref: https://dev.to/takakd/go-safe-truncate-string-9h0
func TruncateString(str string, length int) string {
	if length <= 0 {
		return ""
	}

	truncated := ""
	count := 0
	for _, char := range str {
		truncated += string(char)
		count++
		if count >= length {
			break
		}
	}
	return truncated
}
