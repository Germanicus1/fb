package filter

import (
	"testing"

	"github.com/Germanicus1/fb/models"
)

// TestFilterTicketsByBinName tests Story 7: Filter Tickets by Bin Name
//
// Acceptance Criteria:
// - User can filter tickets by bin name
// - Filtering matches bin name exactly (case-insensitive comparison)
// - Only tickets with BinName matching the filter value are displayed
// - Tickets from all boards are included if they match the bin filter
// - Empty result set is handled gracefully
// - Filter works with bin names containing spaces and special characters
func TestFilterTicketsByBinName(t *testing.T) {
	t.Run("Given tickets When filtering by exact bin name Then return matching tickets", func(t *testing.T) {
		// Arrange
		tickets := []models.Ticket{
			{ID: "1", Name: "Ticket 1", BinName: "In Progress", BinID: "bin1"},
			{ID: "2", Name: "Ticket 2", BinName: "Done", BinID: "bin2"},
			{ID: "3", Name: "Ticket 3", BinName: "In Progress", BinID: "bin1"},
		}

		// Act
		filtered := FilterByBinName(tickets, "In Progress")

		// Assert
		if len(filtered) != 2 {
			t.Fatalf("Expected 2 tickets, got %d", len(filtered))
		}
		if filtered[0].ID != "1" || filtered[1].ID != "3" {
			t.Errorf("Expected tickets 1 and 3, got %v", filtered)
		}
	})

	t.Run("Given tickets When filtering with case-insensitive match Then return matching tickets", func(t *testing.T) {
		// Arrange
		tickets := []models.Ticket{
			{ID: "1", Name: "Ticket 1", BinName: "In Progress", BinID: "bin1"},
			{ID: "2", Name: "Ticket 2", BinName: "Done", BinID: "bin2"},
		}

		// Act
		filtered := FilterByBinName(tickets, "in progress")

		// Assert
		if len(filtered) != 1 {
			t.Fatalf("Expected 1 ticket, got %d", len(filtered))
		}
		if filtered[0].ID != "1" {
			t.Errorf("Expected ticket 1, got %s", filtered[0].ID)
		}
	})

	t.Run("Given tickets When filtering by non-existent bin Then return empty list", func(t *testing.T) {
		// Arrange
		tickets := []models.Ticket{
			{ID: "1", Name: "Ticket 1", BinName: "In Progress", BinID: "bin1"},
			{ID: "2", Name: "Ticket 2", BinName: "Done", BinID: "bin2"},
		}

		// Act
		filtered := FilterByBinName(tickets, "Blocked")

		// Assert
		if len(filtered) != 0 {
			t.Errorf("Expected 0 tickets, got %d", len(filtered))
		}
	})

	t.Run("Given tickets When filtering with bin name containing spaces Then match correctly", func(t *testing.T) {
		// Arrange
		tickets := []models.Ticket{
			{ID: "1", Name: "Ticket 1", BinName: "Ready for Review", BinID: "bin1"},
			{ID: "2", Name: "Ticket 2", BinName: "Done", BinID: "bin2"},
		}

		// Act
		filtered := FilterByBinName(tickets, "Ready for Review")

		// Assert
		if len(filtered) != 1 {
			t.Fatalf("Expected 1 ticket, got %d", len(filtered))
		}
	})

	t.Run("Given empty ticket list When filtering Then return empty list", func(t *testing.T) {
		// Arrange
		tickets := []models.Ticket{}

		// Act
		filtered := FilterByBinName(tickets, "In Progress")

		// Assert
		if len(filtered) != 0 {
			t.Errorf("Expected 0 tickets, got %d", len(filtered))
		}
	})
}

