// Package prompt reads lines of user input. At a terminal, lines can be edited
// with the arrow keys and recalled with Up and Down.
package prompt

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

// ErrInterrupted is returned when the user presses Ctrl-C at a prompt.
var ErrInterrupted = errors.New("cancelled")

// LineReader shows a prompt and reads one line, without its line ending.
// At the end of input it returns io.EOF, with any final partial line.
type LineReader interface {
	ReadLine(prompt string) (string, error)
}

// NewLineReader returns a line editor when in is a terminal, else a plain
// reader that writes the prompt to out.
func NewLineReader(in io.Reader, out io.Writer) LineReader {
	if f, ok := in.(*os.File); ok && term.IsTerminal(int(f.Fd())) {
		fd := int(f.Fd())
		r := newTermReader(struct {
			io.Reader
			io.Writer
		}{f, out}, func() (func(), error) {
			old, err := term.MakeRaw(fd)
			if err != nil {
				return nil, err
			}
			return func() { term.Restore(fd, old) }, nil
		})
		if w, _, err := term.GetSize(fd); err == nil && w > 0 {
			r.t.SetSize(w, 0)
		}
		return r
	}
	return &plainReader{in: bufio.NewReader(in), out: out}
}

type plainReader struct {
	in  *bufio.Reader
	out io.Writer
}

func (r *plainReader) ReadLine(prompt string) (string, error) {
	io.WriteString(r.out, prompt)
	s, err := r.in.ReadString('\n')
	return strings.TrimRight(s, "\r\n"), err
}

// termReader edits lines in raw mode, which it enters only while reading so
// other output keeps normal line endings.
type termReader struct {
	t       *term.Terminal
	tty     io.Writer
	keys    *ctrlCWatcher
	makeRaw func() (restore func(), err error)
}

func newTermReader(tty io.ReadWriter, makeRaw func() (func(), error)) *termReader {
	keys := &ctrlCWatcher{r: tty}
	t := term.NewTerminal(struct {
		io.Reader
		io.Writer
	}{keys, tty}, "")
	return &termReader{t: t, tty: tty, keys: keys, makeRaw: makeRaw}
}

func (r *termReader) ReadLine(prompt string) (string, error) {
	restore, err := r.makeRaw()
	if err != nil {
		return "", err
	}
	defer restore()
	r.t.SetPrompt(prompt)
	r.keys.pressed = false
	line, err := r.t.ReadLine()
	if errors.Is(err, term.ErrPasteIndicator) {
		err = nil
	}
	if err == io.EOF && r.keys.pressed {
		io.WriteString(r.tty, "^C\r\n")
		return "", ErrInterrupted
	}
	return line, err
}

// ctrlCWatcher notes Ctrl-C, which term.Terminal reports as io.EOF just like
// Ctrl-D.
type ctrlCWatcher struct {
	r       io.Reader
	pressed bool
}

func (w *ctrlCWatcher) Read(p []byte) (int, error) {
	n, err := w.r.Read(p)
	if bytes.IndexByte(p[:n], 3) >= 0 {
		w.pressed = true
	}
	return n, err
}
