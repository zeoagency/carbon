package models

import (
	"testing"
)

func TestUrlSetAdd(t *testing.T) {
	s := NewURLSet()

	test := []string{
		"https://boratanrikulu.dev/postgresql-nedir-nasil-calisir/",
		"https://boratanrikulu.dev/smtp-nasil-calisir-ve-postfix-kurulumu/",
		"https://tools.zeo.org/carbon",
	}

	testDup := append(test, test...)

	s.Add(testDup...)
	if len(s.URLs) != len(test) {
		t.Fatal("Error: Set Data Structure issue.")
	}
}

func TestUrlSetAddShouldFail(t *testing.T) {
	s := NewURLSet()

	failTest := []string{
		"http://nonurl/blablabla",
		"aaaaa",
	}

	s.Add(failTest...)
	if len(s.Fails) != len(failTest) {
		t.Fatal("Error: Fail list issue.")
	}
}
