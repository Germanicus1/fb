package commands

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/Germanicus1/fb/api"
	"github.com/Germanicus1/fb/models"
)

// fakeFlowBoards is the shared fake server for command tests of `fb create`.
// It serves users, bins, ticket types and the create endpoint, and records
// every POST body. Extend it by adding fields/routes; do not fork it.
type fakeFlowBoards struct {
	t      *testing.T
	Server *httptest.Server

	UserID      string
	Bins        []models.Bin
	TicketTypes []models.TicketType

	mu      sync.Mutex
	Creates []models.CreateTicketPayload
	// LastBody is the raw JSON of the latest create POST.
	LastBody string
}

func newFakeFlowBoards(t *testing.T) *fakeFlowBoards {
	t.Helper()
	f := &fakeFlowBoards{
		t:      t,
		UserID: "user-me",
		Bins: []models.Bin{
			{ID: "binBacklog1", Name: "peter.backlog"},
			{ID: "binDoing1", Name: "Doing"},
		},
		TicketTypes: []models.TicketType{
			{ID: "typeGeneral", Name: "General"},
			{ID: "typeBug", Name: "Bug"},
		},
	}
	f.Server = httptest.NewServer(http.HandlerFunc(f.handle))
	t.Cleanup(f.Server.Close)
	return f
}

func (f *fakeFlowBoards) handle(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/users/"):
		json.NewEncoder(w).Encode(models.User{ID: f.UserID, Email: strings.TrimPrefix(r.URL.Path, "/users/")})
	case r.Method == "GET" && r.URL.Path == "/bins":
		json.NewEncoder(w).Encode(f.Bins)
	case r.Method == "GET" && r.URL.Path == "/ticket-types":
		json.NewEncoder(w).Encode(f.TicketTypes)
	case r.Method == "POST" && strings.HasPrefix(r.URL.Path, "/tickets/"):
		body, _ := io.ReadAll(r.Body)
		var p models.CreateTicketPayload
		if err := json.Unmarshal(body, &p); err != nil {
			http.Error(w, "bad json", 400)
			return
		}
		f.mu.Lock()
		f.Creates = append(f.Creates, p)
		f.LastBody = string(body)
		f.mu.Unlock()
		w.WriteHeader(204)
	default:
		http.NotFound(w, r)
	}
}

// createRun is the result of invoking the create command through its seam.
type createRun struct {
	Err    error
	Stdout string
	Home   string
}

// runCreate invokes the create command with HOME redirected to a temp dir.
// setup, if non-nil, may seed files under Home (e.g. bin_context.json).
func (f *fakeFlowBoards) runCreate(opts CreateOptions, stdin string, isTerminal bool, setup func(home string)) createRun {
	f.t.Helper()
	home := f.t.TempDir()
	f.t.Setenv("HOME", home)
	if setup != nil {
		setup(home)
	}
	var out bytes.Buffer
	err := ExecuteCreate(CreateDeps{
		Client:     api.NewClientWithBaseURL("test-key", f.Server.URL),
		UserEmail:  "me@example.com",
		Stdin:      strings.NewReader(stdin),
		Stdout:     &out,
		IsTerminal: isTerminal,
	}, opts)
	return createRun{Err: err, Stdout: out.String(), Home: home}
}

func writeBinContext(t *testing.T, home, id, name string) {
	t.Helper()
	dir := filepath.Join(home, ".fb")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	data, _ := json.Marshal(map[string]string{"bin_id": id, "bin_name": name})
	if err := os.WriteFile(filepath.Join(dir, "bin_context.json"), data, 0600); err != nil {
		t.Fatal(err)
	}
}
