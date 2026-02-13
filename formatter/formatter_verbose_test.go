package formatter

import (
	"strings"
	"testing"
	"time"

	"github.com/Germanicus1/fb/models"
)

// Story 2.1: Verbose Mode Shows Detailed Output - Acceptance Tests

// TestStory2_1_VerboseShowsDetailedStructure verifies verbose output structure
// Acceptance Criterion 1: When I run `fb --verbose`, I see detailed output with this structure per ticket:
//   - Line 1: [ID] Name
//   - Line 2: Status: {status_id}
//   - Line 3: Created: {date} (if available)
//   - Line 4: Updated: {date} (if available)
//   - Line 5: Due: {date} (if available)
//   - Line 6+: Description: {text} (word-wrapped if long)
//   - Blank line between tickets
func TestStory2_1_VerboseShowsDetailedStructure(t *testing.T) {
	// Given: I have tickets with full details
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "First Ticket",
			BinName:     "In Progress",
			Description: "This is a description",
			CreatedAt:   time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2026, 2, 10, 14, 0, 0, 0, time.UTC),
		},
	}

	// When: I format in verbose mode
	output := FormatTickets(tickets)

	// Then: Output shows detailed structure
	lines := strings.Split(output, "\n")

	// Find the ticket header line
	headerLineIdx := -1
	for i, line := range lines {
		if strings.Contains(line, "[TICKET-001]") {
			headerLineIdx = i
			break
		}
	}

	if headerLineIdx == -1 {
		t.Fatal("Could not find ticket header in output")
	}

	// Verify structure elements are present
	if !strings.Contains(output, "[TICKET-001] First Ticket") {
		t.Error("Should contain ticket header: [ID] Name")
	}
	if !strings.Contains(output, "Status:") {
		t.Error("Should contain Status field")
	}
	if !strings.Contains(output, "Created:") {
		t.Error("Should contain Created field")
	}
	if !strings.Contains(output, "Updated:") {
		t.Error("Should contain Updated field")
	}
	if !strings.Contains(output, "Description:") {
		t.Error("Should contain Description field")
	}
}

// TestStory2_1_VerboseOutputIdenticalToCurrentDefault verifies backward compatibility
// Acceptance Criterion 2: The detailed output is identical to the current default `fb` output format
func TestStory2_1_VerboseOutputIdenticalToCurrentDefault(t *testing.T) {
	// Given: I have tickets
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "Test Ticket",
			BinName:     "To Do",
			Description: "Test description",
			CreatedAt:   time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC),
		},
	}

	// When: I format in verbose mode (using FormatTickets which is current default)
	output := FormatTickets(tickets)

	// Then: Output should have all the detailed fields that exist currently
	// This test verifies FormatTickets (verbose) maintains current behavior

	// Header line
	if !strings.Contains(output, "Found 1 ticket(s) assigned to you:") {
		t.Error("Should have header line")
	}

	// Ticket details
	if !strings.Contains(output, "[TICKET-001]") {
		t.Error("Should contain ticket ID")
	}
	if !strings.Contains(output, "Test Ticket") {
		t.Error("Should contain ticket name")
	}
	if !strings.Contains(output, "Status:") {
		t.Error("Should contain status field")
	}
	if !strings.Contains(output, "Description:") {
		t.Error("Should contain description field")
	}
}

