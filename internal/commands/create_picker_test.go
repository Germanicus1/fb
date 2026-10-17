package commands

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Germanicus1/fb/internal/state"
	"github.com/Germanicus1/fb/models"
)

// manyBins returns ~1000 bins: peter.backlog, zebra-solo, 3 Team-Alpha-N,
// 25 bulk-N and 970 filler-N.
func manyBins() []models.Bin {
	bins := []models.Bin{{ID: "binBacklog1", Name: "peter.backlog"}, {ID: "binZebra", Name: "zebra-solo"}}
	for i := 1; i <= 3; i++ {
		bins = append(bins, models.Bin{ID: fmt.Sprintf("binAlpha%d", i), Name: fmt.Sprintf("Team-Alpha-%d", i)})
	}
	for i := 1; i <= 25; i++ {
		bins = append(bins, models.Bin{ID: fmt.Sprintf("binBulk%d", i), Name: fmt.Sprintf("bulk-%d", i)})
	}
	for i := 0; i < 970; i++ {
		bins = append(bins, models.Bin{ID: fmt.Sprintf("binFill%d", i), Name: fmt.Sprintf("filler-%d", i)})
	}
	return bins
}

func ctxSetup(t *testing.T) func(string) {
	return func(h string) { writeBinContext(t, h, "binBacklog1", "peter.backlog") }
}

func TestCreateBinPicker(t *testing.T) {
	t.Run("interactive: Enter picks Bin context shown as default", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		f.Bins = manyBins()
		r := f.runCreate(CreateOptions{}, "Title\n\n\n", true, ctxSetup(t))
		if r.Err != nil {
			t.Fatal(r.Err)
		}
		if f.Creates[0].BinID != "binBacklog1" {
			t.Errorf("bin = %q", f.Creates[0].BinID)
		}
		if !strings.Contains(r.Stdout, "Bin [peter.backlog]: ") {
			t.Errorf("prompt: %q", r.Stdout)
		}
		if strings.Index(r.Stdout, "Description") > strings.Index(r.Stdout, "Bin [") {
			t.Errorf("bin prompt must follow description: %q", r.Stdout)
		}
	})

	t.Run("interactive: prompted even with Bin context; typed input overrides", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		f.Bins = manyBins()
		r := f.runCreate(CreateOptions{}, "Title\n\nzebra\n", true, ctxSetup(t))
		if r.Err != nil {
			t.Fatal(r.Err)
		}
		if f.Creates[0].BinID != "binZebra" {
			t.Errorf("bin = %q", f.Creates[0].BinID)
		}
		if !strings.Contains(r.Stdout, "→ zebra-solo") {
			t.Errorf("echo missing: %q", r.Stdout)
		}
	})

	t.Run("interactive with --bin does not prompt for bin", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(CreateOptions{Bin: "peter.backlog"}, "Title\n\n", true, ctxSetup(t))
		if r.Err != nil || strings.Contains(r.Stdout, "Bin [") || strings.Contains(r.Stdout, "Bin: ") {
			t.Errorf("err=%v out=%q", r.Err, r.Stdout)
		}
	})

	t.Run("no context: prompt is plain Bin:, Enter re-prompts", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		f.Bins = manyBins()
		r := f.runCreate(CreateOptions{}, "Title\n\n\nzebra\n", true, nil)
		if r.Err != nil {
			t.Fatal(r.Err)
		}
		if strings.Count(r.Stdout, "Bin: ") != 2 || strings.Contains(r.Stdout, "Bin [") {
			t.Errorf("out: %q", r.Stdout)
		}
	})

	t.Run("2-20 matches list; number picks; invalid number re-prompts", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		f.Bins = manyBins()
		r := f.runCreate(CreateOptions{}, "Title\n\nalpha\n9\nx\n2\n", true, nil)
		if r.Err != nil {
			t.Fatal(r.Err)
		}
		if f.Creates[0].BinID != "binAlpha2" {
			t.Errorf("bin = %q", f.Creates[0].BinID)
		}
		if !strings.Contains(r.Stdout, "1) Team-Alpha-1") || !strings.Contains(r.Stdout, "3) Team-Alpha-3") {
			t.Errorf("list: %q", r.Stdout)
		}
	})

	t.Run(">20 matches and 0 matches re-prompt with messages", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		f.Bins = manyBins()
		r := f.runCreate(CreateOptions{}, "Title\n\nfiller\nnomatchxyz\nzebra\n", true, nil)
		if r.Err != nil {
			t.Fatal(r.Err)
		}
		if !strings.Contains(r.Stdout, "970 bins match") || !strings.Contains(r.Stdout, "type more") {
			t.Errorf("too many: %q", r.Stdout)
		}
		if !strings.Contains(r.Stdout, "No bin matches") {
			t.Errorf("no match: %q", r.Stdout)
		}
		if f.Creates[0].BinID != "binZebra" {
			t.Errorf("bin = %q", f.Creates[0].BinID)
		}
	})

	t.Run("EOF at bin prompt errors, nothing created", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(CreateOptions{}, "Title\n\n", true, nil)
		if r.Err == nil || len(f.Creates) != 0 {
			t.Errorf("err=%v creates=%d", r.Err, len(f.Creates))
		}
	})

	t.Run("one-liner at terminal, no bin, no context: picker", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		f.Bins = manyBins()
		r := f.runCreate(CreateOptions{Title: "T"}, "zebra\n", true, nil)
		if r.Err != nil {
			t.Fatal(r.Err)
		}
		if f.Creates[0].BinID != "binZebra" || !strings.Contains(r.Stdout, "Bin: ") {
			t.Errorf("bin=%q out=%q", f.Creates[0].BinID, r.Stdout)
		}
	})

	t.Run("one-liner at terminal with Bin context does not prompt", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(CreateOptions{Title: "T"}, "", true, ctxSetup(t))
		if r.Err != nil || strings.Contains(r.Stdout, "Bin: ") || strings.Contains(r.Stdout, "Bin [") {
			t.Errorf("err=%v out=%q", r.Err, r.Stdout)
		}
	})

	t.Run("one-liner without terminal and no bin still errors", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(CreateOptions{Title: "T"}, "zebra\n", false, nil)
		if r.Err == nil || !strings.Contains(r.Err.Error(), "--bin") || len(f.Creates) != 0 {
			t.Errorf("err=%v", r.Err)
		}
	})

	t.Run("flags still apply alongside interactive mode (--assign, --checkout)", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		f.Users = map[string]string{"bob@example.com": "user-bob"}
		r := f.runCreate(CreateOptions{Assign: []string{"bob@example.com"}, Checkout: true}, "Title\n\n\n", true, ctxSetup(t))
		if r.Err != nil {
			t.Fatal(r.Err)
		}
		got := f.Creates[0]
		if len(got.AssignedIDs) != 1 || got.AssignedIDs[0] != "user-bob" {
			t.Errorf("assigned = %v", got.AssignedIDs)
		}
		if strings.Contains(r.Stdout, "Assign") {
			t.Errorf("prompted for assignee: %q", r.Stdout)
		}
		if co, err := state.LoadCheckout(); err != nil || co.TicketID != got.ID {
			t.Errorf("checkout not saved: %v", err)
		}
	})
}
