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

func TestFormatDayTime(t *testing.T) {
	loc := time.Now().Location()
	for i := 0; i < 24; i++ {
		c := time.Date(2020, 9, 6, i, i+1, i+2, 0, loc)
		s := FormatDayTime(c, i%2 == 0)
		t.Log(c, s)
	}
}

func TestDayOfWeekCN(t *testing.T) {
	for i := 0; i < 7; i++ {
		t.Log(i, "->", DayOfWeekCN(i))
	}
}
