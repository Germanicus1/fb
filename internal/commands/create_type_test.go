package commands

import (
	"strings"
	"testing"
)

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
		f.ConfigTicketType = "Bug"
		if err := f.runCreate(CreateOptions{Title: "T", Bin: "binDoing1"}, "", false, nil).Err; err != nil {
			t.Fatal(err)
		}
		if f.Creates[0].TicketTypeID != "typeBug" {
			t.Errorf("got %+v", f.Creates[0])
		}
	})
	t.Run("--type overrides config", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		f.ConfigTicketType = "Bug"
		if err := f.runCreate(CreateOptions{Title: "T", Bin: "binDoing1", Type: "General"}, "", false, nil).Err; err != nil {
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