// TestStory2_1_AllTicketFieldsDisplayed verifies all fields shown
// Acceptance Criterion 3: All ticket fields are displayed (status, dates, description)
func TestStory2_1_AllTicketFieldsDisplayed(t *testing.T) {
	// Given: I have a ticket with all fields populated
	tickets := []models.Ticket{
		{
			ID:          "TICKET-123",
			Name:        "Complete Ticket",
			BinName:     "In Progress",
			Description: "Full description text here",
			CreatedAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			DueDate:     time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	// When: I format in verbose mode
	output := FormatTickets(tickets)

	// Then: All fields are displayed
	if !strings.Contains(output, "TICKET-123") {
		t.Error("Should display ticket ID")
	}
	if !strings.Contains(output, "Complete Ticket") {
		t.Error("Should display ticket name")
	}
	if !strings.Contains(output, "Status:") && !strings.Contains(output, "In Progress") {
		t.Error("Should display status")
	}
	if !strings.Contains(output, "Created:") && !strings.Contains(output, "2026-01-01") {
		t.Error("Should display created date")
	}
	if !strings.Contains(output, "Updated:") && !strings.Contains(output, "2026-01-15") {
		t.Error("Should display updated date")
	}
	if !strings.Contains(output, "Due:") && !strings.Contains(output, "2026-03-01") {
		t.Error("Should display due date")
	}
	if !strings.Contains(output, "Description:") && !strings.Contains(output, "Full description") {
		t.Error("Should display description")
	}
}

// TestStory2_1_DescriptionsWordWrapped verifies word wrapping
// Acceptance Criterion 4: Descriptions are word-wrapped at appropriate width
func TestStory2_1_DescriptionsWordWrapped(t *testing.T) {
	// Given: I have a ticket with a long description
	longDesc := "This is a very long description that contains more than eighty characters and should be wrapped to multiple lines for better readability in the terminal"
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "Test",
			BinName:     "To Do",
			Description: longDesc,
		},
	}

	// When: I format in verbose mode
	output := FormatTickets(tickets)

	// Then: Description should be word-wrapped
	lines := strings.Split(output, "\n")

	// Find description lines (lines containing description text)
	descriptionLineCount := 0
	for _, line := range lines {
		if strings.Contains(line, "Description:") ||
		   (strings.HasPrefix(line, "    ") && !strings.Contains(line, ":")) {
			descriptionLineCount++
		}
	}

	if descriptionLineCount < 2 {
		t.Error("Long description should be wrapped across multiple lines")
	}

	// Verify no single line is excessively long
	for _, line := range lines {
		if len(line) > 100 {
			t.Errorf("Line too long (%d chars), wrapping may not be working: %s", len(line), line)
		}
	}
}

// TestStory2_1_SixtyFourTicketsApproximately400Lines verifies output length
// Acceptance Criterion 5: When I have 64 tickets, verbose output is approximately 400-500 lines (same as current)
func TestStory2_1_SixtyFourTicketsApproximately400Lines(t *testing.T) {
	// Given: I have 64 tickets with descriptions
	tickets := make([]models.Ticket, 64)
	for i := 0; i < 64; i++ {
		tickets[i] = models.Ticket{
			ID:          "TICKET-" + string(rune('A'+i%26)),
			Name:        "Ticket " + string(rune('A'+i%26)),
			BinName:     "In Progress",
			Description: "This is a description for the ticket",
			CreatedAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		}
	}

	// When: I format in verbose mode
	output := FormatTickets(tickets)

	// Then: Output should be approximately 400-500 lines
	lines := strings.Split(output, "\n")
	lineCount := len(lines) - 1 // Remove trailing empty

	if lineCount < 300 || lineCount > 600 {
		t.Errorf("Expected approximately 400-500 lines for 64 tickets, got %d", lineCount)
	}
}

// TestStory2_1_VerboseMuchLongerThanMinimal verifies verbose is detailed
func TestStory2_1_VerboseMuchLongerThanMinimal(t *testing.T) {
	// Given: I have 10 tickets
	tickets := make([]models.Ticket, 10)
	for i := 0; i < 10; i++ {
		tickets[i] = models.Ticket{
			ID:          "TICKET-00" + string(rune('0'+i)),
			Name:        "Ticket " + string(rune('0'+i)),
			BinName:     "To Do",
			Description: "Description for ticket",
			CreatedAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		}
	}

	// When: I format in both modes
	verboseOutput := FormatTickets(tickets)
	minimalOutput := FormatTicketsMinimal(tickets)

	// Then: Verbose should be significantly longer
	verboseLines := len(strings.Split(verboseOutput, "\n"))
	minimalLines := len(strings.Split(minimalOutput, "\n"))

	if verboseLines <= minimalLines {
		t.Errorf("Verbose output (%d lines) should be longer than minimal (%d lines)", verboseLines, minimalLines)
	}

	// Verbose should be at least 3x longer
	if verboseLines < minimalLines*3 {
		t.Errorf("Verbose should be at least 3x longer than minimal. Got verbose: %d, minimal: %d", verboseLines, minimalLines)
	}
}

