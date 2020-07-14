package cbl

import (
	"testing"
)

func TestGenGUID(t *testing.T) {
	t.Log(GenUUIDV1())
	t.Log(GenUUIDV4())
}

func TestGenSessionID(t *testing.T) {
	t.Log(GenUSessionID())
	t.Log(GenRSessionID())
}
