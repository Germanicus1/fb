package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetBins tests Story 1: Fetch bins from API with names
//
// User Story:
// As a user, I want the system to fetch all bins from the Flow Boards API
// so that I can reference bins by their friendly names instead of IDs.
//
// Acceptance Criteria:
// - System calls GET /bins endpoint
// - Response contains array of bins with ID and Name fields
// - Bin IDs are correctly mapped to bin names
// - API errors are properly handled
// - Empty bin list is handled gracefully
func TestGetBins(t *testing.T) {
	t.Run("Given API endpoint When fetching bins Then return list of bins with IDs and names", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify the request is for the bins endpoint
			if r.URL.Path != "/bins" {
				t.Errorf("Expected path /bins, got %s", r.URL.Path)
			}
			if r.Method != "GET" {
				t.Errorf("Expected GET method, got %s", r.Method)
			}

			// Return mock bins data matching API spec
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{"_id": "cx7oRn0CK1SoAMn0x", "name": "K+Dev.Doing", "color": "#FF5733"},
				{"_id": "bin123", "name": "In Progress", "color": "#00FF00"},
				{"_id": "bin456", "name": "Done", "color": "#0000FF"}
			]`))
		}))
		defer server.Close()

		client := NewClient("test-key")
		client.baseURL = server.URL

		// Act
		bins, err := client.GetBins()

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(bins) != 3 {
			t.Fatalf("Expected 3 bins, got %d", len(bins))
		}

		// Verify first bin
		if bins[0].ID != "cx7oRn0CK1SoAMn0x" {
			t.Errorf("Expected bin ID cx7oRn0CK1SoAMn0x, got %s", bins[0].ID)
		}
		if bins[0].Name != "K+Dev.Doing" {
			t.Errorf("Expected bin name K+Dev.Doing, got %s", bins[0].Name)
		}

		// Verify second bin
		if bins[1].ID != "bin123" {
			t.Errorf("Expected bin ID bin123, got %s", bins[1].ID)
		}
		if bins[1].Name != "In Progress" {
			t.Errorf("Expected bin name In Progress, got %s", bins[1].Name)
		}
	})

	t.Run("Given API returns empty array When fetching bins Then return empty list without error", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("[]"))
		}))
		defer server.Close()

		client := NewClient("test-key")
		client.baseURL = server.URL

		// Act
		bins, err := client.GetBins()

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(bins) != 0 {
			t.Errorf("Expected empty bin list, got %d bins", len(bins))
		}
	})

	t.Run("Given API returns error When fetching bins Then return error", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}))
		defer server.Close()

		client := NewClient("test-key")
		client.baseURL = server.URL

		// Act
		bins, err := client.GetBins()

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if bins != nil {
			t.Errorf("Expected nil bins on error, got %v", bins)
		}
	})

	t.Run("Given base URL not discovered When fetching bins Then return error", func(t *testing.T) {
		// Arrange
		client := NewClient("test-key")
		// Don't call DiscoverRestPrefix, so baseURL is empty

		// Act
		bins, err := client.GetBins()

		// Assert
		if err == nil {
			t.Fatal("Expected error for missing base URL, got nil")
		}
		if bins != nil {
			t.Errorf("Expected nil bins on error, got %v", bins)
		}
	})
}
