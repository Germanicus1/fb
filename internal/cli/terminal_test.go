package cli

import (
	"os"
	"testing"
)

func TestIsTerminalDevNullIsNotTerminal(t *testing.T) {
	f, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if isTerminal(f) {
		t.Error("/dev/null reported as a terminal")
	}
}
