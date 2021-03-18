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

func TestInt64SliceToString(t *testing.T) {
	a := []int64{1, 2, 3, 4, 5, 6}
	b := Int64SliceToString(a, ", ")
	assert.Equal(t, "1, 2, 3, 4, 5, 6", b)
}

func TestIntSliceToString(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6}
	b := IntSliceToString(a, ", ")
	assert.Equal(t, "1, 2, 3, 4, 5, 6", b)
}
