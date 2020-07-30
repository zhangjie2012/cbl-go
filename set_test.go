package cbl

import "testing"

func TestSliceIntersection(t *testing.T) {
	a := []string{"a", "b", "c", "d"}
	b := []string{"b", "c", "d", "e"}
	inter := SliceIntersection(a, b)
	t.Log(inter)
}

func TestSliceDifference(t *testing.T) {
	a := []string{"a", "b", "c", "d", "f"}
	b := []string{"b", "c", "d", "e"}
	diff := SliceDifference(a, b)
	t.Log(diff) // expect ["a", "f"]
}
