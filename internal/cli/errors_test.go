package cli

import (
	"bytes"
	"errors"
	"testing"
)

func TestPrintErrorWritesError(t *testing.T) {
	var buf bytes.Buffer
	PrintError(&buf, errors.New("failed to find bin 'Foo': bin not found: Foo"))
	want := "Error: failed to find bin 'Foo': bin not found: Foo\n"
	if buf.String() != want {
		t.Errorf("got %q, want %q", buf.String(), want)
	}
}

func TestPrintErrorSkipsFlagErrorsAlreadyPrinted(t *testing.T) {
	var buf bytes.Buffer
	PrintError(&buf, errReported{errors.New("flag provided but not defined: -x")})
	if buf.Len() != 0 {
		t.Errorf("expected no output, got %q", buf.String())
	}
}
