package commands

import (
	"strings"
	"testing"
)

func TestCreateAssign(t *testing.T) {
	t.Run("--assign replaces me and does not double-encode @", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		f.Users = map[string]string{"me@example.com": "user-me", "a@x.com": "uA", "b@x.com": "uB"}
		r := f.runCreate(CreateOptions{Title: "T", Bin: "peter.backlog", Assign: []string{"a@x.com", "b@x.com"}}, "", false, nil)
		if r.Err != nil {
			t.Fatal(r.Err)
		}
		if got := f.Creates[0].AssignedIDs; len(got) != 2 || got[0] != "uA" || got[1] != "uB" {
			t.Fatalf("assigned_ids = %v", got)
		}
		for _, p := range f.UserPaths {
			if strings.Contains(p, "%40") {
				t.Fatalf("double-encoded @: %s", p)
			}
		}
	})

	t.Run("without --assign assigns me", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		r := f.runCreate(CreateOptions{Title: "T", Bin: "peter.backlog"}, "", false, nil)
		if r.Err != nil {
			t.Fatal(r.Err)
		}
		if got := f.Creates[0].AssignedIDs; len(got) != 1 || got[0] != "user-me" {
			t.Fatalf("assigned_ids = %v", got)
		}
	})

	t.Run("unknown emails all listed in one error and nothing created", func(t *testing.T) {
		f := newFakeFlowBoards(t)
		f.Users = map[string]string{"me@example.com": "user-me", "a@x.com": "uA"}
		r := f.runCreate(CreateOptions{Title: "T", Bin: "peter.backlog", Assign: []string{"nope1@x.com", "a@x.com", "nope2@x.com"}}, "", false, nil)
		if r.Err == nil {
			t.Fatal("expected error")
		}
		for _, e := range []string{"nope1@x.com", "nope2@x.com"} {
			if !strings.Contains(r.Err.Error(), e) {
				t.Fatalf("error %q missing %s", r.Err, e)
			}
		}
		if len(f.Creates) != 0 {
			t.Fatal("create was sent")
		}
	})
}
