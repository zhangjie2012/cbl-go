package cbl

func SliceIntersection(a []string, b []string) []string {
	hash := make(map[string]bool)
	for _, v := range a {
		hash[v] = true
	}

	inter := make([]string, 0)
	for _, v := range b {
		if hash[v] {
			inter = append(inter, v)
		}
	}

	return inter
}

func SliceDifference(a []string, b []string) []string {
	hash := make(map[string]bool)
	for _, v := range b {
		hash[v] = true
	}

	diff := make([]string, 0)
	for _, v := range a {
		if !hash[v] {
			diff = append(diff, v)
		}
	}

	return diff
}
