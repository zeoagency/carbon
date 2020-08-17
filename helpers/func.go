package helpers

import "strings"

// ReplaceAnyWithSpace replace all given keys with space.
func ReplaceAnyWithSpace(s string, keys ...string) string {
	for _, k := range keys {
		s = strings.ReplaceAll(s, k, " ")
	}
	return strings.TrimSpace(s)
}

// StringSliceContains tells whether a contains x.
func StringSliceContains(a []string, x string) (int, bool) {
	for i, n := range a {
		if x == n {
			return i, true
		}
	}
	return -1, false
}
