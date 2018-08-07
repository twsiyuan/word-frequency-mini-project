package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
)

type token []byte

func (t token) Equal(v token) bool {
	if t == nil && v == nil {
		return true
	} else if t == nil && v != nil {
		return false
	} else if t != nil && v == nil {
		return false
	} else if len(t) != len(v) {
		return false
	}

	for i := 0; i < len(t); i++ {
		if t[i] != v[i] {
			return false
		}
	}
	return true
}

func (t token) String() string {
	return string(t)
}

func (t token) Hash() int {
	if t == nil {
		return 0
	}
	hash := 0
	for i, b := range t {
		hash += i * int(b)
	}
	return hash
}

func (t token) Less(v token) bool {
	for i := 0; i < len(t); i++ {
		if i >= len(v) {
			return true
		}
		if t[i] != v[i] {
			return t[i] < v[i]
		}
	}
	return true
}

func (t token) Clone() token {
	b := make([]byte, len(t))
	copy(b, t)
	return token(b)
}

func processText(raws []byte) []byte {
	// ASCII only
	for idx, b := range raws {
		raws[idx] = processTextByte(b)
	}
	return raws
}

func processTextByte(b byte) byte {
	// ASCII only
	if b >= 'a' && b <= 'z' {
		return b
	} else if b >= 'A' && b <= 'Z' {
		return b + 32
	} else {
		return 0
	}
}

func tokenizeReader(ctx context.Context, reader io.Reader) chan token {
	// Iterator pattern
	c := make(chan token)
	go func() {
		defer close(c)
		tbuf := make([]byte, 0)
		tbuf2 := make([]byte, 0)
		rbuf := make([]byte, 1024)
		for true {
			n, err := reader.Read(rbuf)
			if err == io.EOF {
				if len(tbuf) > 0 {
					c <- token(tbuf)
				}
				return
			} else if err != nil {
				panic(err)
			}

			for i := 0; i < n; i++ {
				t := token(nil)
				b := processTextByte(rbuf[i])

				if b == 0 {
					if len(tbuf) > 0 {
						t = token(tbuf)
					}
				} else {
					tbuf = append(tbuf, b)
				}

				if t != nil {
					select {
					case c <- t:
						break
					case <-ctx.Done():
						return
					}
					// Goroutine need switch buffer to avoid error (did not clone buffer)
					tbuf, tbuf2 = tbuf2, tbuf
					tbuf = tbuf[:0]
				}
			}
		}
	}()

	return c
}

func tokenize(ctx context.Context, raws []byte) chan token {
	return tokenizeReader(ctx, bytes.NewReader(raws))
}

type frequency struct {
	Token token
	Count int
}

type frequencySorter []frequency

func (fs frequencySorter) Len() int {
	return len(fs)
}

func (fs frequencySorter) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

func (fs frequencySorter) Less(i, j int) bool {
	vi := fs[i]
	vj := fs[j]
	if vi.Count == vj.Count {
		return !vi.Token.Less(vj.Token)
	}
	return fs[i].Count > fs[j].Count
}

func newFrequencyCounter(c int) (*frequencyCounter, error) {
	if c <= 0 {
		return nil, errors.New("Bad Cap")
	}
	fc := frequencyCounter{
		cap: make([][]frequency, c),
	}
	return &fc, nil
}

type frequencyCounter struct {
	cap   [][]frequency
	count int
}

func (fc frequencyCounter) Count() int {
	return fc.count
}

func (fc *frequencyCounter) Get(t token) *frequency {
	hash := t.Hash()
	idx := hash % len(fc.cap)
	data := fc.cap[idx]
	if data != nil {
		for i := 0; i < len(data); i++ {
			f := data[i]
			if t.Equal(f.Token) {
				return &f
			}
		}
	}
	return &frequency{
		Token: t,
		Count: 0,
	}
}

func (fc *frequencyCounter) Add(t token) {
	hash := t.Hash()
	idx := hash % len(fc.cap)
	data := fc.cap[idx]
	if data == nil {
		data = make([]frequency, 0)
		fc.cap[idx] = data
	}

	// Try to find a hit
	for i := 0; i < len(data); i++ {
		f := &data[i]
		if t.Equal(f.Token) {
			f.Count++
			return
		}
	}

	// No hit, add one
	fc.cap[idx] = append(data, frequency{
		Token: t.Clone(),
		Count: 1,
	})
	fc.count++
}

func (fc frequencyCounter) List() []frequency {
	v := make([]frequency, 0, fc.count)
	for i := 0; i < len(fc.cap); i++ {
		v = append(v, fc.cap[i]...)
	}
	return v
}

func main() {
	var file string
	flag.StringVar(&file, "file", "", "process file name")
	flag.Parse()

	if len(file) <= 0 {
		fmt.Fprintf(os.Stderr, "Must setup argument 'file'\n")
		os.Exit(1)
	}

	f, err := os.Open(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot open file, %s\n", err.Error())
		os.Exit(1)
	}

	// Tokenize and count (streaming read)
	fc, _ := newFrequencyCounter(1000)
	for t := range tokenizeReader(context.Background(), f) {
		fc.Add(t)
	}

	// Sort, and output top 20
	s := fc.List()
	sort.Sort(frequencySorter(s))
	for i := 0; i < len(s) && i < 20; i++ {
		f := s[i]
		fmt.Fprintf(os.Stdout, "%d %s\n", f.Count, f.Token)
	}
}
