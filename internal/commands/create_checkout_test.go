package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Germanicus1/fb/internal/state"
)

func writeCheckoutFile(t *testing.T, home, id, name string) {
	t.Helper()
	dir := filepath.Join(home, ".fb")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	data := `{"ticket_id":"` + id + `","ticket_name":"` + name + `","bin_id":"b","bin_name":"b"}`
	if err := os.WriteFile(filepath.Join(dir, "checkout.json"), []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestCreateCheckout(t *testing.T) {
	t.Run("--checkout saves Checkout and Bin context", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(CreateOptions{Title: "Fix", Bin: "peter.backlog", Checkout: true}, "", false, nil)
		if r.Err != nil {
			t.Fatal(r.Err)
		}
		co, err := state.LoadCheckout()
		if err != nil || co.TicketID != f.Creates[0].ID || co.TicketName != "Fix" || co.BinID != "binBacklog1" {
			t.Fatalf("checkout = %+v err %v", co, err)
		}
		bc, err := state.LoadBinContext()
		if err != nil || bc.BinID != "binBacklog1" || bc.BinName != "peter.backlog" {
			t.Fatalf("bin context = %+v err %v", bc, err)
		}
	})

	t.Run("existing Checkout without --force refuses before create", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(CreateOptions{Title: "Fix", Bin: "peter.backlog", Checkout: true}, "", false,
			func(h string) { writeCheckoutFile(t, h, "old", "Old ticket") })
		if r.Err == nil || !strings.Contains(r.Err.Error(), "ticket already checked out: Old ticket\nUse 'fb clear' or 'fb checkout --force'") {
			t.Fatalf("err = %v", r.Err)
		}
		if len(f.Creates) != 0 {
			t.Errorf("create sent")
		}
	})

	t.Run("--force replaces Checkout", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(CreateOptions{Title: "Fix", Bin: "peter.backlog", Checkout: true, Force: true}, "", false,
			func(h string) { writeCheckoutFile(t, h, "old", "Old ticket") })
		if r.Err != nil {
			t.Fatal(r.Err)
		}
		co, _ := state.LoadCheckout()
		if co == nil || co.TicketID != f.Creates[0].ID {
			t.Fatalf("checkout = %+v", co)
		}
	})

	t.Run("plain create leaves Bin context and Checkout alone", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(CreateOptions{Title: "Fix", Bin: "peter.backlog"}, "", false,
			func(h string) { writeBinContext(t, h, "binDoing1", "Doing") })
		if r.Err != nil {
			t.Fatal(r.Err)
		}
		bc, _ := state.LoadBinContext()
		if bc.BinID != "binDoing1" {
			t.Errorf("bin context changed: %+v", bc)
		}
		if _, err := state.LoadCheckout(); err == nil {
			t.Errorf("checkout written")
		}
	})

	t.Run("failure saving Checkout still prints success line then errors", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(CreateOptions{Title: "Fix", Bin: "peter.backlog", Checkout: true}, "", false, func(h string) {
			// a directory where checkout.json should go makes the save fail
			if err := os.MkdirAll(filepath.Join(h, ".fb", "checkout.json"), 0700); err != nil {
				t.Fatal(err)
			}
		})
		if r.Err == nil {
			t.Fatal("want error")
		}
		if !strings.HasPrefix(r.Stdout, "✓ Created "+f.Creates[0].ID+" - Fix  [peter.backlog]\n") {
			t.Errorf("stdout = %q", r.Stdout)
		}
	})
}
