package commands

import (
	"strings"
	"testing"
)

func TestCreateBinSingleWordName(t *testing.T) {
	for _, ref := range []string{"Doing", "doing"} {
		f := newFakeFlowBoards(t)
		r := f.runCreate(CreateOptions{Title: "T", Bin: ref}, "", false, nil)
		if r.Err != nil {
			t.Fatalf("%s: unexpected error: %v", ref, r.Err)
		}
		if f.Creates[0].BinID != "binDoing1" || !strings.HasSuffix(r.Stdout, "  [Doing]\n") {
			t.Errorf("%s: payload %+v stdout %q", ref, f.Creates[0], r.Stdout)
		}
	}
}

func TestCreateBinUnknown(t *testing.T) {
	t.Run("unknown name fails before create", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(CreateOptions{Title: "T", Bin: "no such.bin"}, "", false, nil)
		if r.Err == nil || r.Err.Error() != "failed to find bin 'no such.bin': bin not found: no such.bin" || len(f.Creates) != 0 {
			t.Fatalf("err=%v creates=%d", r.Err, len(f.Creates))
		}
	})
	t.Run("unknown id-like ref is sent and 403 names it", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		f.CreateFault = "403"
		r := f.runCreate(CreateOptions{Title: "T", Bin: "zzUnknown9"}, "", false, nil)
		if r.Err == nil || r.Err.Error() != "bin not found or no access: zzUnknown9" {
			t.Fatalf("err=%v", r.Err)
		}
	})
}

func TestCreate403ShowsBinName(t *testing.T) {
	t.Run("--bin by id", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		f.CreateFault = "403"
		r := f.runCreate(CreateOptions{Title: "T", Bin: "binDoing1"}, "", false, nil)
		if r.Err == nil || r.Err.Error() != "bin not found or no access: Doing" {
			t.Fatalf("err=%v", r.Err)
		}
	})
	t.Run("Bin context", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		f.CreateFault = "403"
		r := f.runCreate(CreateOptions{Title: "T"}, "", false, func(home string) {
			writeBinContext(t, home, "binDoing1", "Doing")
		})
		if r.Err == nil || r.Err.Error() != "bin not found or no access: Doing" {
			t.Fatalf("err=%v", r.Err)
		}
	})
}