// TestStory2_1_VerboseShowsStatusForAllTickets verifies status display
func TestStory2_1_VerboseShowsStatusForAllTickets(t *testing.T) {
	// Given: I have tickets with different statuses
	tickets := []models.Ticket{
		{ID: "TICKET-001", Name: "First", BinName: "To Do"},
		{ID: "TICKET-002", Name: "Second", BinName: "In Progress"},
		{ID: "TICKET-003", Name: "Third", BinName: "Done"},
	}

	// When: I format in verbose mode
	output := FormatTickets(tickets)

	// Then: Status is shown for all tickets
	statusCount := strings.Count(output, "Status:")
	if statusCount != 3 {
		t.Errorf("Expected 'Status:' to appear 3 times (once per ticket), got %d", statusCount)
	}

	// And: Status values are shown
	if !strings.Contains(output, "To Do") {
		t.Error("Should contain 'To Do' status")
	}
	if !strings.Contains(output, "In Progress") {
		t.Error("Should contain 'In Progress' status")
	}
	if !strings.Contains(output, "Done") {
		t.Error("Should contain 'Done' status")
	}
}

// TestStory2_1_VerboseShowsDescriptionForAllTickets verifies description display
func TestStory2_1_VerboseShowsDescriptionForAllTickets(t *testing.T) {
	// Given: I have tickets with descriptions
	tickets := []models.Ticket{
		{ID: "TICKET-001", Name: "First", BinName: "To Do", Description: "First description"},
		{ID: "TICKET-002", Name: "Second", BinName: "To Do", Description: "Second description"},
	}

	// When: I format in verbose mode
	output := FormatTickets(tickets)

	// Then: Description label appears for all tickets
	descCount := strings.Count(output, "Description:")
	if descCount != 2 {
		t.Errorf("Expected 'Description:' to appear 2 times, got %d", descCount)
	}

	// And: Description content is shown
	if !strings.Contains(output, "First description") {
		t.Error("Should contain first description text")
	}
	if !strings.Contains(output, "Second description") {
		t.Error("Should contain second description text")
	}
}

// TestStory2_1_VerboseTicketsSeparatedByBlankLine verifies separation
func TestStory2_1_VerboseTicketsSeparatedByBlankLine(t *testing.T) {
	// Given: I have multiple tickets
	tickets := []models.Ticket{
		{ID: "TICKET-001", Name: "First", BinName: "To Do"},
		{ID: "TICKET-002", Name: "Second", BinName: "To Do"},
		{ID: "TICKET-003", Name: "Third", BinName: "To Do"},
	}

	// When: I format in verbose mode
	output := FormatTickets(tickets)

	// Then: Tickets should be separated by blank lines
	lines := strings.Split(output, "\n")

	// Find ticket header lines
	ticketHeaderIndices := []int{}
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "[TICKET-") {
			ticketHeaderIndices = append(ticketHeaderIndices, i)
		}
	}

	// Verify second and third tickets have blank line before them
	if len(ticketHeaderIndices) >= 2 {
		secondTicketIdx := ticketHeaderIndices[1]
		if secondTicketIdx > 0 && strings.TrimSpace(lines[secondTicketIdx-1]) != "" {
			t.Error("Second ticket should have blank line before it")
		}
	}

	if len(ticketHeaderIndices) >= 3 {
		thirdTicketIdx := ticketHeaderIndices[2]
		if thirdTicketIdx > 0 && strings.TrimSpace(lines[thirdTicketIdx-1]) != "" {
			t.Error("Third ticket should have blank line before it")
		}
	}
}

// TestStory2_1_VerboseEmptyListHandling verifies empty list in verbose
func TestStory2_1_VerboseEmptyListHandling(t *testing.T) {
	// Given: I have no tickets
	tickets := []models.Ticket{}

	// When: I format in verbose mode
	output := FormatTickets(tickets)

	// Then: Shows same message as minimal mode
	expectedMessage := "No tickets assigned to you."
	if output != expectedMessage {
		t.Errorf("Verbose mode should show '%s' for empty list, got: %s", expectedMessage, output)
	}
}
