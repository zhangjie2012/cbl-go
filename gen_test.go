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

func TestGenVerifyCode(t *testing.T) {
	t.Log(GenVerifyCodeNumber(8))
	t.Log(GenVerifyCodeCommon(8))

	t.Log(GenVerifyCodeAny("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 8))
	t.Log(GenVerifyCodeAny("!@#$%^&*()", 8))
}
