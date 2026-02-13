package formatter

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Germanicus1/fb/models"
)

// Story 1.1: Basic Minimal List Output - Acceptance Tests

// TestStory1_1_MinimalOutputStructure verifies the basic structure of minimal output
// Acceptance Criterion 1: When I run `fb` command with no arguments, I see output with this structure:
//   - First line: "Found N ticket(s) assigned to you:"
//   - Second line: Blank
//   - Remaining lines: One line per ticket in format `[ID] Name`
func TestStory1_1_MinimalOutputStructure(t *testing.T) {
	// Given: I have multiple tickets assigned to me
	tickets := []models.Ticket{
		{ID: "TICKET-001", Name: "First Ticket"},
		{ID: "TICKET-002", Name: "Second Ticket"},
		{ID: "TICKET-003", Name: "Third Ticket"},
	}

	// When: I format tickets in minimal mode
	output := FormatTicketsMinimal(tickets)

	// Then: I see output with the expected structure
	lines := strings.Split(output, "\n")

	// First line should be the header
	if !strings.Contains(lines[0], "Found 3 ticket(s) assigned to you:") {
		t.Errorf("First line should be 'Found 3 ticket(s) assigned to you:', got: %s", lines[0])
	}

	// Second line should be blank
	if strings.TrimSpace(lines[1]) != "" {
		t.Errorf("Second line should be blank, got: %s", lines[1])
	}

	// Remaining lines should be tickets
	if !strings.Contains(lines[2], "[TICKET-001] First Ticket") {
		t.Errorf("Third line should be '[TICKET-001] First Ticket', got: %s", lines[2])
	}
	if !strings.Contains(lines[3], "[TICKET-002] Second Ticket") {
		t.Errorf("Fourth line should be '[TICKET-002] Second Ticket', got: %s", lines[3])
	}
	if !strings.Contains(lines[4], "[TICKET-003] Third Ticket") {
		t.Errorf("Fifth line should be '[TICKET-003] Third Ticket', got: %s", lines[4])
	}
}

// TestStory1_1_OneLinePerTicket verifies exactly one line per ticket with no additional details
// Acceptance Criterion 2: The output shows exactly one line per ticket with no additional details
func TestStory1_1_OneLinePerTicket(t *testing.T) {
	// Given: I have tickets with status, dates, and descriptions
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "First Ticket",
			BinName:     "In Progress",
			Description: "This is a detailed description",
		},
		{
			ID:          "TICKET-002",
			Name:        "Second Ticket",
			BinName:     "To Do",
			Description: "Another description",
		},
	}

	// When: I format tickets in minimal mode
	output := FormatTicketsMinimal(tickets)

	// Then: Each ticket appears on exactly one line
	lines := strings.Split(output, "\n")

	// Should have header (1) + blank (1) + 2 tickets + trailing newline = 5 elements total
	// The blank line is an empty string, so we have: ["Found...", "", "[TICKET-001]...", "[TICKET-002]...", ""]
	if len(lines) != 5 {
		t.Errorf("Expected 5 lines total (header + blank + 2 tickets + trailing), got %d: %v", len(lines), lines)
	}

	// Verify the structure more explicitly
	if !strings.Contains(lines[0], "Found 2 ticket(s)") {
		t.Error("First line should be header")
	}
	if lines[1] != "" {
		t.Errorf("Second line should be blank, got: '%s'", lines[1])
	}
	if !strings.Contains(lines[2], "[TICKET-001]") {
		t.Error("Third line should be first ticket")
	}
	if !strings.Contains(lines[3], "[TICKET-002]") {
		t.Error("Fourth line should be second ticket")
	}
}

// TestStory1_1_ClearSeparation verifies ticket ID and name are clearly separated
// Acceptance Criterion 3: The ticket ID and name are clearly separated by formatting (brackets around ID, space before name)
func TestStory1_1_ClearSeparation(t *testing.T) {
	// Given: I have a ticket
	tickets := []models.Ticket{
		{ID: "TICKET-123", Name: "Test Ticket"},
	}

	// When: I format tickets in minimal mode
	output := FormatTicketsMinimal(tickets)

	// Then: The ID is wrapped in brackets and followed by a space and the name
	if !strings.Contains(output, "[TICKET-123] Test Ticket") {
		t.Errorf("Output should contain '[TICKET-123] Test Ticket', got: %s", output)
	}

	// Brackets should surround the ID
	if !strings.Contains(output, "[TICKET-123]") {
		t.Error("Ticket ID should be wrapped in brackets")
	}
}

