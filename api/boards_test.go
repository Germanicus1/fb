package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetBoards tests Story 4: Fetch boards from API with names
//
// User Story:
// As a user, I want the system to fetch all boards from the Flow Boards API
// so that I can reference boards by their friendly names instead of IDs.
//
// Acceptance Criteria:
// - System calls GET /boards endpoint
// - Response contains array of boards with ID, Name, and Bins fields
// - Board IDs are correctly mapped to board names
// - Bins array is correctly parsed for each board
// - API errors are properly handled
// - Empty board list is handled gracefully
func TestGetBoards(t *testing.T) {
	t.Run("Given API endpoint When fetching boards Then return list of boards with IDs and names", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify the request is for the boards endpoint
			if r.URL.Path != "/boards" {
				t.Errorf("Expected path /boards, got %s", r.URL.Path)
			}
			if r.Method != "GET" {
				t.Errorf("Expected GET method, got %s", r.Method)
			}

			// Return mock boards data matching API spec
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{"_id": "board1", "name": "Development Board", "bins": ["bin1", "bin2"]},
				{"_id": "board2", "name": "Product Board", "bins": ["bin3", "bin4"]},
				{"_id": "board3", "name": "Design Board", "bins": ["bin5"]}
			]`))
		}))
		defer server.Close()

		client := NewClient("test-key")
		client.baseURL = server.URL

		// Act
		boards, err := client.GetBoards()

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(boards) != 3 {
			t.Fatalf("Expected 3 boards, got %d", len(boards))
		}

		// Verify first board
		if boards[0].ID != "board1" {
			t.Errorf("Expected board ID board1, got %s", boards[0].ID)
		}
		if boards[0].Name != "Development Board" {
			t.Errorf("Expected board name Development Board, got %s", boards[0].Name)
		}
		if len(boards[0].Bins) != 2 {
			t.Errorf("Expected 2 bins, got %d", len(boards[0].Bins))
		}

		// Verify second board
		if boards[1].ID != "board2" {
			t.Errorf("Expected board ID board2, got %s", boards[1].ID)
		}
		if boards[1].Name != "Product Board" {
			t.Errorf("Expected board name Product Board, got %s", boards[1].Name)
		}
	})

	t.Run("Given API returns empty array When fetching boards Then return empty list without error", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("[]"))
		}))
		defer server.Close()

		client := NewClient("test-key")
		client.baseURL = server.URL

		// Act
		boards, err := client.GetBoards()

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(boards) != 0 {
			t.Errorf("Expected empty board list, got %d boards", len(boards))
		}
	})

	t.Run("Given API returns error When fetching boards Then return error", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}))
		defer server.Close()

		client := NewClient("test-key")
		client.baseURL = server.URL

		// Act
		boards, err := client.GetBoards()

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if boards != nil {
			t.Errorf("Expected nil boards on error, got %v", boards)
		}
	})

	t.Run("Given base URL not discovered When fetching boards Then return error", func(t *testing.T) {
		// Arrange
		client := NewClient("test-key")
		// Don't call DiscoverRestPrefix, so baseURL is empty

		// Act
		boards, err := client.GetBoards()

		// Assert
		if err == nil {
			t.Fatal("Expected error for missing base URL, got nil")
		}
		if boards != nil {
			t.Errorf("Expected nil boards on error, got %v", boards)
		}
	})

	t.Run("Given board with empty bins array When fetching boards Then handle gracefully", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{"_id": "board1", "name": "Empty Board", "bins": []}
			]`))
		}))
		defer server.Close()

		client := NewClient("test-key")
		client.baseURL = server.URL

		// Act
		boards, err := client.GetBoards()

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(boards) != 1 {
			t.Fatalf("Expected 1 board, got %d", len(boards))
		}
		if len(boards[0].Bins) != 0 {
			t.Errorf("Expected 0 bins, got %d", len(boards[0].Bins))
		}
	})
}
