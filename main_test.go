package main

import (
	"context"
	"testing"
)

func TestProccess(t *testing.T) {
	raws := make([]byte, 0, 256)
	for i := 0; i < 256; i++ {
		raws = append(raws, byte(i))
	}
	raws = processText(raws)
	for _, c := range raws {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == 0 {

		} else {
			t.Errorf("Unexpected charater: %d", c)
		}
	}
}

func TestTokenize(t *testing.T) {
	// TODO: More random
	ctx := context.Background()
	raws := processText([]byte("A B CD A-B?C Hello?XXXえぇZ"))
	data := ""
	for t := range tokenize(ctx, raws) {
		if len(data) != 0 {
			data += ","
		}
		data += t.String()
	}
	if data != "a,b,cd,a,b,c,hello,xxx,z" {
		t.Errorf("Unexpected result, %s", data)
	}
}

func TestFrequencyCounter(t *testing.T) {
	// More random
	{
		fc, _ := newFrequencyCounter(1)
		t1 := token("Hello")
		fc.Add(t1)
		fc.Add(t1)
		fc.Add(t1)
		t2 := token("World")
		fc.Add(t2)
		t3 := token("!")

		if n := fc.Get(t1); n.Count != 3 {
			t.Errorf("Unexpected count, t1 got %d", n)
		}
		if n := fc.Get(t2); n.Count != 1 {
			t.Errorf("Unexpected count, t2 got %d", n)
		}
		if n := fc.Get(t3); n.Count != 0 {
			t.Errorf("Unexpected count, t3 got %d", n)
		}
	}

	{
		fc, _ := newFrequencyCounter(1000)
		t1 := token("Hello")
		fc.Add(t1)
		fc.Add(t1)
		fc.Add(t1)
		t2 := token("World")
		fc.Add(t2)
		t3 := token("!")

		if n := fc.Get(t1); n.Count != 3 {
			t.Errorf("Unexpected count, t1 got %d", n)
		}
		if n := fc.Get(t2); n.Count != 1 {
			t.Errorf("Unexpected count, t2 got %d", n)
		}
		if n := fc.Get(t3); n.Count != 0 {
			t.Errorf("Unexpected count, t3 got %d", n)
		}
	}
}
