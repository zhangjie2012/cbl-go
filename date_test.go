package cbl

import (
	"testing"
	"time"
)

func TestStartOfDay(t *testing.T) {
	c := time.Now()
	c1 := StartOfDay(c)
	t.Log(c, c1)
}
