package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestLookupBinIDByName tests Story 2: Look up bin ID from bin name
//
// User Story:
// As a user, I want to provide a bin name and get the corresponding bin ID
// so that I can use friendly names instead of remembering IDs.
//
// Acceptance Criteria:
// - Given a bin name, return the corresponding bin ID
// - Matching is case-insensitive
// - Returns error if bin name not found
// - Handles bins with special characters in names
// - Returns first match if multiple bins have same name
func TestLookupBinIDByName(t *testing.T) {
	t.Run("Given bin name When looking up ID Then return matching bin ID", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{"_id": "cx7oRn0CK1SoAMn0x", "name": "K+Dev.Doing"},
				{"_id": "bin123", "name": "In Progress"},
				{"_id": "bin456", "name": "Done"}
			]`))
		}))
		defer server.Close()

		client := NewClient("test-key")
		client.baseURL = server.URL

		// Act
		binID, err := client.LookupBinIDByName("K+Dev.Doing")

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if binID != "cx7oRn0CK1SoAMn0x" {
			t.Errorf("Expected bin ID cx7oRn0CK1SoAMn0x, got %s", binID)
		}
	})

	t.Run("Given bin name with different case When looking up ID Then return matching bin ID", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{"_id": "bin123", "name": "In Progress"}
			]`))
		}))
		defer server.Close()

		client := NewClient("test-key")
		client.baseURL = server.URL

		// Act
		binID, err := client.LookupBinIDByName("in progress")

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if binID != "bin123" {
			t.Errorf("Expected bin ID bin123, got %s", binID)
		}
	})

	t.Run("Given non-existent bin name When looking up ID Then return error", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{"_id": "bin123", "name": "In Progress"}
			]`))
		}))
		defer server.Close()

		client := NewClient("test-key")
		client.baseURL = server.URL

		// Act
		binID, err := client.LookupBinIDByName("Nonexistent")

		// Assert
		if err == nil {
			t.Fatal("Expected error for non-existent bin, got nil")
		}
		if binID != "" {
			t.Errorf("Expected empty bin ID on error, got %s", binID)
		}
	})

	t.Run("Given bin name with special characters When looking up ID Then return matching bin ID", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{"_id": "special123", "name": "K+Dev.Doing"},
				{"_id": "special456", "name": "Review & Merge"}
			]`))
		}))
		defer server.Close()

		client := NewClient("test-key")
		client.baseURL = server.URL

		// Act
		binID, err := client.LookupBinIDByName("Review & Merge")

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if binID != "special456" {
			t.Errorf("Expected bin ID special456, got %s", binID)
		}
	})

	t.Run("Given empty bin list When looking up ID Then return error", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[]`))
		}))
		defer server.Close()

		client := NewClient("test-key")
		client.baseURL = server.URL

		// Act
		binID, err := client.LookupBinIDByName("Any Name")

		// Assert
		if err == nil {
			t.Fatal("Expected error for empty bin list, got nil")
		}
		if binID != "" {
			t.Errorf("Expected empty bin ID on error, got %s", binID)
		}
	})
}
