package commands

import (
	"strings"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Run("--bin by name posts the right payload and prints confirmation", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(CreateOptions{Title: "Fix login", Bin: "peter.backlog"}, "", false, nil)
		if r.Err != nil {
			t.Fatalf("unexpected error: %v", r.Err)
		}
		if len(f.Creates) != 1 {
			t.Fatalf("want 1 create, got %d", len(f.Creates))
		}
		p := f.Creates[0]
		if p.Name != "Fix login" || p.BinID != "binBacklog1" || p.TicketTypeID != "typeGeneral" ||
			len(p.AssignedIDs) != 1 || p.AssignedIDs[0] != "user-me" || len(p.ID) != 17 {
			t.Errorf("bad payload: %+v", p)
		}
		want := "✓ Created " + p.ID + " - Fix login  [peter.backlog]\n"
		if r.Stdout != want {
			t.Errorf("stdout = %q, want %q", r.Stdout, want)
		}
	})

	t.Run("--bin by id", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(CreateOptions{Title: "T", Bin: "binDoing1"}, "", false, nil)
		if r.Err != nil {
			t.Fatalf("unexpected error: %v", r.Err)
		}
		if f.Creates[0].BinID != "binDoing1" || !strings.HasSuffix(r.Stdout, "  [Doing]\n") {
			t.Errorf("payload %+v stdout %q", f.Creates[0], r.Stdout)
		}
	})

	t.Run("no --bin uses Bin context", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(CreateOptions{Title: "T"}, "", false, func(home string) {
			writeBinContext(t, home, "binDoing1", "Doing")
		})
		if r.Err != nil {
			t.Fatalf("unexpected error: %v", r.Err)
		}
		if f.Creates[0].BinID != "binDoing1" {
			t.Errorf("bin = %s", f.Creates[0].BinID)
		}
	})

	t.Run("no --bin and no Bin context errors naming --bin, no create", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(CreateOptions{Title: "T"}, "", false, nil)
		if r.Err == nil || !strings.Contains(r.Err.Error(), "--bin") {
			t.Errorf("err = %v", r.Err)
		}
		if len(f.Creates) != 0 {
			t.Error("create was sent")
		}
	})

	t.Run("empty or whitespace title errors before any API call", func(t *testing.T) {
		for _, title := range []string{"", "   \t"} {
			f := newFakeFlowBoards(t)
			f.Server.Close() // any API call would fail with a different error
			r := f.runCreate(CreateOptions{Title: title, Bin: "Doing"}, "", false, nil)
			if r.Err == nil || !strings.Contains(r.Err.Error(), "title") {
				t.Errorf("title %q: err = %v", title, r.Err)
			}
		}
	})

	t.Run("403 on create is rewritten to bin not found or no access", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		f.CreateFault = "403"
		r := f.runCreate(CreateOptions{Title: "T", Bin: "binDoing1"}, "", false, nil)
		if r.Err == nil || r.Err.Error() != "bin not found or no access: Doing" {
			t.Errorf("err = %v", r.Err)
		}
	})

	t.Run("network error but ticket exists is success", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		f.CreateFault, f.CreateFaultPersists = "drop", true
		r := f.runCreate(CreateOptions{Title: "T", Bin: "peter.backlog"}, "", false, nil)
		if r.Err != nil {
			t.Fatalf("unexpected error: %v", r.Err)
		}
		if !strings.HasPrefix(r.Stdout, "✓ Created ") {
			t.Errorf("stdout = %q", r.Stdout)
		}
	})

	t.Run("network error and ticket absent returns original error", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		f.CreateFault = "drop"
		r := f.runCreate(CreateOptions{Title: "T", Bin: "peter.backlog"}, "", false, nil)
		if r.Err == nil || !strings.Contains(r.Err.Error(), "request failed") {
			t.Errorf("err = %v", r.Err)
		}
		if r.Stdout != "" {
			t.Errorf("stdout = %q", r.Stdout)
		}
	})

	t.Run("-v prints the created ticket after the success line", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(CreateOptions{Title: "Verbose me", Bin: "peter.backlog", Verbose: true}, "", false, nil)
		if r.Err != nil {
			t.Fatalf("unexpected error: %v", r.Err)
		}
		id := f.Creates[0].ID
		if !strings.HasPrefix(r.Stdout, "✓ Created "+id) || !strings.Contains(r.Stdout, "Ticket ID: "+id) ||
			!strings.Contains(r.Stdout, "Ticket Name: Verbose me") {
			t.Errorf("stdout = %q", r.Stdout)
		}
	})
}
