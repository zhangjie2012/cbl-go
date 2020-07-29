package datasize

import (
	"testing"
)

func TestB(t *testing.T) {
	n := 1000000000000000
	bs := New(n)

	t.Log(bs.KB())
	t.Log(bs.MB())
	t.Log(bs.GB())
	t.Log(bs.TB())
	t.Log(bs.PB())
	t.Log(bs.EB())

	t.Log(bs.KiB())
	t.Log(bs.MiB())
	t.Log(bs.GiB())
	t.Log(bs.TiB())
	t.Log(bs.PiB())
	t.Log(bs.EiB())
}

func TestFormat2String(t *testing.T) {
	t.Log(float2String(12000))
	t.Log(float2String(12000.120))
	t.Log(float2String(12000.120304))
}

func TestString(t *testing.T) {
	{
		n := 1024
		bs := New(n)
		t.Log(bs.String1())
		t.Log(bs.String2())
	}
	{
		n := 2000
		bs := New(n)
		t.Log(bs.String1())
		t.Log(bs.String2())
	}
}
