package main

import (
	"bytes"
	"context"
	"testing"
)

func TestTokenize(t *testing.T) {
	// TODO: More random
	ctx := context.Background()
	raws := []byte("A B CD A-B?C Hello?XXXえぇZ")
	data := ""
	for t := range tokenizeFromReader(ctx, bytes.NewReader(raws)) {
		if len(data) != 0 {
			data += ","
		}
		data += t.String()
	}
	if data != "a,b,cd,a,b,c,hello,xxx,z" {
		t.Errorf("Unexpected result, %s", data)
	}
}