// TestFilterTicketsByBinID tests Story 8: Filter Tickets by Bin ID
//
// Acceptance Criteria:
// - --bin flag accepts both bin names and bin IDs
// - When a value matches a BinID exactly, filtering uses BinID
// - When a value doesn't match any BinID, filtering falls back to BinName matching
// - User can filter using bin ID for precision
// - ID-based filtering is case-sensitive and exact match
// - Name-based filtering remains case-insensitive
func TestFilterTicketsByBinID(t *testing.T) {
	t.Run("Given tickets When filtering by bin ID Then return matching tickets", func(t *testing.T) {
		// Arrange
		tickets := []models.Ticket{
			{ID: "1", Name: "Ticket 1", BinName: "In Progress", BinID: "bin-12345"},
			{ID: "2", Name: "Ticket 2", BinName: "Done", BinID: "bin-67890"},
			{ID: "3", Name: "Ticket 3", BinName: "In Progress", BinID: "bin-12345"},
		}

		// Act
		filtered := FilterByBinName(tickets, "bin-12345")

		// Assert
		if len(filtered) != 2 {
			t.Fatalf("Expected 2 tickets, got %d", len(filtered))
		}
		if filtered[0].ID != "1" || filtered[1].ID != "3" {
			t.Errorf("Expected tickets 1 and 3")
		}
	})

	t.Run("Given tickets When BinID matches exactly Then use BinID match", func(t *testing.T) {
		// Arrange
		tickets := []models.Ticket{
			{ID: "1", Name: "Ticket 1", BinName: "test", BinID: "test"},
			{ID: "2", Name: "Ticket 2", BinName: "other", BinID: "bin-123"},
		}

		// Act - Filter by value that matches both BinID of ticket 1 and BinName of ticket 1
		filtered := FilterByBinName(tickets, "test")

		// Assert - Should match ticket 1 by both ID and name
		if len(filtered) != 1 {
			t.Fatalf("Expected 1 ticket, got %d", len(filtered))
		}
		if filtered[0].ID != "1" {
			t.Errorf("Expected ticket 1")
		}
	})

	t.Run("Given tickets When BinID doesn't match Then fall back to BinName", func(t *testing.T) {
		// Arrange
		tickets := []models.Ticket{
			{ID: "1", Name: "Ticket 1", BinName: "In Progress", BinID: "bin-12345"},
			{ID: "2", Name: "Ticket 2", BinName: "Done", BinID: "bin-67890"},
		}

		// Act - Filter by name (no ID match)
		filtered := FilterByBinName(tickets, "in progress")

		// Assert - Should match by name (case-insensitive)
		if len(filtered) != 1 {
			t.Fatalf("Expected 1 ticket, got %d", len(filtered))
		}
		if filtered[0].ID != "1" {
			t.Errorf("Expected ticket 1")
		}
	})

	t.Run("Given tickets When filtering by BinID Then match is case-sensitive", func(t *testing.T) {
		// Arrange
		tickets := []models.Ticket{
			{ID: "1", Name: "Ticket 1", BinName: "Test", BinID: "BIN-123"},
			{ID: "2", Name: "Ticket 2", BinName: "Done", BinID: "bin-456"},
		}

		// Act - Try with different case
		filtered := FilterByBinName(tickets, "bin-123")

		// Assert - Should not match (ID is case-sensitive)
		// But might match by name if "bin-123" case-insensitively matches any BinName
		// In this case, no match expected
		if len(filtered) != 0 {
			t.Errorf("Expected 0 tickets (case-sensitive ID match), got %d", len(filtered))
		}
	})

	t.Run("Given tickets When value matches neither ID nor Name Then return empty", func(t *testing.T) {
		// Arrange
		tickets := []models.Ticket{
			{ID: "1", Name: "Ticket 1", BinName: "In Progress", BinID: "bin-12345"},
			{ID: "2", Name: "Ticket 2", BinName: "Done", BinID: "bin-67890"},
		}

		// Act
		filtered := FilterByBinName(tickets, "nonexistent")

		// Assert
		if len(filtered) != 0 {
			t.Errorf("Expected 0 tickets, got %d", len(filtered))
		}
	})
}
