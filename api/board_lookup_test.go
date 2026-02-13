package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestLookupBoardIDByName tests Story 5: Filter tickets by board name
//
// User Story:
// As a user, I want to provide a board name and get the corresponding board ID
// so that I can filter tickets by board using friendly names.
//
// Acceptance Criteria:
// - Given a board name, return the corresponding board ID
// - Matching is case-insensitive
// - Returns error if board name not found
// - Handles boards with special characters in names
// - Can use board ID to filter tickets via server-side filtering
func TestLookupBoardIDByName(t *testing.T) {
	t.Run("Given board name When looking up ID Then return matching board ID", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{"_id": "board123", "name": "Development Board", "bins": []},
				{"_id": "board456", "name": "Product Board", "bins": []},
				{"_id": "board789", "name": "Design Board", "bins": []}
			]`))
		}))
		defer server.Close()

		client := NewClient("test-key")
		client.baseURL = server.URL

		// Act
		boardID, err := client.LookupBoardIDByName("Development Board")

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if boardID != "board123" {
			t.Errorf("Expected board ID board123, got %s", boardID)
		}
	})

	t.Run("Given board name with different case When looking up ID Then return matching board ID", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{"_id": "board123", "name": "Product Board", "bins": []}
			]`))
		}))
		defer server.Close()

		client := NewClient("test-key")
		client.baseURL = server.URL

		// Act
		boardID, err := client.LookupBoardIDByName("product board")

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if boardID != "board123" {
			t.Errorf("Expected board ID board123, got %s", boardID)
		}
	})

	t.Run("Given non-existent board name When looking up ID Then return error", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{"_id": "board123", "name": "Product Board", "bins": []}
			]`))
		}))
		defer server.Close()

		client := NewClient("test-key")
		client.baseURL = server.URL

		// Act
		boardID, err := client.LookupBoardIDByName("Nonexistent")

		// Assert
		if err == nil {
			t.Fatal("Expected error for non-existent board, got nil")
		}
		if boardID != "" {
			t.Errorf("Expected empty board ID on error, got %s", boardID)
		}
	})

	t.Run("Given empty board list When looking up ID Then return error", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[]`))
		}))
		defer server.Close()

		client := NewClient("test-key")
		client.baseURL = server.URL

		// Act
		boardID, err := client.LookupBoardIDByName("Any Name")

		// Assert
		if err == nil {
			t.Fatal("Expected error for empty board list, got nil")
		}
		if boardID != "" {
			t.Errorf("Expected empty board ID on error, got %s", boardID)
		}
	})
}

// TestSearchTicketsWithBoardFilter tests filtering tickets by board ID
func TestSearchTicketsWithBoardFilter(t *testing.T) {
	t.Run("Given board ID When searching tickets Then include boards parameter in request", func(t *testing.T) {
		// Arrange
		requestReceived := false
		var requestURL string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestReceived = true
			requestURL = r.URL.String()
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{"_id": "ticket1", "name": "Ticket 1"}
			]`))
		}))
		defer server.Close()

		client := NewClient("test-key")
		client.baseURL = server.URL

		// Act
		_, err := client.SearchTicketsWithFilters([]string{"user123"}, "", "board456")

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !requestReceived {
			t.Fatal("Expected request to be made")
		}
		if !strings.Contains(requestURL, "boards=board456") {
			t.Errorf("Expected URL to contain boards=board456, got %s", requestURL)
		}
	})

	t.Run("Given board and bin IDs When searching Then combine both filters", func(t *testing.T) {
		// Arrange
		var requestURL string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestURL = r.URL.String()
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[]`))
		}))
		defer server.Close()

		client := NewClient("test-key")
		client.baseURL = server.URL

		// Act
		_, err := client.SearchTicketsWithFilters([]string{"user1"}, "bin123", "board456")

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !strings.Contains(requestURL, "bins=bin123") {
			t.Errorf("Expected URL to contain bins=bin123, got %s", requestURL)
		}
		if !strings.Contains(requestURL, "boards=board456") {
			t.Errorf("Expected URL to contain boards=board456, got %s", requestURL)
		}
	})
}
