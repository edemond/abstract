package util

import (
	"fmt"
	"strings"
)

// strings.Join, but for ints
func JoinInts(ints []int, sep string) string {
	strs := make([]string, len(ints))
	for i, s := range ints {
		strs[i] = fmt.Sprintf("%v", s)
	}
	return strings.Join(strs, sep)
}
