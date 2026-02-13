package commands

import (
	"strings"
	"testing"

	"github.com/Germanicus1/fb/models"
)

// TestListBinsCommand tests Story 6: Add --list-bins command
//
// User Story:
// As a user, I want to see all available bins with their names
// so that I know which bin names I can use for filtering.
//
// Acceptance Criteria:
// - --list-bins flag shows all bins
// - Each bin displays ID and name
// - Output is formatted for easy reading
// - Empty bin list shows appropriate message
func TestListBinsCommand(t *testing.T) {
	t.Run("Given bins When formatting list Then display all bins with IDs and names", func(t *testing.T) {
		bins := []models.Bin{
			{ID: "bin3", Name: "Done"},
			{ID: "bin1", Name: "In Progress"},
			{ID: "bin2", Name: "K+Dev.Doing"},
		}

		output := formatBinList(bins)

		// Assert
		if output == "" {
			t.Fatal("Expected non-empty output")
		}

		// Verify all bins are in the output
		if !strings.Contains(output, "In Progress") {
			t.Error("Expected output to contain 'In Progress'")
		}
		if !strings.Contains(output, "K+Dev.Doing") {
			t.Error("Expected output to contain 'K+Dev.Doing'")
		}
		if !strings.Contains(output, "Done") {
			t.Error("Expected output to contain 'Done'")
		}

		// Verify IDs are in the output
		if !strings.Contains(output, "bin1") {
			t.Error("Expected output to contain bin1 ID")
		}
		if !strings.Contains(output, "bin2") {
			t.Error("Expected output to contain bin2 ID")
		}
		if !strings.Contains(output, "bin3") {
			t.Error("Expected output to contain bin3 ID")
		}
	})

	t.Run("Given empty bin list When formatting list Then show appropriate message", func(t *testing.T) {
		bins := []models.Bin{}

		output := formatBinList(bins)

		// Assert
		if output == "" {
			t.Fatal("Expected non-empty output for empty list")
		}
		if !strings.Contains(strings.ToLower(output), "no bins") {
			t.Error("Expected message about no bins found")
		}
	})

	t.Run("Given bins When formatting list Then include header", func(t *testing.T) {
		bins := []models.Bin{
			{ID: "bin1", Name: "In Progress"},
		}

		output := formatBinList(bins)

		// Assert - should have some kind of header or title
		if !strings.Contains(strings.ToLower(output), "bins") && !strings.Contains(strings.ToLower(output), "available") {
			t.Error("Expected output to contain header with 'bins' or 'available'")
		}
	})
}
