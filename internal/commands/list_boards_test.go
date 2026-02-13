package commands

import (
	"strings"
	"testing"

	"github.com/Germanicus1/fb/models"
)

// TestListBoardsCommand tests Story 7: Add --list-boards command
//
// User Story:
// As a user, I want to see all available boards with their names
// so that I know which board names I can use for filtering.
//
// Acceptance Criteria:
// - --list-boards flag shows all boards
// - Each board displays ID and name
// - Output is formatted for easy reading
// - Empty board list shows appropriate message
func TestListBoardsCommand(t *testing.T) {
	t.Run("Given boards When formatting list Then display all boards with IDs and names", func(t *testing.T) {
		boards := []models.Board{
			{ID: "board3", Name: "Design Board", Bins: []string{}},
			{ID: "board1", Name: "Development Board", Bins: []string{"bin1"}},
			{ID: "board2", Name: "Product Board", Bins: []string{"bin2", "bin3"}},
		}

		output := formatBoardList(boards)

		// Assert
		if output == "" {
			t.Fatal("Expected non-empty output")
		}

		// Verify all boards are in the output
		if !strings.Contains(output, "Development Board") {
			t.Error("Expected output to contain 'Development Board'")
		}
		if !strings.Contains(output, "Product Board") {
			t.Error("Expected output to contain 'Product Board'")
		}
		if !strings.Contains(output, "Design Board") {
			t.Error("Expected output to contain 'Design Board'")
		}

		// Verify IDs are in the output
		if !strings.Contains(output, "board1") {
			t.Error("Expected output to contain board1 ID")
		}
		if !strings.Contains(output, "board2") {
			t.Error("Expected output to contain board2 ID")
		}
		if !strings.Contains(output, "board3") {
			t.Error("Expected output to contain board3 ID")
		}
	})

	t.Run("Given empty board list When formatting list Then show appropriate message", func(t *testing.T) {
		boards := []models.Board{}

		output := formatBoardList(boards)

		// Assert
		if output == "" {
			t.Fatal("Expected non-empty output for empty list")
		}
		if !strings.Contains(strings.ToLower(output), "no boards") {
			t.Error("Expected message about no boards found")
		}
	})

	t.Run("Given boards When formatting list Then include header", func(t *testing.T) {
		boards := []models.Board{
			{ID: "board1", Name: "Development Board", Bins: []string{}},
		}

		output := formatBoardList(boards)

		// Assert - should have some kind of header or title
		if !strings.Contains(strings.ToLower(output), "boards") && !strings.Contains(strings.ToLower(output), "available") {
			t.Error("Expected output to contain header with 'boards' or 'available'")
		}
	})
}
