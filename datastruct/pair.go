package datastruct

import (
	"fmt"
	"sort"
	"strings"
)

// Pair a string key/value pair
type Pair struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

type Pairs []Pair

// NewPairsFromMap map {key, value} to Pairs (keep sequence stable)
func NewPairsFromMap(m map[string]string) Pairs {
	pairs := make([]Pair, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, Pair{Key: k, Val: v})
	}
	// sort by key
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Key < pairs[j].Key
	})
	return Pairs(pairs)
}

// ToString Pair to a string, "kvSep" for key/value delimiter; "pSep" for Pair delimiter
// if "kvSep" is "", default set "="
// if "pSep" is "", default set ", "
func (pairs Pairs) ToString(kvSep, pSep string) string {
	if kvSep == "" {
		kvSep = "="
	}
	if pSep == "" {
		pSep = ", "
	}

	ss := make([]string, 0, len(pairs))
	for _, p := range pairs {
		ss = append(ss, fmt.Sprintf("%s%s%s", p.Key, kvSep, p.Val))
	}
	return strings.Join(ss, pSep)
}
