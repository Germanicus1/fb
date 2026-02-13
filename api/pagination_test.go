package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetBinsPagination tests pagination support for bins endpoint
//
// User Story:
// As a user, I want GetBins() to return ALL bins (not just the first page),
// so that I can filter by any bin name regardless of pagination.
//
// Acceptance Criteria:
// 1. GetBins() fetches all pages of bins until no more data
// 2. Uses query parameters: ?max-results=1000 to minimize API calls
// 3. Follows page-token in responses to get next page
// 4. Existing tests still pass
// 5. New tests verify pagination works correctly
func TestGetBinsPagination(t *testing.T) {
	t.Run("Given bins across multiple pages When fetching bins Then return all bins from all pages", func(t *testing.T) {
		// Arrange
		requestCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify the request is for the bins endpoint
			if r.URL.Path != "/bins" {
				t.Errorf("Expected path /bins, got %s", r.URL.Path)
			}
			if r.Method != "GET" {
				t.Errorf("Expected GET method, got %s", r.Method)
			}

			// Verify max-results parameter is set
			maxResults := r.URL.Query().Get("max-results")
			if maxResults != "1000" {
				t.Errorf("Expected max-results=1000, got %s", maxResults)
			}

			requestCount++
			pageToken := r.URL.Query().Get("page-token")

			w.WriteHeader(http.StatusOK)

			// First page - return 2 bins with page-token
			if pageToken == "" {
				w.Write([]byte(`{
					"results": [
						{"_id": "bin1", "name": "Bin One"},
						{"_id": "bin2", "name": "Bin Two"}
					],
					"page-token": "token123"
				}`))
				return
			}

			// Second page - return 2 more bins with page-token
			if pageToken == "token123" {
				w.Write([]byte(`{
					"results": [
						{"_id": "cx7oRn0CK1SoAMn0x", "name": "K+Dev.Doing"},
						{"_id": "bin4", "name": "Bin Four"}
					],
					"page-token": "token456"
				}`))
				return
			}

			// Third page - return 1 bin with no page-token (last page)
			if pageToken == "token456" {
				w.Write([]byte(`{
					"results": [
						{"_id": "bin5", "name": "Bin Five"}
					]
				}`))
				return
			}

			t.Errorf("Unexpected page-token: %s", pageToken)
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

		// Should have made 3 requests (3 pages)
		if requestCount != 3 {
			t.Errorf("Expected 3 API requests, got %d", requestCount)
		}

		// Should have all 5 bins from all pages
		if len(bins) != 5 {
			t.Fatalf("Expected 5 bins total, got %d", len(bins))
		}

		// Verify bins from all pages are present
		expectedBins := map[string]string{
			"bin1":           "Bin One",
			"bin2":           "Bin Two",
			"cx7oRn0CK1SoAMn0x": "K+Dev.Doing",
			"bin4":           "Bin Four",
			"bin5":           "Bin Five",
		}

		for i, bin := range bins {
			expectedName, ok := expectedBins[bin.ID]
			if !ok {
				t.Errorf("Unexpected bin ID at index %d: %s", i, bin.ID)
			}
			if bin.Name != expectedName {
				t.Errorf("Expected bin %s to have name %s, got %s", bin.ID, expectedName, bin.Name)
			}
		}
	})

	t.Run("Given single page of bins When fetching bins Then return all bins without extra requests", func(t *testing.T) {
		// Arrange
		requestCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCount++
			w.WriteHeader(http.StatusOK)
			// Return bins without page-token (single page)
			w.Write([]byte(`{
				"results": [
					{"_id": "bin1", "name": "Only Bin"}
				]
			}`))
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

		// Should have made only 1 request
		if requestCount != 1 {
			t.Errorf("Expected 1 API request, got %d", requestCount)
		}

		if len(bins) != 1 {
			t.Fatalf("Expected 1 bin, got %d", len(bins))
		}
	})

	t.Run("Given empty results on first page When fetching bins Then return empty list", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"results": []}`))
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

	t.Run("Given API error on second page When fetching bins Then return error", func(t *testing.T) {
		// Arrange
		requestCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCount++
			pageToken := r.URL.Query().Get("page-token")

			if pageToken == "" {
				// First page succeeds
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"results": [{"_id": "bin1", "name": "Bin One"}],
					"page-token": "token123"
				}`))
				return
			}

			// Second page fails
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
			t.Fatal("Expected error on second page, got nil")
		}
		if bins != nil {
			t.Errorf("Expected nil bins on error, got %v", bins)
		}
	})
}

