package commands

import (
	"strings"
	"testing"
)

func TestCreateDescription(t *testing.T) {
	run := func(t *testing.T, desc, stdin string) *fakeFlowBoards {
		f := newFakeFlowBoards(t)
		r := f.runCreate(CreateOptions{Title: "T", Bin: "peter.backlog", Description: desc}, stdin, false, nil)
		if r.Err != nil {
			t.Fatalf("unexpected error: %v", r.Err)
		}
		if len(f.Creates) != 1 {
			t.Fatalf("want 1 create, got %d", len(f.Creates))
		}
		return f
	}

	t.Run("-d text sent verbatim", func(t *testing.T) {
		f := run(t, `a\nb`, "")
		if got := f.Creates[0].Description; got != `a\nb` {
			t.Errorf("description = %q", got)
		}
	})

	t.Run("-d - reads all stdin verbatim", func(t *testing.T) {
		f := run(t, "-", "line1\n  line2\n\nline4\n")
		if got := f.Creates[0].Description; got != "line1\n  line2\n\nline4\n" {
			t.Errorf("description = %q", got)
		}
	})

	t.Run("-d - with empty stdin sends no description", func(t *testing.T) {
		f := run(t, "-", "")
		if strings.Contains(f.LastBody, `"description"`) {
			t.Errorf("body has description key: %s", f.LastBody)
		}
	})

	t.Run("no -d sends no description key", func(t *testing.T) {
		f := run(t, "", "ignored stdin")
		if strings.Contains(f.LastBody, `"description"`) {
			t.Errorf("body has description key: %s", f.LastBody)
		}
	})
}
