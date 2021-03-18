package cbl

import (
	"fmt"
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

// Int64SliceToString a int64 slice to string
// Ref: https://stackoverflow.com/questions/37532255/one-liner-to-transform-int-into-string
func Int64SliceToString(a []int64, delim string) string {
	return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), delim), "[]")
}

// IntSliceToString a int64 slice to string
// Ref: https://stackoverflow.com/questions/37532255/one-liner-to-transform-int-into-string
func IntSliceToString(a []int, delim string) string {
	return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), delim), "[]")
}