// TestGetBoardsPagination tests pagination support for boards endpoint
//
// User Story:
// As a user, I want GetBoards() to return ALL boards (not just the first page),
// so that I can filter by any board name regardless of pagination.
//
// Acceptance Criteria:
// 1. GetBoards() fetches all pages of boards until no more data
// 2. Uses query parameters: ?max-results=1000 to minimize API calls
// 3. Follows page-token in responses to get next page
// 4. Existing tests still pass
// 5. New tests verify pagination works correctly
func TestGetBoardsPagination(t *testing.T) {
	t.Run("Given boards across multiple pages When fetching boards Then return all boards from all pages", func(t *testing.T) {
		// Arrange
		requestCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify the request is for the boards endpoint
			if r.URL.Path != "/boards" {
				t.Errorf("Expected path /boards, got %s", r.URL.Path)
			}
			if r.Method != "GET" {
				t.Errorf("Expected GET method, got %s", r.Method)
			}

			// Verify max-results parameter is set
			maxResults := r.URL.Query().Get("max-results")
			if maxResults != "1000" {
				t.Errorf("Expected max-results=1000, got %s", maxResults)
			}

			requestCount++
			pageToken := r.URL.Query().Get("page-token")

			w.WriteHeader(http.StatusOK)

			// First page - return 2 boards with page-token
			if pageToken == "" {
				w.Write([]byte(`{
					"results": [
						{"_id": "board1", "name": "Board One", "bins": ["bin1"]},
						{"_id": "board2", "name": "Board Two", "bins": ["bin2", "bin3"]}
					],
					"page-token": "token123"
				}`))
				return
			}

			// Second page - return 2 more boards with no page-token (last page)
			if pageToken == "token123" {
				w.Write([]byte(`{
					"results": [
						{"_id": "board3", "name": "Board Three", "bins": []},
						{"_id": "board4", "name": "Board Four", "bins": ["bin4"]}
					]
				}`))
				return
			}

			t.Errorf("Unexpected page-token: %s", pageToken)
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

		// Should have made 2 requests (2 pages)
		if requestCount != 2 {
			t.Errorf("Expected 2 API requests, got %d", requestCount)
		}

		// Should have all 4 boards from all pages
		if len(boards) != 4 {
			t.Fatalf("Expected 4 boards total, got %d", len(boards))
		}

		// Verify boards from all pages are present
		expectedBoards := map[string]string{
			"board1": "Board One",
			"board2": "Board Two",
			"board3": "Board Three",
			"board4": "Board Four",
		}

		for i, board := range boards {
			expectedName, ok := expectedBoards[board.ID]
			if !ok {
				t.Errorf("Unexpected board ID at index %d: %s", i, board.ID)
			}
			if board.Name != expectedName {
				t.Errorf("Expected board %s to have name %s, got %s", board.ID, expectedName, board.Name)
			}
		}
	})

	t.Run("Given single page of boards When fetching boards Then return all boards without extra requests", func(t *testing.T) {
		// Arrange
		requestCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCount++
			w.WriteHeader(http.StatusOK)
			// Return boards without page-token (single page)
			w.Write([]byte(`{
				"results": [
					{"_id": "board1", "name": "Only Board", "bins": []}
				]
			}`))
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

		// Should have made only 1 request
		if requestCount != 1 {
			t.Errorf("Expected 1 API request, got %d", requestCount)
		}

		if len(boards) != 1 {
			t.Fatalf("Expected 1 board, got %d", len(boards))
		}
	})

	t.Run("Given empty results on first page When fetching boards Then return empty list", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"results": []}`))
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

	t.Run("Given API error on second page When fetching boards Then return error", func(t *testing.T) {
		// Arrange
		requestCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCount++
			pageToken := r.URL.Query().Get("page-token")

			if pageToken == "" {
				// First page succeeds
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"results": [{"_id": "board1", "name": "Board One", "bins": []}],
					"page-token": "token123"
				}`))
				return
			}

			// Second page fails
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
			t.Fatal("Expected error on second page, got nil")
		}
		if boards != nil {
			t.Errorf("Expected nil boards on error, got %v", boards)
		}
	})
}

// TestBackwardsCompatibilityWithOldAPIResponse tests that the new pagination code
// still works with the old API response format (array instead of object with results)
func TestBackwardsCompatibilityBins(t *testing.T) {
	t.Run("Given old API response format When fetching bins Then parse successfully", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			// Old format: direct array
			w.Write([]byte(`[
				{"_id": "bin1", "name": "Bin One"},
				{"_id": "bin2", "name": "Bin Two"}
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
		if len(bins) != 2 {
			t.Fatalf("Expected 2 bins, got %d", len(bins))
		}
		if bins[0].ID != "bin1" {
			t.Errorf("Expected bin1, got %s", bins[0].ID)
		}
	})
}

func TestBackwardsCompatibilityBoards(t *testing.T) {
	t.Run("Given old API response format When fetching boards Then parse successfully", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			// Old format: direct array
			w.Write([]byte(`[
				{"_id": "board1", "name": "Board One", "bins": []},
				{"_id": "board2", "name": "Board Two", "bins": ["bin1"]}
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
		if len(boards) != 2 {
			t.Fatalf("Expected 2 boards, got %d", len(boards))
		}
		if boards[0].ID != "board1" {
			t.Errorf("Expected board1, got %s", boards[0].ID)
		}
	})
}
