package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestSearchTicketsWithBinFilter tests Story 3: Filter tickets by bin name using server-side filtering
//
// User Story:
// As a user, I want to filter tickets by bin name on the server side
// so that only relevant tickets are returned from the API.
//
// Acceptance Criteria:
// - SearchTickets accepts optional bin ID parameter
// - When bin ID provided, adds bins= query parameter to API request
// - Server filters tickets before returning them
// - Reduces data transfer and client-side processing
// - Works alongside user filtering
func TestSearchTicketsWithBinFilter(t *testing.T) {
	t.Run("Given bin ID When searching tickets Then include bins parameter in request", func(t *testing.T) {
		// Arrange
		requestReceived := false
		var requestURL string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestReceived = true
			requestURL = r.URL.String()
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{"_id": "ticket1", "name": "Ticket 1", "bin_id": "bin123"}
			]`))
		}))
		defer server.Close()

		client := NewClient("test-key")
		client.baseURL = server.URL

		// Act
		_, err := client.SearchTicketsWithFilters([]string{"user123"}, "bin123", "")

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !requestReceived {
			t.Fatal("Expected request to be made")
		}
		if !strings.Contains(requestURL, "bins=bin123") {
			t.Errorf("Expected URL to contain bins=bin123, got %s", requestURL)
		}
		if !strings.Contains(requestURL, "users=user123") {
			t.Errorf("Expected URL to contain users=user123, got %s", requestURL)
		}
	})

	t.Run("Given no bin filter When searching tickets Then omit bins parameter", func(t *testing.T) {
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
		_, err := client.SearchTicketsWithFilters([]string{"user123"}, "", "")

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if strings.Contains(requestURL, "bins=") {
			t.Errorf("Expected URL to not contain bins parameter, got %s", requestURL)
		}
	})

	t.Run("Given bin ID and user IDs When searching Then combine both filters", func(t *testing.T) {
		// Arrange
		var requestURL string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestURL = r.URL.String()
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{"_id": "ticket1", "name": "Ticket 1", "bin_id": "bin456", "assigned_ids": ["user1"]}
			]`))
		}))
		defer server.Close()

		client := NewClient("test-key")
		client.baseURL = server.URL

		// Act
		tickets, err := client.SearchTicketsWithFilters([]string{"user1", "user2"}, "bin456", "")

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !strings.Contains(requestURL, "bins=bin456") {
			t.Errorf("Expected URL to contain bins=bin456, got %s", requestURL)
		}
		if !strings.Contains(requestURL, "users=") {
			t.Errorf("Expected URL to contain users parameter, got %s", requestURL)
		}
		if len(tickets) != 1 {
			t.Errorf("Expected 1 ticket, got %d", len(tickets))
		}
	})

	t.Run("Given empty bin ID When searching Then treat as no filter", func(t *testing.T) {
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
		_, err := client.SearchTicketsWithFilters([]string{"user123"}, "", "")

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if strings.Contains(requestURL, "bins=") {
			t.Errorf("Expected URL to not contain bins parameter for empty bin ID, got %s", requestURL)
		}
	})
}
