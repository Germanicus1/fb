package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/Germanicus1/fb/models"
)

func TestNewClientWithBaseURL_SkipsDiscovery(t *testing.T) {
	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		w.Write([]byte(`[]`))
	}))
	defer server.Close()

	client := NewClientWithBaseURL("key", server.URL)
	if _, err := client.ListTicketTypes(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paths) != 1 || paths[0] != "/ticket-types" {
		t.Errorf("expected single request to /ticket-types, got %v", paths)
	}
}

func TestCreateTicket_SendsDocumentedRequestAndAccepts204(t *testing.T) {
	var method, path string
	var body map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		json.Unmarshal(raw, &body)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClientWithBaseURL("key", server.URL)
	err := client.CreateTicket(models.CreateTicketPayload{
		ID: "abc123", Name: "Fix login", TicketTypeID: "tt1", BinID: "bin1",
		Description: "Repro", AssignedIDs: []string{"u1", "u2"},
	})
	if err != nil {
		t.Fatalf("expected success on 204, got %v", err)
	}
	if method != "POST" || path != "/tickets/abc123" {
		t.Errorf("got %s %s", method, path)
	}
	want := `{"_id":"abc123","assigned_ids":["u1","u2"],"bin_id":"bin1","description":"Repro","name":"Fix login","ticketType_id":"tt1"}`
	got, _ := json.Marshal(body) // map marshals with sorted keys
	if string(got) != want {
		t.Errorf("body = %s, want %s", got, want)
	}
}

func TestCreateTicket_OmitsEmptyDescriptionAndSendsEmptyAssignees(t *testing.T) {
	var raw string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		raw = string(b)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClientWithBaseURL("key", server.URL)
	if err := client.CreateTicket(models.CreateTicketPayload{ID: "x", Name: "n", TicketTypeID: "t", BinID: "b"}); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(raw, "description") {
		t.Errorf("description should be omitted: %s", raw)
	}
	if !strings.Contains(raw, `"assigned_ids":[]`) {
		t.Errorf("assigned_ids must always be sent as array: %s", raw)
	}
}

func TestCreateTicket_ErrorIncludesStatusAndBody(t *testing.T) {
	cases := []struct {
		status int
		body   string
	}{
		{400, `{"name":[["is_required"]]}`},
		{403, `Access denied`},
	}
	for _, c := range cases {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(c.status)
			w.Write([]byte(c.body))
		}))
		err := NewClientWithBaseURL("key", server.URL).CreateTicket(models.CreateTicketPayload{ID: "x"})
		server.Close()
		if err == nil {
			t.Fatalf("status %d: expected error", c.status)
		}
		if !strings.Contains(err.Error(), strconv.Itoa(c.status)) || !strings.Contains(err.Error(), c.body) {
			t.Errorf("status %d: error should include status and body, got: %v", c.status, err)
		}
	}
}

func TestListTicketTypes_ParsesDisabledFlag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"_id":"a","name":"General"},{"_id":"b","name":"Old","disabled":true}]`))
	}))
	defer server.Close()

	types, err := NewClientWithBaseURL("key", server.URL).ListTicketTypes()
	if err != nil {
		t.Fatal(err)
	}
	if len(types) != 2 || types[0].ID != "a" || types[0].Name != "General" || types[0].Disabled ||
		types[1].ID != "b" || !types[1].Disabled {
		t.Errorf("unexpected types: %+v", types)
	}
}

func TestGetTicket_FetchesById(t *testing.T) {
	var path string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		w.Write([]byte(`{"_id":"abc","name":"Fix","bin_id":"b1","assigned_ids":["u1"]}`))
	}))
	defer server.Close()

	ticket, err := NewClientWithBaseURL("key", server.URL).GetTicket("abc")
	if err != nil {
		t.Fatal(err)
	}
	if path != "/tickets/abc" || ticket.ID != "abc" || ticket.Name != "Fix" || ticket.BinID != "b1" {
		t.Errorf("path=%s ticket=%+v", path, ticket)
	}
}

func TestGenerateID_17Alphanumeric(t *testing.T) {
	re := regexp.MustCompile(`^[A-Za-z0-9]{17}$`)
	seen := map[string]bool{}
	for i := 0; i < 200; i++ {
		id := GenerateID()
		if !re.MatchString(id) {
			t.Fatalf("bad id %q", id)
		}
		seen[id] = true
	}
	if len(seen) < 200 {
		t.Error("ids should be unique")
	}
}
