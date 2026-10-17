package commands

import (
	"strings"
	"testing"
)

func TestCreateInteractive(t *testing.T) {
	bin := CreateOptions{Bin: "peter.backlog"}

	t.Run("prompts title then description, multi-line", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(bin, "Title\nline1\nline2\n\n", true, nil)
		if r.Err != nil {
			t.Fatal(r.Err)
		}
		if len(f.Creates) != 1 {
			t.Fatalf("creates = %d", len(f.Creates))
		}
		if got := f.Creates[0]; got.Name != "Title" || got.Description != "line1\nline2" {
			t.Errorf("got name=%q desc=%q", got.Name, got.Description)
		}
		ti := strings.Index(r.Stdout, "Title: ")
		di := strings.Index(r.Stdout, "Description (optional; finish with an empty line):")
		if ti < 0 || di < ti {
			t.Errorf("prompt order wrong: %q", r.Stdout)
		}
		if !strings.Contains(r.Stdout, "✓ Created ") || !strings.Contains(r.Stdout, "- Title  [peter.backlog]") {
			t.Errorf("confirmation missing: %q", r.Stdout)
		}
	})

	t.Run("empty title re-prompts", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(bin, "\n   \nReal\n\n", true, nil)
		if r.Err != nil {
			t.Fatal(r.Err)
		}
		if len(f.Creates) != 1 || f.Creates[0].Name != "Real" {
			t.Errorf("creates = %+v", f.Creates)
		}
		if n := strings.Count(r.Stdout, "Title: "); n != 3 {
			t.Errorf("title prompts = %d", n)
		}
	})

	t.Run("empty first description line means none", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(bin, "Title\n\n", true, nil)
		if r.Err != nil {
			t.Fatal(r.Err)
		}
		if strings.Contains(f.LastBody, `"description"`) {
			t.Errorf("body has description: %s", f.LastBody)
		}
	})

	t.Run("description ended by EOF is accepted", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(bin, "Title\nonly line", true, nil)
		if r.Err != nil {
			t.Fatal(r.Err)
		}
		if f.Creates[0].Description != "only line" {
			t.Errorf("desc = %q", f.Creates[0].Description)
		}
	})

	t.Run("EOF before a title fails, nothing created", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(bin, "", true, nil)
		if r.Err == nil || !strings.Contains(r.Err.Error(), "title required") {
			t.Errorf("err = %v", r.Err)
		}
		if len(f.Creates) != 0 {
			t.Error("created")
		}
	})

	t.Run("no args and no terminal: title required, nothing read or sent", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(bin, "Title\nbody\n\n", false, nil)
		if r.Err == nil || r.Err.Error() != "title required" {
			t.Errorf("err = %v", r.Err)
		}
		if r.Stdout != "" || len(f.Creates) != 0 {
			t.Errorf("stdout=%q creates=%d", r.Stdout, len(f.Creates))
		}
	})

	t.Run("no prompts when a title is given", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		o := bin
		o.Title = "T"
		r := f.runCreate(o, "", true, nil)
		if r.Err != nil {
			t.Fatal(r.Err)
		}
		if strings.Contains(r.Stdout, "Title: ") || strings.Contains(r.Stdout, "Description") {
			t.Errorf("prompted: %q", r.Stdout)
		}
	})
}
