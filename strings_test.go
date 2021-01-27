package cbl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringSplit(t *testing.T) {
	s := ""
	ss := StringSplit(s, ",")
	t.Log(ss) // []

	s = "hello,world"
	ss = StringSplit(s, ",")
	t.Log(ss) // ["hello", "world"]
}

func TestTruncateString(t *testing.T) {
	dataList := [][]interface{}{
		{"hello, world", 5, "hello"},
		{"hello, world", 100, "hello, world"},
		{"为天地立心，为生民立命", 5, "为天地立心"},
	}
	for _, dl := range dataList {
		r := TruncateString(dl[0].(string), dl[1].(int))
		assert.EqualValues(t, r, dl[2])
	}
}
