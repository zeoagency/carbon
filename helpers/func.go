package helpers

import "strings"

// ReplaceAnyWithSpace replace all given keys with space.
func ReplaceAnyWithSpace(s string, keys ...string) string {
	for _, k := range keys {
		s = strings.ReplaceAll(s, k, " ")
	}
	return strings.TrimSpace(s)
}

// StringContains tells whether a contains x.
func StringContains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
