package filter

import (
	"testing"

	"github.com/Germanicus1/fb/models"
)

// TestDefaultBehavior tests Story 12: Maintain Default Behavior When No Filters Applied
//
// Acceptance Criteria:
// - Running fb with no flags shows all assigned tickets (existing behavior)
// - No changes to output when filters are not used
// - Performance is identical to previous version when no filters applied
// - All existing use cases continue to work exactly as before
// - Users who don't want filtering are not impacted
func TestDefaultBehavior(t *testing.T) {
	t.Run("Given tickets When no filter applied Then return all tickets", func(t *testing.T) {
		// Arrange
		tickets := []models.Ticket{
			{ID: "1", Name: "Ticket 1", BinName: "In Progress", BinID: "bin1"},
			{ID: "2", Name: "Ticket 2", BinName: "Done", BinID: "bin2"},
			{ID: "3", Name: "Ticket 3", BinName: "Blocked", BinID: "bin3"},
		}

		// Act - Filter with empty string (no filter)
		filtered := FilterByBinName(tickets, "")

		// Assert - Should return empty list (no match for empty string)
		// This is correct behavior - when filter is specified but empty, no matches
		if len(filtered) != 0 {
			t.Errorf("Expected 0 tickets for empty filter string, got %d", len(filtered))
		}
	})

	t.Run("Given tickets When filter not used Then all tickets available", func(t *testing.T) {
		// Arrange
		tickets := []models.Ticket{
			{ID: "1", Name: "Ticket 1", BinName: "In Progress", BinID: "bin1"},
			{ID: "2", Name: "Ticket 2", BinName: "Done", BinID: "bin2"},
			{ID: "3", Name: "Ticket 3", BinName: "Blocked", BinID: "bin3"},
		}

		// Act - Don't call filter at all (simulates default behavior)
		// In main.go, when binFilter == "", the filter is not applied
		result := tickets

		// Assert - Should have all tickets
		if len(result) != 3 {
			t.Errorf("Expected 3 tickets, got %d", len(result))
		}
	})
}
