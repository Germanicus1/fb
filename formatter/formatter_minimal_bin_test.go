package formatter

import (
	"strings"
	"testing"

	"github.com/Germanicus1/fb/models"
)

// Story 1.2: Minimal Output with Bin Filtering - Acceptance Tests

// TestStory1_2_MinimalFormatWithBinFilter verifies minimal format works with filtering
// Acceptance Criterion 1: When I run `fb --bin "In Progress"`, the output shows minimal format (ID + Name only)
func TestStory1_2_MinimalFormatWithBinFilter(t *testing.T) {
	// Given: I have tickets in different bins
	tickets := []models.Ticket{
		{ID: "TICKET-001", Name: "First Ticket", BinName: "In Progress"},
		{ID: "TICKET-002", Name: "Second Ticket", BinName: "In Progress"},
		{ID: "TICKET-003", Name: "Third Ticket", BinName: "Done"},
	}

	// When: I format filtered tickets (simulating --bin "In Progress") in minimal mode
	filteredTickets := []models.Ticket{tickets[0], tickets[1]}
	output := FormatTicketsMinimal(filteredTickets)

	// Then: Output shows minimal format (ID + Name only)
	if !strings.Contains(output, "[TICKET-001] First Ticket") {
		t.Error("Output should contain first ticket in minimal format")
	}
	if !strings.Contains(output, "[TICKET-002] Second Ticket") {
		t.Error("Output should contain second ticket in minimal format")
	}

	// And: No status, dates, or descriptions
	if strings.Contains(output, "In Progress") {
		t.Error("Minimal output should not contain bin name 'In Progress'")
	}
	if strings.Contains(output, "Status:") {
		t.Error("Minimal output should not contain 'Status:'")
	}
}

// TestStory1_2_FilteredResultsShowCorrectCount verifies header shows filtered count
// Acceptance Criterion 3: The header shows correct count for filtered results: "Found N ticket(s) assigned to you:"
func TestStory1_2_FilteredResultsShowCorrectCount(t *testing.T) {
	// Given: I have 5 total tickets but only 2 match the filter
	filteredTickets := []models.Ticket{
		{ID: "TICKET-001", Name: "First", BinName: "In Progress"},
		{ID: "TICKET-002", Name: "Second", BinName: "In Progress"},
	}

	// When: I format the filtered tickets in minimal mode
	output := FormatTicketsMinimal(filteredTickets)

	// Then: The header shows the count of filtered tickets (2), not total (5)
	if !strings.Contains(output, "Found 2 ticket(s) assigned to you:") {
		t.Errorf("Header should show filtered count of 2, got: %s", output)
	}
}

// TestStory1_2_AllFilteredTicketsDisplayed verifies all filtered tickets shown
// Acceptance Criterion 4: All filtered tickets are displayed in minimal format (one line per ticket)
func TestStory1_2_AllFilteredTicketsDisplayed(t *testing.T) {
	// Given: I have 3 tickets that match a filter
	filteredTickets := []models.Ticket{
		{ID: "TICKET-100", Name: "Task Alpha"},
		{ID: "TICKET-101", Name: "Task Beta"},
		{ID: "TICKET-102", Name: "Task Gamma"},
	}

	// When: I format the filtered tickets in minimal mode
	output := FormatTicketsMinimal(filteredTickets)

	// Then: All 3 filtered tickets are displayed
	if !strings.Contains(output, "TICKET-100") {
		t.Error("Should display TICKET-100")
	}
	if !strings.Contains(output, "TICKET-101") {
		t.Error("Should display TICKET-101")
	}
	if !strings.Contains(output, "TICKET-102") {
		t.Error("Should display TICKET-102")
	}

	// Each ticket on one line in minimal format
	lines := strings.Split(output, "\n")
	ticketLineCount := 0
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "[TICKET-") {
			ticketLineCount++
		}
	}

	if ticketLineCount != 3 {
		t.Errorf("Expected 3 ticket lines, got %d", ticketLineCount)
	}
}

// TestStory1_2_NoVerboseDetailsInFilteredOutput verifies no verbose details appear
// Acceptance Criterion 5: When combined with bin filtering, no verbose details appear in the output
func TestStory1_2_NoVerboseDetailsInFilteredOutput(t *testing.T) {
	// Given: I have filtered tickets with full details
	filteredTickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "Filtered Ticket",
			BinName:     "In Progress",
			Description: "This is a detailed description that should not appear",
		},
	}

	// When: I format the filtered tickets in minimal mode
	output := FormatTicketsMinimal(filteredTickets)

	// Then: No verbose details appear
	if strings.Contains(output, "Status:") {
		t.Error("Should not contain 'Status:'")
	}
	if strings.Contains(output, "Description:") {
		t.Error("Should not contain 'Description:'")
	}
	if strings.Contains(output, "Created:") {
		t.Error("Should not contain 'Created:'")
	}
	if strings.Contains(output, "Updated:") {
		t.Error("Should not contain 'Updated:'")
	}
	if strings.Contains(output, "detailed description") {
		t.Error("Should not contain description text")
	}
}

