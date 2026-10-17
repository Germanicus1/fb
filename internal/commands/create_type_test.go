package commands

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Germanicus1/fb/api"
)

func (f *fakeFlowBoards) runCreateWithConfigType(opts CreateOptions, cfgType string) error {
	f.t.Helper()
	f.t.Setenv("HOME", f.t.TempDir())
	var out bytes.Buffer
	return ExecuteCreate(CreateDeps{
		Client:           api.NewClientWithBaseURL("test-key", f.Server.URL),
		UserEmail:        "me@example.com",
		Stdin:            strings.NewReader(""),
		Stdout:           &out,
		ConfigTicketType: cfgType,
	}, opts)
}

func TestCreateTicketType(t *testing.T) {
	t.Run("--type by name case-insensitive", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(CreateOptions{Title: "T", Bin: "binDoing1", Type: "bug"}, "", false, nil)
		if r.Err != nil || f.Creates[0].TicketTypeID != "typeBug" {
			t.Fatalf("err=%v creates=%+v", r.Err, f.Creates)
		}
	})
	t.Run("--type by id", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(CreateOptions{Title: "T", Bin: "binDoing1", Type: "typeBug"}, "", false, nil)
		if r.Err != nil || f.Creates[0].TicketTypeID != "typeBug" {
			t.Fatalf("err=%v creates=%+v", r.Err, f.Creates)
		}
	})
	t.Run("config ticket_type used without flag", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		if err := f.runCreateWithConfigType(CreateOptions{Title: "T", Bin: "binDoing1"}, "Bug"); err != nil {
			t.Fatal(err)
		}
		if f.Creates[0].TicketTypeID != "typeBug" {
			t.Errorf("got %+v", f.Creates[0])
		}
	})
	t.Run("--type overrides config", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		if err := f.runCreateWithConfigType(CreateOptions{Title: "T", Bin: "binDoing1", Type: "General"}, "Bug"); err != nil {
			t.Fatal(err)
		}
		if f.Creates[0].TicketTypeID != "typeGeneral" {
			t.Errorf("got %+v", f.Creates[0])
		}
	})
	t.Run("default is General", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(CreateOptions{Title: "T", Bin: "binDoing1"}, "", false, nil)
		if r.Err != nil || f.Creates[0].TicketTypeID != "typeGeneral" {
			t.Fatalf("err=%v", r.Err)
		}
	})
	t.Run("unknown type fails before create", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(CreateOptions{Title: "T", Bin: "binDoing1", Type: "nope"}, "", false, nil)
		if r.Err == nil || !strings.Contains(r.Err.Error(), "not found") || len(f.Creates) != 0 {
			t.Fatalf("err=%v creates=%d", r.Err, len(f.Creates))
		}
	})
	t.Run("disabled type fails before create", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		f.TicketTypes[1].Disabled = true
		r := f.runCreate(CreateOptions{Title: "T", Bin: "binDoing1", Type: "Bug"}, "", false, nil)
		if r.Err == nil || r.Err.Error() != "ticket type disabled: Bug" || len(f.Creates) != 0 {
			t.Fatalf("err=%v creates=%d", r.Err, len(f.Creates))
		}
	})
}
