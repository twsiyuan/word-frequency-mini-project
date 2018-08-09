package main

import (
	"testing"
)

func TestToken(t *testing.T) {
	t1 := token([]byte{1})
	t2 := token([]byte{1})
	t3 := token([]byte{1, 2})
	t4 := token([]byte{1, 0})
	t5 := token([]byte{0, 0})

	if !t1.Equal(t2) || t1.Hash() != t2.Hash() {
		t.Errorf("Unexpected result of t1 & t2")
	}

	if t1.Equal(t3) || t1.Hash() == t3.Hash() {
		t.Errorf("Unexpected result of t1 & t3")
	}

	if t1.Equal(t4) || t1.Hash() != t4.Hash() {
		t.Errorf("Unexpected result of t1 & t4")
	}

	if t4.Equal(t5) || t4.Hash() == t5.Hash() {
		t.Errorf("Unexpected result of t4 & t5")
	}
}
