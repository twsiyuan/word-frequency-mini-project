package main

import (
	"testing"
)

func TestMap(t *testing.T) {
	m, err := NewMap(1)
	if err != nil {
		t.Fatal(err)
	}

	if m.Count() != 0 {
		t.Error("Unexpected map count")
	}

	t1 := token([]byte{1})
	t2 := token([]byte{1})
	t3 := token([]byte{1, 1})

	v1 := "Hello World"
	v2 := "Hello Sky"

	if m.Get(t1) != nil {
		t.Error("Unexpected map value")
	}

	m.Set(t1, v1)
	if m.Get(t1).(string) != v1 {
		t.Error("Unexpected map value")
	}

	m.Set(t2, v2)
	if m.Get(t2).(string) != v2 {
		t.Error("Unexpected map value")
	} else if m.Get(t1).(string) != v2 {
		t.Error("Unexpected map value")
	}

	m.Set(t3, v1)
	if m.Count() != 2 {
		t.Error("Unexpected map count")
	}
}