// TestStory1_2_SingleFilteredTicket verifies single filtered ticket works
func TestStory1_2_SingleFilteredTicket(t *testing.T) {
	// Given: Filter results in only one ticket
	filteredTickets := []models.Ticket{
		{ID: "TICKET-999", Name: "Only Match"},
	}

	// When: I format the filtered ticket in minimal mode
	output := FormatTicketsMinimal(filteredTickets)

	// Then: Shows correct count and format
	if !strings.Contains(output, "Found 1 ticket(s) assigned to you:") {
		t.Error("Should show count of 1")
	}
	if !strings.Contains(output, "[TICKET-999] Only Match") {
		t.Error("Should display the single ticket in minimal format")
	}
}

// TestStory1_2_EmptyFilterResults verifies empty filter results handled gracefully
func TestStory1_2_EmptyFilterResults(t *testing.T) {
	// Given: Filter results in zero tickets
	filteredTickets := []models.Ticket{}

	// When: I format the empty filtered results in minimal mode
	output := FormatTicketsMinimal(filteredTickets)

	// Then: Shows appropriate "no tickets" message
	if !strings.Contains(output, "No tickets assigned to you.") {
		t.Errorf("Should show 'No tickets assigned to you.' for empty filter results, got: %s", output)
	}
}

// TestStory1_2_FilterPreservesMinimalFormat verifies filtering doesn't change format
func TestStory1_2_FilterPreservesMinimalFormat(t *testing.T) {
	// Given: I have the same tickets, some filtered and some not
	allTickets := []models.Ticket{
		{ID: "TICKET-001", Name: "First Ticket"},
		{ID: "TICKET-002", Name: "Second Ticket"},
		{ID: "TICKET-003", Name: "Third Ticket"},
	}
	filteredTickets := []models.Ticket{
		{ID: "TICKET-001", Name: "First Ticket"},
		{ID: "TICKET-002", Name: "Second Ticket"},
	}

	// When: I format both unfiltered and filtered in minimal mode
	unfilteredOutput := FormatTicketsMinimal(allTickets)
	filteredOutput := FormatTicketsMinimal(filteredTickets)

	// Then: The format structure is identical (just different counts and tickets)
	// Both should have header + blank line + ticket lines
	unfilteredLines := strings.Split(strings.TrimSpace(unfilteredOutput), "\n")
	filteredLines := strings.Split(strings.TrimSpace(filteredOutput), "\n")

	// Header line format should be same
	if !strings.HasPrefix(unfilteredLines[0], "Found") {
		t.Error("Unfiltered should have 'Found' header")
	}
	if !strings.HasPrefix(filteredLines[0], "Found") {
		t.Error("Filtered should have 'Found' header")
	}

	// Second line should be blank in both
	if len(unfilteredLines) > 1 && unfilteredLines[1] != "" {
		t.Error("Unfiltered second line should be blank")
	}
	if len(filteredLines) > 1 && filteredLines[1] != "" {
		t.Error("Filtered second line should be blank")
	}

	// Ticket lines should follow [ID] Name format
	if !strings.HasPrefix(unfilteredLines[2], "[TICKET-") {
		t.Error("Unfiltered tickets should use [ID] Name format")
	}
	if !strings.HasPrefix(filteredLines[2], "[TICKET-") {
		t.Error("Filtered tickets should use [ID] Name format")
	}
}

// TestStory1_2_ManyFilteredTickets verifies minimal format works with many filtered results
func TestStory1_2_ManyFilteredTickets(t *testing.T) {
	// Given: Filter results in 20 tickets
	filteredTickets := make([]models.Ticket, 20)
	for i := 0; i < 20; i++ {
		filteredTickets[i] = models.Ticket{
			ID:      "TICKET-" + string(rune('A'+i)),
			Name:    "Filtered Task " + string(rune('A'+i)),
			BinName: "In Progress",
		}
	}

	// When: I format the filtered tickets in minimal mode
	output := FormatTicketsMinimal(filteredTickets)

	// Then: All 20 tickets are shown in minimal format
	if !strings.Contains(output, "Found 20 ticket(s)") {
		t.Error("Should show count of 20")
	}

	// Verify output is compact (approximately 22 lines: header + blank + 20 tickets)
	lines := strings.Split(output, "\n")
	lineCount := len(lines) - 1 // Remove trailing empty from split

	if lineCount < 21 || lineCount > 23 {
		t.Errorf("Expected approximately 22 lines, got %d", lineCount)
	}
}
