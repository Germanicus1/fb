package cli

import (
	"os"

	"golang.org/x/term"
)

// isTerminal reports whether f is an interactive terminal (not merely a
// character device such as /dev/null).
func isTerminal(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
}
