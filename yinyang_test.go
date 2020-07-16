package cbl

import (
	"testing"
)

func TestToString(t *testing.T) {
	y, _ := ConvYangYin(2020, 6, 25)
	t.Log(y.ToString1())
	t.Log(y.ToString2())
	t.Log(y.ToString3())
}

func TestConvYinYang(t *testing.T) {
	y, err := ConvYinYang(2020, 5, 0, 5)
	t.Log(y, err) // 2020-06-25 <nil>
}

func TestConvYangYin(t *testing.T) {
	y, err := ConvYangYin(2020, 6, 25)
	t.Log(y.ToString1(), err)
}
