package sde_test

import (
	"testing"

	"github.com/fffzlfk/sde"
)

func TestSDE(t *testing.T) {
	s, err := sde.NewSDE("string.ns", "string.ns.hs")
	if err != nil {
		t.Fatal(err)
	}
	off, err := s.Encode("hello")
	if err != nil {
		t.Fatal(err)
	}
	str, err := s.Decode(off)
	if err != nil {
		t.Fatal(err)
	}
	if str != "hello" {
		t.Errorf("expected: %s, got: %s\n", "hello", str)
	}
	off, err = s.Encode("world")
	if err != nil {
		t.Fatal(err)
	}
	str, err = s.Decode(off)
	if err != nil {
		t.Fatal(err)
	}
	if str != "world" {
		t.Errorf("expected: %s, got: %s\n", "world", str)
	}
}
