package models

import (
	"testing"
)

func TestKeywordSetAdd(t *testing.T) {
	k := NewKeywordSet()

	test := []string{
		"boratanrikulu blog yazıları",
	}

	testDup := append(test, test...)
	k.Add(testDup...)

	if len(k.Keywords) != len(test) {
		t.Fatal("Error: Set Data Structure issue.")
	}
}
