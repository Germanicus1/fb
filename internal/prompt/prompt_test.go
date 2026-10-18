package prompt

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestPlainReaderEchoesPromptAndReadsLines(t *testing.T) {
	var out bytes.Buffer
	r := NewLineReader(strings.NewReader("first\r\nlast"), &out)

	line, err := r.ReadLine("Title: ")
	if line != "first" || err != nil {
		t.Fatalf("got %q, %v; want \"first\", nil", line, err)
	}
	line, err = r.ReadLine("Bin: ")
	if line != "last" || err != io.EOF {
		t.Fatalf("got %q, %v; want \"last\", io.EOF", line, err)
	}
	if out.String() != "Title: Bin: " {
		t.Errorf("output %q, want %q", out.String(), "Title: Bin: ")
	}
}

// fakeTTY is keyboard input and the screen of a terminal in raw mode.
type fakeTTY struct {
	io.Reader
	io.Writer
}

func newTestTermReader(keys string) (*termReader, *int) {
	rawCalls := 0
	r := newTermReader(&fakeTTY{strings.NewReader(keys), io.Discard}, func() (func(), error) {
		rawCalls++
		return func() {}, nil
	})
	return r, &rawCalls
}

func TestTermReaderEditsWithArrowKeys(t *testing.T) {
	// type "ab", Left twice, type "x", Enter
	r, rawCalls := newTestTermReader("ab\x1b[D\x1b[Dx\r")

	line, err := r.ReadLine("Title: ")
	if line != "xab" || err != nil {
		t.Fatalf("got %q, %v; want \"xab\", nil", line, err)
	}
	if *rawCalls != 1 {
		t.Errorf("raw mode entered %d times, want 1", *rawCalls)
	}
}

func TestTermReaderKeepsTypeAheadForNextLine(t *testing.T) {
	r, _ := newTestTermReader("one\rtwo\r")

	for _, want := range []string{"one", "two"} {
		if line, err := r.ReadLine(""); line != want || err != nil {
			t.Fatalf("got %q, %v; want %q, nil", line, err, want)
		}
	}
}

func TestTermReaderCtrlCInterrupts(t *testing.T) {
	r, _ := newTestTermReader("ab\x03")

	if _, err := r.ReadLine("Title: "); !errors.Is(err, ErrInterrupted) {
		t.Fatalf("err %v, want ErrInterrupted", err)
	}
}

func TestTermReaderCtrlDOnEmptyLineIsEOF(t *testing.T) {
	r, _ := newTestTermReader("\x04")

	if _, err := r.ReadLine(""); err != io.EOF {
		t.Fatalf("err %v, want io.EOF", err)
	}
}
