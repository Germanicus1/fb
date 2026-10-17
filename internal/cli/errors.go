package cli

import (
	"errors"
	"fmt"
	"io"
)

// errReported wraps an error that has already been shown to the user, such
// as a flag parse error the flag package printed along with usage.
type errReported struct{ error }

func (e errReported) Unwrap() error { return e.error }

// PrintError writes err to w as "Error: <err>", unless it was already shown.
func PrintError(w io.Writer, err error) {
	var reported errReported
	if errors.As(err, &reported) {
		return
	}
	fmt.Fprintf(w, "Error: %v\n", err)
}
