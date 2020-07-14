package cbl

import "testing"

func TestB2i(t *testing.T) {
	t.Log(B2i(true), B2i(false))
}

func TestI2b(t *testing.T) {
	t.Log(I2b(1), I2b(0), I2b(-1))
}
