package datastruct

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPairs(t *testing.T) {
	var (
		m  = map[string]string{"a": "aaa", "b": "bbb", "c": "ccc"}
		s  = ""
		s1 = "a=aaa, b=bbb, c=ccc"
		s2 = "a==aaa; b==bbb; c==ccc"
	)

	pairs := NewPairsFromMap(m)

	s = pairs.ToString("", "")
	assert.Equal(t, s1, s)

	s = pairs.ToString("==", "; ")
	assert.Equal(t, s2, s)
}
