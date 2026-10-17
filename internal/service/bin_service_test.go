package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Germanicus1/fb/api"
)

func newBinsClient(t *testing.T) *api.Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bins" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"results":[{"_id":"binBacklog1","name":"Backlog"},{"_id":"binPeter1","name":"peter.backlog"}]}`))
	}))
	t.Cleanup(server.Close)
	return api.NewClientWithBaseURL("test-key", server.URL)
}

func TestLookupBin_OneWordNameResolvesByName(t *testing.T) {
	bin, err := LookupBin(newBinsClient(t), "backlog")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bin.ID != "binBacklog1" {
		t.Errorf("got id %q, want binBacklog1", bin.ID)
	}
}

func TestLookupBin_UnknownOneWordRefIsNotFound(t *testing.T) {
	_, err := LookupBin(newBinsClient(t), "Foo")
	want := "failed to find bin 'Foo': bin not found: Foo"
	if err == nil || err.Error() != want {
		t.Fatalf("got error %v, want %q", err, want)
	}
}

func TestLookupBin_ResolvesByID(t *testing.T) {
	bin, err := LookupBin(newBinsClient(t), "binPeter1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bin.Name != "peter.backlog" {
		t.Errorf("got name %q, want peter.backlog", bin.Name)
	}
}
