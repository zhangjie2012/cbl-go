package cbl

// B2i bool to int
func B2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

func I2b(i int) bool {
	return i != 0
}
