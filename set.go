package cbl

// SliceUniq 去重
func SliceUniq(a []string) []string {
	m := make(map[string]bool)
	for _, v := range a {
		m[v] = true
	}

	u := make([]string, 0)
	for k, _ := range m {
		u = append(u, k)
	}
	return u
}

func SliceUniqInt64(a []int64) []int64 {
	m := make(map[int64]bool)
	for _, v := range a {
		m[v] = true
	}

	u := make([]int64, 0)
	for k := range m {
		u = append(u, k)
	}
	return u
}

func SliceUniqInt(a []int) []int {
	m := make(map[int]bool)
	for _, v := range a {
		m[v] = true
	}

	u := make([]int, 0)
	for k := range m {
		u = append(u, k)
	}
	return u
}

// SliceUnion 并集
func SliceUnion(a []string, b []string) []string {
	m := make(map[string]bool)

	c := append(a, b...)
	for _, v := range c {
		m[v] = true
	}

	u := make([]string, 0)
	for k, _ := range m {
		u = append(u, k)
	}

	return u
}

// SliceIntersection 交集
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

// SliceDifference 差集
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
