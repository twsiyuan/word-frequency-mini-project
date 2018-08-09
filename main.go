package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
)

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

func tokenizeFromReader(ctx context.Context, reader io.Reader) chan token {
	// Iterator pattern
	c := make(chan token)
	go func() {
		defer close(c)
		tbuf := make([]byte, 0)
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
						temp := make([]byte, len(tbuf))
						copy(temp, tbuf)
						tbuf = tbuf[:0]
						t = token(temp)
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
				}
			}
		}
	}()

	return c
}

func main() {
	var fileName string
	flag.StringVar(&fileName, "file", "", "process file name")
	flag.Parse()

	if len(fileName) <= 0 {
		fmt.Fprintf(os.Stderr, "Must setup argument 'file'\n")
		os.Exit(1)
	}

	// Read from file
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot open file, %s\n", err.Error())
		os.Exit(1)
	}
	defer file.Close()

	// Tokenize and count (streaming read)
	mmap, err := NewMap(1000)
	if err != nil {
		panic(err)
	}
	for token := range tokenizeFromReader(context.Background(), file) {
		freq, ok := mmap.Get(token).(*frequency)
		if !ok {
			mmap.Set(token, &frequency{
				Token: token,
				Count: 1,
			})
		} else {
			freq.Count++
		}
	}

	// Sort, and output top 20
	values := mmap.Values(make([]interface{}, 0))
	freqs := make([]*frequency, len(values))
	for i := 0; i < len(values); i++ {
		freqs[i] = values[i].(*frequency)
	}

	sort.Sort(frequencySorter(freqs))
	for i := 0; i < len(freqs) && i < 20; i++ {
		freq := freqs[i]
		fmt.Fprintf(os.Stdout, "%d %s\n", freq.Count, freq.Token)
	}
}
