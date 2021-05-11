package cbl

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestDurationString(t *testing.T) {
	loc := time.Now().Location()
	current := time.Date(2021, 5, 11, 17, 50, 30, 0, loc)

	{
		p := time.Date(2019, 2, 11, 0, 0, 0, 0, loc)
		d := current.Sub(p)
		assert.Equal(t, "2 年", DurationString(d))
	}

	{
		p := time.Date(2020, 4, 11, 0, 0, 0, 0, loc)
		d := current.Sub(p)
		assert.Equal(t, "13 月", DurationString(d))
	}

	{
		p := time.Date(2021, 5, 05, 0, 0, 0, 0, loc)
		d := current.Sub(p)
		assert.Equal(t, "6 天", DurationString(d))
	}

	{
		p := time.Date(2021, 5, 11, 0, 0, 0, 0, loc)
		d := current.Sub(p)
		assert.Equal(t, "17 小时", DurationString(d))
	}

	{
		p := time.Date(2021, 5, 11, 17, 0, 0, 0, loc)
		d := current.Sub(p)
		assert.Equal(t, "50 分钟", DurationString(d))
	}

	{
		p := time.Date(2021, 5, 11, 17, 50, 0, 0, loc)
		d := current.Sub(p)
		assert.Equal(t, "30 秒", DurationString(d))
	}

	{

		p := time.Date(2021, 5, 11, 17, 50, 30, 0, loc)
		d := current.Sub(p)
		assert.Equal(t, "0 秒", DurationString(d))
	}
}