// TestStory1_1_NoStatusDateDescription verifies minimal output contains no extra details
// Acceptance Criterion 4: The output contains no status information, no dates, and no descriptions
func TestStory1_1_NoStatusDateDescription(t *testing.T) {
	// Given: I have tickets with status, dates, and descriptions
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "Test Ticket",
			BinName:     "In Progress",
			Description: "Detailed description here",
		},
	}

	// When: I format tickets in minimal mode
	output := FormatTicketsMinimal(tickets)

	// Then: Output should NOT contain status, dates, or description
	if strings.Contains(output, "Status:") {
		t.Error("Minimal output should not contain 'Status:'")
	}

	if strings.Contains(output, "Created:") {
		t.Error("Minimal output should not contain 'Created:'")
	}

	if strings.Contains(output, "Updated:") {
		t.Error("Minimal output should not contain 'Updated:'")
	}

	if strings.Contains(output, "Due:") {
		t.Error("Minimal output should not contain 'Due:'")
	}

	if strings.Contains(output, "Description:") {
		t.Error("Minimal output should not contain 'Description:'")
	}

	if strings.Contains(output, "In Progress") {
		t.Error("Minimal output should not contain status value 'In Progress'")
	}

	if strings.Contains(output, "Detailed description here") {
		t.Error("Minimal output should not contain description text")
	}
}

// TestStory1_1_SixtyFourTicketsApproximately66Lines verifies line count for 64 tickets
// Acceptance Criterion 5: When I have 64 tickets, the output is approximately 66 lines (header + blank + 64 ticket lines)
func TestStory1_1_SixtyFourTicketsApproximately66Lines(t *testing.T) {
	// Given: I have 64 tickets assigned to me
	tickets := make([]models.Ticket, 64)
	for i := 0; i < 64; i++ {
		tickets[i] = models.Ticket{
			ID:   fmt.Sprintf("TICKET-%03d", i+1),
			Name: fmt.Sprintf("Ticket %d", i+1),
		}
	}

	// When: I format tickets in minimal mode
	output := FormatTicketsMinimal(tickets)

	// Then: The output is approximately 66 lines (header + blank + 64 ticket lines)
	lines := strings.Split(output, "\n")

	// Count non-empty lines (header + blank is counted as 2, + 64 tickets = 66)
	// But blank line will be empty string after split, so we count all lines before final newline
	lineCount := len(lines) - 1 // Remove the trailing empty string from final newline

	if lineCount < 65 || lineCount > 67 {
		t.Errorf("Expected approximately 66 lines for 64 tickets, got %d", lineCount)
	}
}

// TestStory1_1_AllTicketsShown verifies all assigned tickets are displayed
// Acceptance Criterion 6: All tickets assigned to me are shown in the list
func TestStory1_1_AllTicketsShown(t *testing.T) {
	// Given: I have multiple tickets with different IDs
	tickets := []models.Ticket{
		{ID: "TICKET-001", Name: "First"},
		{ID: "TICKET-002", Name: "Second"},
		{ID: "TICKET-003", Name: "Third"},
		{ID: "TICKET-004", Name: "Fourth"},
		{ID: "TICKET-005", Name: "Fifth"},
	}

	// When: I format tickets in minimal mode
	output := FormatTicketsMinimal(tickets)

	// Then: All ticket IDs appear in the output
	for _, ticket := range tickets {
		if !strings.Contains(output, ticket.ID) {
			t.Errorf("Output should contain ticket ID %s", ticket.ID)
		}
		if !strings.Contains(output, ticket.Name) {
			t.Errorf("Output should contain ticket name %s", ticket.Name)
		}
	}
}

// TestStory1_1_EmptyTicketList verifies handling of zero tickets
func TestStory1_1_EmptyTicketList(t *testing.T) {
	// Given: I have no tickets assigned to me
	tickets := []models.Ticket{}

	// When: I format tickets in minimal mode
	output := FormatTicketsMinimal(tickets)

	// Then: I see a clear "no tickets" message
	if !strings.Contains(output, "No tickets assigned to you.") {
		t.Errorf("Empty ticket list should show 'No tickets assigned to you.', got: %s", output)
	}
}

// TestStory1_1_MinimalFormatReducesScrolling verifies output is significantly shorter
func TestStory1_1_MinimalFormatReducesScrolling(t *testing.T) {
	// Given: I have 10 tickets
	tickets := make([]models.Ticket, 10)
	for i := 0; i < 10; i++ {
		tickets[i] = models.Ticket{
			ID:          fmt.Sprintf("TICKET-%03d", i+1),
			Name:        fmt.Sprintf("Ticket %d", i+1),
			BinName:     "In Progress",
			Description: "This is a description that would normally be displayed",
		}
	}

	// When: I format tickets in minimal mode vs verbose mode
	minimalOutput := FormatTicketsMinimal(tickets)
	verboseOutput := FormatTickets(tickets)

	// Then: Minimal output should be significantly shorter than verbose
	minimalLines := len(strings.Split(minimalOutput, "\n"))
	verboseLines := len(strings.Split(verboseOutput, "\n"))

	if minimalLines >= verboseLines {
		t.Errorf("Minimal output (%d lines) should be shorter than verbose output (%d lines)", minimalLines, verboseLines)
	}

	// Minimal should be at least 50% shorter
	if float64(minimalLines) > float64(verboseLines)*0.5 {
		t.Errorf("Minimal output should be at least 50%% shorter. Got minimal: %d, verbose: %d", minimalLines, verboseLines)
	}
}
