package cbl

import "testing"

func TestStringSplit(t *testing.T) {
	s := ""
	ss := StringSplit(s, ",")
	t.Log(ss) // []

	s = "hello,world"
	ss = StringSplit(s, ",")
	t.Log(ss) // ["hello", "world"]
}
