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

func TestEndOfDay(t *testing.T) {
	c := time.Now()
	c1 := EndOfDay(c)
	t.Log(c, c1)
}
