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

	err := s.Add(testDup...)
	if err != nil {
		t.Fatal(err)
	}

	if len(s.URLs) != len(test) {
		t.Fatal("Error: Set Data Structure issue.")
	}
}

func TestUrlSetAddShouldFail(t *testing.T) {
	s := NewURLSet()

	failTest := []string{"http://nonurl/blablabla"}

	err := s.Add(failTest...)
	if err == nil {
		t.Fatal(err)
	}
}
