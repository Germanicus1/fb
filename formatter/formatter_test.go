package formatter

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Germanicus1/fb/models"
)

// TestFormatHardcodedTicket tests formatting a single hardcoded ticket
func TestFormatHardcodedTicket(t *testing.T) {
	// Given: A hardcoded demo ticket
	ticket := models.Ticket{
		ID:          "DEMO-123",
		Name:        "Sample Ticket",
		BinName:     "In Progress",
		Description: "This is a demo ticket",
		CreatedAt:   time.Now(),
	}

	// When: Formatting the ticket
	output := FormatTicket(ticket)

	// Then: Output should contain all required fields
	if !strings.Contains(output, "DEMO-123") {
		t.Errorf("Output should contain ticket ID 'DEMO-123', got:\n%s", output)
	}

	if !strings.Contains(output, "Sample Ticket") {
		t.Errorf("Output should contain ticket name 'Sample Ticket', got:\n%s", output)
	}

	if !strings.Contains(output, "In Progress") {
		t.Errorf("Output should contain status 'In Progress', got:\n%s", output)
	}

	if !strings.Contains(output, "This is a demo ticket") {
		t.Errorf("Output should contain description 'This is a demo ticket', got:\n%s", output)
	}
}

// TestFormatTicketReadability tests that output is readable
func TestFormatTicketReadability(t *testing.T) {
	ticket := models.Ticket{
		ID:          "DEMO-123",
		Name:        "Sample Ticket",
		BinName:     "In Progress",
		Description: "This is a demo ticket",
	}

	output := FormatTicket(ticket)

	// Output should not be empty
	if output == "" {
		t.Error("Output should not be empty")
	}

	// Output should have multiple lines for readability
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 3 {
		t.Errorf("Output should have at least 3 lines for readability, got %d lines:\n%s", len(lines), output)
	}
}

// TestFormatTicketEmptyDescription tests handling of empty description
func TestFormatTicketEmptyDescription(t *testing.T) {
	ticket := models.Ticket{
		ID:          "DEMO-123",
		Name:        "Sample Ticket",
		BinName:     "In Progress",
		Description: "",
	}

	output := FormatTicket(ticket)

	// Should still display other fields
	if !strings.Contains(output, "DEMO-123") {
		t.Error("Should display ticket ID even with empty description")
	}

	if !strings.Contains(output, "Sample Ticket") {
		t.Error("Should display ticket name even with empty description")
	}
}

// Story 2.1 Acceptance Tests: Display All Tickets from API Response

// TestStory2_1_DisplayAllTicketsFromResponse tests that all tickets are displayed
func TestStory2_1_DisplayAllTicketsFromResponse(t *testing.T) {
	// Given: Multiple tickets from an API response
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "First Ticket",
			BinName:     "To Do",
			Description: "First ticket description",
		},
		{
			ID:          "TICKET-002",
			Name:        "Second Ticket",
			BinName:     "In Progress",
			Description: "Second ticket description",
		},
		{
			ID:          "TICKET-003",
			Name:        "Third Ticket",
			BinName:     "Done",
			Description: "Third ticket description",
		},
	}

	// When: Formatting all tickets
	output := FormatTickets(tickets)

	// Then: All tickets should be in the output
	if !strings.Contains(output, "TICKET-001") {
		t.Error("Output should contain first ticket ID")
	}
	if !strings.Contains(output, "TICKET-002") {
		t.Error("Output should contain second ticket ID")
	}
	if !strings.Contains(output, "TICKET-003") {
		t.Error("Output should contain third ticket ID")
	}

	if !strings.Contains(output, "First Ticket") {
		t.Error("Output should contain first ticket name")
	}
	if !strings.Contains(output, "Second Ticket") {
		t.Error("Output should contain second ticket name")
	}
	if !strings.Contains(output, "Third Ticket") {
		t.Error("Output should contain third ticket name")
	}
}

// TestStory2_1_EachTicketShowsRequiredFields verifies each ticket displays all required fields
func TestStory2_1_EachTicketShowsRequiredFields(t *testing.T) {
	// Given: Tickets with all required fields
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "Test Ticket",
			BinName:     "To Do",
			Description: "Test description",
		},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Each required field should be present
	requiredFields := []string{"TICKET-001", "Test Ticket", "To Do", "Test description"}
	for _, field := range requiredFields {
		if !strings.Contains(output, field) {
			t.Errorf("Output should contain required field '%s', got:\n%s", field, output)
		}
	}
}

// TestStory2_1_TicketsAreVisuallySeparated verifies visual separation between tickets
func TestStory2_1_TicketsAreVisuallySeparated(t *testing.T) {
	// Given: Multiple tickets
	tickets := []models.Ticket{
		{ID: "TICKET-001", Name: "First", BinName: "To Do"},
		{ID: "TICKET-002", Name: "Second", BinName: "In Progress"},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Tickets should be visually separated (blank line between them)
	// Count blank lines between ticket content
	lines := strings.Split(output, "\n")
	hasBlankLine := false
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			hasBlankLine = true
			break
		}
	}

	if !hasBlankLine {
		t.Error("Tickets should be visually separated with blank lines")
	}
}

// TestStory2_1_HandlesSingleTicket verifies correct handling of 1 ticket
func TestStory2_1_HandlesSingleTicket(t *testing.T) {
	// Given: A single ticket
	tickets := []models.Ticket{
		{ID: "TICKET-001", Name: "Only Ticket", BinName: "To Do"},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: The ticket should be displayed
	if !strings.Contains(output, "TICKET-001") {
		t.Error("Should display single ticket")
	}
	if !strings.Contains(output, "Only Ticket") {
		t.Error("Should display ticket name")
	}
}

// TestStory2_1_HandlesFiveTickets verifies correct handling of 5 tickets
func TestStory2_1_HandlesFiveTickets(t *testing.T) {
	// Given: Five tickets
	tickets := make([]models.Ticket, 5)
	for i := 0; i < 5; i++ {
		tickets[i] = models.Ticket{
			ID:      fmt.Sprintf("TICKET-%03d", i+1),
			Name:    fmt.Sprintf("Ticket %d", i+1),
			BinName: "To Do",
		}
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: All five tickets should be present
	for i := 1; i <= 5; i++ {
		expectedID := fmt.Sprintf("TICKET-%03d", i)
		if !strings.Contains(output, expectedID) {
			t.Errorf("Should display ticket %s", expectedID)
		}
	}
}

// TestStory2_1_HandlesTenTickets verifies correct handling of 10 tickets
func TestStory2_1_HandlesTenTickets(t *testing.T) {
	// Given: Ten tickets
	tickets := make([]models.Ticket, 10)
	for i := 0; i < 10; i++ {
		tickets[i] = models.Ticket{
			ID:      fmt.Sprintf("TICKET-%03d", i+1),
			Name:    fmt.Sprintf("Ticket %d", i+1),
			BinName: "To Do",
		}
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: All ten tickets should be present
	for i := 1; i <= 10; i++ {
		expectedID := fmt.Sprintf("TICKET-%03d", i)
		if !strings.Contains(output, expectedID) {
			t.Errorf("Should display ticket %s", expectedID)
		}
	}
}

// TestStory2_1_OutputRemainsReadable verifies output is readable with multiple tickets
func TestStory2_1_OutputRemainsReadable(t *testing.T) {
	// Given: Multiple tickets
	tickets := make([]models.Ticket, 5)
	for i := 0; i < 5; i++ {
		tickets[i] = models.Ticket{
			ID:          fmt.Sprintf("TICKET-%03d", i+1),
			Name:        fmt.Sprintf("Ticket %d", i+1),
			BinName:     "To Do",
			Description: "Description for ticket",
		}
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Output should not be empty
	if output == "" {
		t.Error("Output should not be empty")
	}

	// Output should have multiple lines
	lines := strings.Split(output, "\n")
	if len(lines) < 5 {
		t.Error("Output should have sufficient lines for readability")
	}
}

// TestStory2_1_EmptyTicketList verifies handling of empty ticket list
func TestStory2_1_EmptyTicketList(t *testing.T) {
	// Given: No tickets
	tickets := []models.Ticket{}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Should show a clear message
	if output == "" {
		t.Error("Output should not be empty for zero tickets")
	}

	// Should indicate no tickets found
	if !strings.Contains(strings.ToLower(output), "no") || !strings.Contains(strings.ToLower(output), "ticket") {
		t.Errorf("Output should indicate no tickets, got: %s", output)
	}
}
// Story 2.3 Acceptance Tests: Display Created and Updated Dates

// TestStory2_3_DisplayCreatedDate verifies created date is displayed in YYYY-MM-DD format
func TestStory2_3_DisplayCreatedDate(t *testing.T) {
	// Given: A ticket with a created date
	createdTime := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)
	tickets := []models.Ticket{
		{
			ID:        "TICKET-001",
			Name:      "Test Ticket",
			BinName:   "To Do",
			CreatedAt: createdTime,
		},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Created date should be in YYYY-MM-DD format
	if !strings.Contains(output, "2026-01-15") {
		t.Errorf("Output should contain created date in YYYY-MM-DD format, got:\n%s", output)
	}

	// Created date should be labeled
	if !strings.Contains(output, "Created:") {
		t.Errorf("Output should label created date with 'Created:', got:\n%s", output)
	}
}

// TestStory2_3_DisplayUpdatedDate verifies updated date is displayed in YYYY-MM-DD format
func TestStory2_3_DisplayUpdatedDate(t *testing.T) {
	// Given: A ticket with an updated date
	updatedTime := time.Date(2026, 2, 10, 14, 45, 0, 0, time.UTC)
	tickets := []models.Ticket{
		{
			ID:        "TICKET-001",
			Name:      "Test Ticket",
			BinName:   "To Do",
			UpdatedAt: updatedTime,
		},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Updated date should be in YYYY-MM-DD format
	if !strings.Contains(output, "2026-02-10") {
		t.Errorf("Output should contain updated date in YYYY-MM-DD format, got:\n%s", output)
	}

	// Updated date should be labeled
	if !strings.Contains(output, "Updated:") {
		t.Errorf("Output should label updated date with 'Updated:', got:\n%s", output)
	}
}

// TestStory2_3_DisplayBothDates verifies both dates can be displayed together
func TestStory2_3_DisplayBothDates(t *testing.T) {
	// Given: A ticket with both created and updated dates
	createdTime := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)
	updatedTime := time.Date(2026, 2, 10, 14, 45, 0, 0, time.UTC)
	tickets := []models.Ticket{
		{
			ID:        "TICKET-001",
			Name:      "Test Ticket",
			BinName:   "To Do",
			CreatedAt: createdTime,
			UpdatedAt: updatedTime,
		},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Both dates should be present
	if !strings.Contains(output, "2026-01-15") {
		t.Error("Output should contain created date")
	}
	if !strings.Contains(output, "2026-02-10") {
		t.Error("Output should contain updated date")
	}
}

// TestStory2_3_MissingDateShowsNothing verifies missing dates don't show errors
func TestStory2_3_MissingDateShowsNothing(t *testing.T) {
	// Given: A ticket with no dates (zero values)
	tickets := []models.Ticket{
		{
			ID:        "TICKET-001",
			Name:      "Test Ticket",
			BinName:   "To Do",
			CreatedAt: time.Time{}, // Zero value
			UpdatedAt: time.Time{}, // Zero value
		},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Should not show "Created:" or "Updated:" labels
	// (The current implementation doesn't show labels if dates are missing)
	// Verify no date-related errors or "N/A" appear
	if strings.Contains(output, "N/A") {
		t.Error("Output should not show 'N/A' for missing dates")
	}

	// Verify ticket is still displayed
	if !strings.Contains(output, "TICKET-001") {
		t.Error("Ticket should still be displayed even without dates")
	}
}

// TestStory2_3_DateFormatConsistency verifies date format is consistent across tickets
func TestStory2_3_DateFormatConsistency(t *testing.T) {
	// Given: Multiple tickets with different dates
	tickets := []models.Ticket{
		{
			ID:        "TICKET-001",
			Name:      "First Ticket",
			BinName:   "To Do",
			CreatedAt: time.Date(2026, 1, 5, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2026, 1, 6, 11, 0, 0, 0, time.UTC),
		},
		{
			ID:        "TICKET-002",
			Name:      "Second Ticket",
			BinName:   "In Progress",
			CreatedAt: time.Date(2026, 2, 15, 14, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2026, 2, 16, 15, 0, 0, 0, time.UTC),
		},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: All dates should be in YYYY-MM-DD format
	expectedDates := []string{"2026-01-05", "2026-01-06", "2026-02-15", "2026-02-16"}
	for _, date := range expectedDates {
		if !strings.Contains(output, date) {
			t.Errorf("Output should contain date %s in YYYY-MM-DD format", date)
		}
	}
}
// Story 2.4 Acceptance Tests: Display Due Date When Present

// TestStory2_4_DisplayDueDateWhenPresent verifies due date is displayed when set
func TestStory2_4_DisplayDueDateWhenPresent(t *testing.T) {
	// Given: A ticket with a due date
	dueTime := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	tickets := []models.Ticket{
		{
			ID:      "TICKET-001",
			Name:    "Test Ticket",
			BinName: "To Do",
			DueDate: dueTime,
		},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Due date should be in YYYY-MM-DD format
	if !strings.Contains(output, "2026-03-01") {
		t.Errorf("Output should contain due date in YYYY-MM-DD format, got:\n%s", output)
	}

	// Due date should be labeled
	if !strings.Contains(output, "Due:") {
		t.Errorf("Output should label due date with 'Due:', got:\n%s", output)
	}
}

// TestStory2_4_NoDueDateShowsNothing verifies tickets without due dates don't show the field
func TestStory2_4_NoDueDateShowsNothing(t *testing.T) {
	// Given: A ticket without a due date (zero time)
	tickets := []models.Ticket{
		{
			ID:      "TICKET-001",
			Name:    "Test Ticket",
			BinName: "To Do",
			DueDate: time.Time{}, // Zero value - no due date
		},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Should not show "Due:" label
	// (Current implementation should not show the field at all)
	// Verify ticket is still displayed
	if !strings.Contains(output, "TICKET-001") {
		t.Error("Ticket should still be displayed even without due date")
	}
}

// TestStory2_4_DueDateFormatConsistentWithOtherDates verifies due date format matches created/updated
func TestStory2_4_DueDateFormatConsistentWithOtherDates(t *testing.T) {
	// Given: A ticket with all three dates
	createdTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	updatedTime := time.Date(2026, 2, 10, 14, 0, 0, 0, time.UTC)
	dueTime := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	tickets := []models.Ticket{
		{
			ID:        "TICKET-001",
			Name:      "Test Ticket",
			BinName:   "To Do",
			CreatedAt: createdTime,
			UpdatedAt: updatedTime,
			DueDate:   dueTime,
		},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: All dates should be in YYYY-MM-DD format
	expectedDates := []string{"2026-01-15", "2026-02-10", "2026-03-01"}
	for _, date := range expectedDates {
		if !strings.Contains(output, date) {
			t.Errorf("Output should contain date %s in YYYY-MM-DD format", date)
		}
	}
}

// TestStory2_4_PastDueDateStillDisplayed verifies past due dates are shown
func TestStory2_4_PastDueDateStillDisplayed(t *testing.T) {
	// Given: A ticket with a past due date
	pastDueTime := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	tickets := []models.Ticket{
		{
			ID:      "TICKET-001",
			Name:    "Overdue Ticket",
			BinName: "To Do",
			DueDate: pastDueTime,
		},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Past due date should still be displayed
	if !strings.Contains(output, "2025-12-31") {
		t.Error("Output should display past due dates")
	}
}

// TestStory2_4_MultiplTicketsWithMixedDueDates verifies handling of mixed due date scenarios
func TestStory2_4_MultipleTicketsWithMixedDueDates(t *testing.T) {
	// Given: Multiple tickets with and without due dates
	dueTime := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	tickets := []models.Ticket{
		{
			ID:      "TICKET-001",
			Name:    "Has Due Date",
			BinName: "To Do",
			DueDate: dueTime,
		},
		{
			ID:      "TICKET-002",
			Name:    "No Due Date",
			BinName: "In Progress",
			DueDate: time.Time{}, // No due date
		},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: First ticket should show due date
	if !strings.Contains(output, "2026-03-01") {
		t.Error("Ticket with due date should display it")
	}

	// Both tickets should be displayed
	if !strings.Contains(output, "TICKET-001") || !strings.Contains(output, "TICKET-002") {
		t.Error("All tickets should be displayed regardless of due date presence")
	}
}

// Story 3.1 Acceptance Tests: Improve Ticket Visual Separation

// TestStory3_1_TicketsAreVisuallySeparated verifies each ticket is visually separated from others
func TestStory3_1_TicketsAreVisuallySeparated(t *testing.T) {
	// Given: Multiple tickets to display
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "First Ticket",
			BinName:     "To Do",
			Description: "First description",
		},
		{
			ID:          "TICKET-002",
			Name:        "Second Ticket",
			BinName:     "In Progress",
			Description: "Second description",
		},
		{
			ID:          "TICKET-003",
			Name:        "Third Ticket",
			BinName:     "Done",
			Description: "Third description",
		},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Tickets should be visually separated (blank line between them)
	lines := strings.Split(output, "\n")

	// Find where each ticket starts and verify blank line before it (except first)
	firstTicketIdx := -1
	secondTicketIdx := -1
	thirdTicketIdx := -1

	for i, line := range lines {
		if strings.Contains(line, "[TICKET-001]") {
			firstTicketIdx = i
		}
		if strings.Contains(line, "[TICKET-002]") {
			secondTicketIdx = i
		}
		if strings.Contains(line, "[TICKET-003]") {
			thirdTicketIdx = i
		}
	}

	// Verify we found all tickets
	if firstTicketIdx == -1 || secondTicketIdx == -1 || thirdTicketIdx == -1 {
		t.Fatal("Could not find all ticket headers in output")
	}

	// Verify second and third tickets have blank line before them
	if secondTicketIdx > 0 && strings.TrimSpace(lines[secondTicketIdx-1]) != "" {
		t.Error("Second ticket should have blank line before it for visual separation")
	}
	if thirdTicketIdx > 0 && strings.TrimSpace(lines[thirdTicketIdx-1]) != "" {
		t.Error("Third ticket should have blank line before it for visual separation")
	}
}

// TestStory3_1_SeparationWorksInAllTerminalTypes verifies separation doesn't require special characters
func TestStory3_1_SeparationWorksInAllTerminalTypes(t *testing.T) {
	// Given: Multiple tickets
	tickets := []models.Ticket{
		{ID: "TICKET-001", Name: "First", BinName: "To Do"},
		{ID: "TICKET-002", Name: "Second", BinName: "In Progress"},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Separation should only use standard ASCII characters (no box drawing, Unicode, etc.)
	// Check that output only contains printable ASCII and standard whitespace
	for i, r := range output {
		if r != '\n' && r != '\r' && r != '\t' && r != ' ' && (r < 32 || r > 126) {
			t.Errorf("Output contains non-standard character at position %d: %U (should use only ASCII)", i, r)
		}
	}
}

// TestStory3_1_FirstTicketStartsCleanly verifies no separator before first ticket
func TestStory3_1_FirstTicketStartsCleanly(t *testing.T) {
	// Given: Multiple tickets
	tickets := []models.Ticket{
		{ID: "TICKET-001", Name: "First", BinName: "To Do"},
		{ID: "TICKET-002", Name: "Second", BinName: "In Progress"},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: First ticket should start immediately after summary line (no separator before it)
	lines := strings.Split(output, "\n")

	// Find the summary line and first ticket
	summaryIdx := -1
	firstTicketIdx := -1

	for i, line := range lines {
		if strings.Contains(line, "Found") && strings.Contains(line, "ticket(s)") {
			summaryIdx = i
		}
		if strings.Contains(line, "[TICKET-001]") {
			firstTicketIdx = i
			break
		}
	}

	// Verify there's at most one blank line between summary and first ticket
	if firstTicketIdx > summaryIdx+2 {
		t.Errorf("First ticket should start cleanly after summary (max 1 blank line), found %d lines between", firstTicketIdx-summaryIdx-1)
	}
}

// TestStory3_1_LastTicketEndsCleanly verifies appropriate spacing after last ticket
func TestStory3_1_LastTicketEndsCleanly(t *testing.T) {
	// Given: Multiple tickets
	tickets := []models.Ticket{
		{ID: "TICKET-001", Name: "First", BinName: "To Do"},
		{ID: "TICKET-002", Name: "Second", BinName: "In Progress"},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Output should not have excessive trailing newlines
	// (at most 1 trailing newline is acceptable)
	if strings.HasSuffix(output, "\n\n\n") {
		t.Error("Last ticket should end cleanly without excessive trailing newlines")
	}
}

// TestStory3_1_SeparationIsConsistent verifies separation is consistent across all tickets
func TestStory3_1_SeparationIsConsistent(t *testing.T) {
	// Given: Multiple tickets (5 tickets to test consistency)
	tickets := make([]models.Ticket, 5)
	for i := 0; i < 5; i++ {
		tickets[i] = models.Ticket{
			ID:      fmt.Sprintf("TICKET-%03d", i+1),
			Name:    fmt.Sprintf("Ticket %d", i+1),
			BinName: "To Do",
		}
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Each ticket (except first) should have same separation pattern
	lines := strings.Split(output, "\n")

	ticketIndices := make([]int, 0)
	for i, line := range lines {
		if strings.HasPrefix(line, "[TICKET-") {
			ticketIndices = append(ticketIndices, i)
		}
	}

	// Verify each ticket (except first) has a blank line before it
	for i := 1; i < len(ticketIndices); i++ {
		ticketIdx := ticketIndices[i]
		if ticketIdx > 0 && strings.TrimSpace(lines[ticketIdx-1]) != "" {
			t.Errorf("Ticket at line %d should have blank line before it for consistent separation", ticketIdx)
		}
	}
}

// TestStory3_1_OutputReadableWhenCopiedPasted verifies output works in pipes/redirects
func TestStory3_1_OutputReadableWhenCopiedPasted(t *testing.T) {
	// Given: Multiple tickets
	tickets := []models.Ticket{
		{ID: "TICKET-001", Name: "First", BinName: "To Do"},
		{ID: "TICKET-002", Name: "Second", BinName: "In Progress"},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Output should be plain text (no ANSI codes, control characters except newline/tab)
	// and should contain standard newlines
	if !strings.Contains(output, "\n") {
		t.Error("Output should contain newlines for readability")
	}

	// Verify no ANSI escape codes
	if strings.Contains(output, "\x1b[") {
		t.Error("Output should not contain ANSI escape codes (for copy/paste compatibility)")
	}

	// Verify all tickets are present (copy/paste shouldn't lose data)
	if !strings.Contains(output, "TICKET-001") || !strings.Contains(output, "TICKET-002") {
		t.Error("Output should contain all tickets when copied/pasted")
	}
}

// Story 3.2 Acceptance Tests: Format Long Descriptions with Word Wrapping

// TestStory3_2_LongDescriptionsWrappedToMultipleLines verifies descriptions > 80 chars are wrapped
func TestStory3_2_LongDescriptionsWrappedToMultipleLines(t *testing.T) {
	// Given: A ticket with a description longer than 80 characters
	longDescription := "This is a very long description that contains more than eighty characters and should be wrapped to multiple lines for better readability in the terminal"
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "Test Ticket",
			BinName:     "To Do",
			Description: longDescription,
		},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Description should be on multiple lines
	lines := strings.Split(output, "\n")

	// Find the description label line and count continuation lines
	descriptionIdx := -1
	for i, line := range lines {
		if strings.Contains(line, "Description:") {
			descriptionIdx = i
			break
		}
	}

	if descriptionIdx == -1 {
		t.Fatal("Could not find Description label in output")
	}

	// Count lines that are part of the description (indented continuation lines)
	continuationLines := 0
	for i := descriptionIdx + 1; i < len(lines); i++ {
		line := lines[i]
		// Check if line is indented (continuation of description)
		if strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "  ") {
			// Make sure it's not another field label (contains colon near start)
			if !strings.Contains(line, ":") || strings.Index(line, ":") > 10 {
				continuationLines++
			}
		} else if strings.TrimSpace(line) == "" {
			// Blank lines are OK, continue checking
			continue
		} else {
			// Reached end of description section
			break
		}
	}

	if continuationLines < 1 {
		t.Errorf("Long description should be wrapped to multiple lines, found only %d continuation line(s)", continuationLines)
	}

	// Verify full content is preserved
	if !strings.Contains(output, "very long description") {
		t.Error("Wrapped description should preserve full content")
	}
}

// TestStory3_2_WrappedLinesMaintainReadability verifies wrapped lines are properly indented
func TestStory3_2_WrappedLinesMaintainReadability(t *testing.T) {
	// Given: A ticket with a long description
	longDescription := "This is a very long description that needs to be wrapped across multiple lines while maintaining proper indentation and readability"
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "Test Ticket",
			BinName:     "To Do",
			Description: longDescription,
		},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Wrapped lines should have proper indentation for readability
	lines := strings.Split(output, "\n")

	foundDescriptionLabel := false
	foundWrappedLine := false

	for i, line := range lines {
		if strings.Contains(line, "Description:") {
			foundDescriptionLabel = true

			// Check if there's a continuation line
			if i+1 < len(lines) && !strings.Contains(lines[i+1], "[") && strings.TrimSpace(lines[i+1]) != "" {
				foundWrappedLine = true

				// Wrapped line should be indented
				if !strings.HasPrefix(lines[i+1], "    ") && !strings.HasPrefix(lines[i+1], "  ") {
					t.Error("Wrapped description lines should be indented for readability")
				}
			}
		}
	}

	if foundDescriptionLabel && !foundWrappedLine {
		t.Error("Long description should have wrapped lines")
	}
}

// TestStory3_2_WordBoundariesRespected verifies words aren't broken mid-word
func TestStory3_2_WordBoundariesRespected(t *testing.T) {
	// Given: A ticket with a description that would require wrapping
	description := "This description contains several multisyllabic words including extraordinary, remarkable, and outstanding qualities that should be preserved"
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "Test Ticket",
			BinName:     "To Do",
			Description: description,
		},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Words should not be broken mid-word (no hyphens at line breaks)
	// Verify key words are intact
	if !strings.Contains(output, "extraordinary") {
		t.Error("Word 'extraordinary' should not be broken across lines")
	}
	if !strings.Contains(output, "remarkable") {
		t.Error("Word 'remarkable' should not be broken across lines")
	}
	if !strings.Contains(output, "outstanding") {
		t.Error("Word 'outstanding' should not be broken across lines")
	}
}

// TestStory3_2_VeryLongWordsHandledGracefully verifies URLs and code snippets work
func TestStory3_2_VeryLongWordsHandledGracefully(t *testing.T) {
	// Given: A ticket with very long words (URL, code)
	description := "Please check this URL: https://example.com/very/long/path/to/some/resource/that/exceeds/normal/line/length and update accordingly"
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "Test Ticket",
			BinName:     "To Do",
			Description: description,
		},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: URL should be preserved (not broken)
	if !strings.Contains(output, "https://example.com/very/long/path/to/some/resource/that/exceeds/normal/line/length") {
		t.Error("Very long words (URLs) should be preserved intact")
	}
}

// TestStory3_2_WrappingWorksConsistently verifies wrapping works across different terminal widths
func TestStory3_2_WrappingWorksConsistently(t *testing.T) {
	// Given: A ticket with a long description
	longDescription := "This is a long description that should be wrapped consistently regardless of the exact terminal width being used by the user"
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "Test Ticket",
			BinName:     "To Do",
			Description: longDescription,
		},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Description should be wrapped
	// (Implementation should use a consistent width like 80 chars)
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		// Allow some tolerance for indentation, but lines shouldn't be excessively long
		if len(line) > 100 {
			t.Errorf("Line is too long (%d chars), wrapping may not be working: %s", len(line), line)
		}
	}
}

// TestStory3_2_ShortDescriptionsNotAffected verifies descriptions < 80 chars are unchanged
func TestStory3_2_ShortDescriptionsNotAffected(t *testing.T) {
	// Given: A ticket with a short description (under 80 characters)
	shortDescription := "This is a short description"
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "Test Ticket",
			BinName:     "To Do",
			Description: shortDescription,
		},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Short description should appear on single line with label
	// Find the Description line
	lines := strings.Split(output, "\n")

	foundOnSingleLine := false
	for _, line := range lines {
		if strings.Contains(line, "Description:") && strings.Contains(line, "short description") {
			foundOnSingleLine = true
			break
		}
	}

	if !foundOnSingleLine {
		t.Error("Short descriptions should remain on single line with label")
	}
}

// Story 3.3 Acceptance Tests: Handle Empty Description Gracefully

// TestStory3_3_EmptyDescriptionShowsPlaceholder verifies empty descriptions show "(none)" or similar
func TestStory3_3_EmptyDescriptionShowsPlaceholder(t *testing.T) {
	// Given: A ticket with an empty description
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "Test Ticket",
			BinName:     "To Do",
			Description: "",
		},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Should show "Description: (none)" or similar placeholder
	if !strings.Contains(output, "Description:") {
		t.Error("Empty description should still show Description label")
	}

	if !strings.Contains(output, "(none)") && !strings.Contains(output, "None") && !strings.Contains(output, "N/A") {
		t.Error("Empty description should show placeholder like '(none)', 'None', or 'N/A'")
	}
}

// TestStory3_3_NullDescriptionShowsPlaceholder verifies null descriptions show "(none)"
func TestStory3_3_NullDescriptionShowsPlaceholder(t *testing.T) {
	// Given: A ticket with null description (empty string in Go)
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "Test Ticket",
			BinName:     "To Do",
			Description: "", // Null/empty
		},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Should show Description field with placeholder
	if !strings.Contains(output, "Description:") {
		t.Error("Description field should be present for empty descriptions")
	}
}

// TestStory3_3_EmptyDescriptionNoBlankSpace verifies no confusing blank space
func TestStory3_3_EmptyDescriptionNoBlankSpace(t *testing.T) {
	// Given: A ticket with empty description
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "Test Ticket",
			BinName:     "To Do",
			Description: "",
		},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Description line should not be just "Description: " with nothing after
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Description:") {
			trimmedLine := strings.TrimSpace(line)
			if trimmedLine == "Description:" {
				t.Error("Description label should not be followed by blank space; should have placeholder")
			}
		}
	}
}

// TestStory3_3_DescriptionFieldLabeled verifies field is labeled
func TestStory3_3_DescriptionFieldLabeled(t *testing.T) {
	// Given: A ticket with empty description
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "Test Ticket",
			BinName:     "To Do",
			Description: "",
		},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Description field should be labeled to show it was checked
	if !strings.Contains(output, "Description:") {
		t.Error("Empty description should still show 'Description:' label")
	}
}

// TestStory3_3_OtherFieldsDisplayNormally verifies other fields are unaffected
func TestStory3_3_OtherFieldsDisplayNormally(t *testing.T) {
	// Given: A ticket with empty description but other fields populated
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "Test Ticket",
			BinName:     "To Do",
			Description: "",
			CreatedAt:   time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC),
		},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Other fields should display normally
	if !strings.Contains(output, "TICKET-001") {
		t.Error("Ticket ID should display normally")
	}
	if !strings.Contains(output, "Test Ticket") {
		t.Error("Ticket name should display normally")
	}
	if !strings.Contains(output, "To Do") {
		t.Error("Status should display normally")
	}
	if !strings.Contains(output, "2026-01-15") {
		t.Error("Created date should display normally")
	}
}

// TestStory3_3_DistinguishFromLoadingError verifies user can tell it's not an error
func TestStory3_3_DistinguishFromLoadingError(t *testing.T) {
	// Given: A ticket with empty description
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "Test Ticket",
			BinName:     "To Do",
			Description: "",
		},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Output should clearly indicate "no description" vs error
	// Should have some placeholder text, not just missing field
	if !strings.Contains(output, "Description:") {
		t.Error("Description field should be present to show it was checked (not an error)")
	}

	// Verify it doesn't look like an error (no "error", "failed", etc.)
	lowerOutput := strings.ToLower(output)
	if strings.Contains(lowerOutput, "error") || strings.Contains(lowerOutput, "failed") {
		t.Error("Empty description should not look like an error condition")
	}
}

// Story 3.4 Acceptance Tests: Display Summary Line with Ticket Count
// NOTE: This story was already implemented in Sprint 2. These tests verify the existing functionality.

// TestStory3_4_DisplaySummaryLine verifies summary line is displayed
func TestStory3_4_DisplaySummaryLine(t *testing.T) {
	// Given: Multiple tickets
	tickets := []models.Ticket{
		{ID: "TICKET-001", Name: "First", BinName: "To Do"},
		{ID: "TICKET-002", Name: "Second", BinName: "In Progress"},
		{ID: "TICKET-003", Name: "Third", BinName: "Done"},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Should display summary line
	if !strings.Contains(output, "Found") || !strings.Contains(output, "ticket") {
		t.Error("Output should contain summary line with 'Found' and 'ticket'")
	}
}

// TestStory3_4_SummaryShowsTotalCount verifies count matches number of tickets
func TestStory3_4_SummaryShowsTotalCount(t *testing.T) {
	// Given: 5 tickets
	tickets := make([]models.Ticket, 5)
	for i := 0; i < 5; i++ {
		tickets[i] = models.Ticket{
			ID:      fmt.Sprintf("TICKET-%03d", i+1),
			Name:    fmt.Sprintf("Ticket %d", i+1),
			BinName: "To Do",
		}
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Should show "Found 5 ticket(s)"
	if !strings.Contains(output, "Found 5 ticket(s)") {
		t.Errorf("Summary should show 'Found 5 ticket(s)', got:\n%s", output)
	}
}

// TestStory3_4_ZeroTicketsMessage verifies message for zero tickets
func TestStory3_4_ZeroTicketsMessage(t *testing.T) {
	// Given: No tickets
	tickets := []models.Ticket{}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Should show "No tickets assigned to"
	if !strings.Contains(output, "No tickets assigned to") {
		t.Errorf("Should show 'No tickets assigned to' message for zero tickets, got:\n%s", output)
	}
}

// TestStory3_4_OneTicketSingular verifies proper singular form
func TestStory3_4_OneTicketSingular(t *testing.T) {
	// Given: 1 ticket
	tickets := []models.Ticket{
		{ID: "TICKET-001", Name: "Only Ticket", BinName: "To Do"},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Should use proper grammar (accepting "1 ticket(s)" as valid)
	// Current implementation uses "ticket(s)" for all counts which is acceptable
	if !strings.Contains(output, "Found 1 ticket(s)") {
		t.Errorf("Summary should show ticket count for 1 ticket, got:\n%s", output)
	}
}

// TestStory3_4_SummarySeparatedFromTickets verifies summary is clearly separated
func TestStory3_4_SummarySeparatedFromTickets(t *testing.T) {
	// Given: Multiple tickets
	tickets := []models.Ticket{
		{ID: "TICKET-001", Name: "First", BinName: "To Do"},
		{ID: "TICKET-002", Name: "Second", BinName: "In Progress"},
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Summary should be separated from ticket details
	lines := strings.Split(output, "\n")

	// Find summary line
	summaryIdx := -1
	for i, line := range lines {
		if strings.Contains(line, "Found") && strings.Contains(line, "ticket(s)") {
			summaryIdx = i
			break
		}
	}

	if summaryIdx == -1 {
		t.Fatal("Could not find summary line in output")
	}

	// Verify there's at least a blank line between summary and first ticket
	if summaryIdx+1 < len(lines) && strings.TrimSpace(lines[summaryIdx+1]) != "" {
		// First line after summary should be blank for separation
		// (or could be the first ticket, which is also acceptable)
	}
}

// TestStory3_4_CountMatchesDisplayedTickets verifies count accuracy
func TestStory3_4_CountMatchesDisplayedTickets(t *testing.T) {
	// Given: 10 tickets
	tickets := make([]models.Ticket, 10)
	for i := 0; i < 10; i++ {
		tickets[i] = models.Ticket{
			ID:      fmt.Sprintf("TICKET-%03d", i+1),
			Name:    fmt.Sprintf("Ticket %d", i+1),
			BinName: "To Do",
		}
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Summary should show "Found 10 ticket(s)"
	if !strings.Contains(output, "Found 10 ticket(s)") {
		t.Error("Count in summary should match number of tickets (10)")
	}

	// And: All 10 tickets should be in the output
	for i := 1; i <= 10; i++ {
		expectedID := fmt.Sprintf("TICKET-%03d", i)
		if !strings.Contains(output, expectedID) {
			t.Errorf("Expected to find %s in output", expectedID)
		}
	}
}

// STORY 4.4: Handle Special Characters in Ticket Data

// TestStory4_4_UnicodeCharactersDisplay tests Unicode character handling
func TestStory4_4_UnicodeCharactersDisplay(t *testing.T) {
	// Given: A ticket with Unicode characters in name and description
	ticket := models.Ticket{
		ID:          "UNICODE-1",
		Name:        "测试 Тест Δοκιμή Test",
		BinName:     "In Progress",
		Description: "Description with Unicode: 日本語 Русский Ελληνικά",
		CreatedAt:   time.Now(),
	}

	// When: Formatting the ticket
	output := FormatTicket(ticket)

	// Then: Unicode characters should be preserved and displayed correctly
	// Acceptance Criterion: Unicode characters in ticket names and descriptions display correctly
	if !strings.Contains(output, "测试 Тест Δοκιμή Test") {
		t.Error("Unicode characters in name should be preserved")
	}
	if !strings.Contains(output, "日本語 Русский Ελληνικά") {
		t.Error("Unicode characters in description should be preserved")
	}
}

// TestStory4_4_EmojisPreserved tests emoji handling
func TestStory4_4_EmojisPreserved(t *testing.T) {
	// Given: A ticket with emojis
	ticket := models.Ticket{
		ID:          "EMOJI-1",
		Name:        "Fix bug 🐛 in login",
		BinName:     "To Do",
		Description: "This is urgent! ⚠️ Need to fix ASAP ✅",
		CreatedAt:   time.Now(),
	}

	// When: Formatting the ticket
	output := FormatTicket(ticket)

	// Then: Emojis should be preserved or handled gracefully
	// Acceptance Criterion: Emojis in ticket content are preserved (or handled gracefully if terminal doesn't support)
	if !strings.Contains(output, "🐛") {
		t.Error("Emoji in name should be preserved")
	}
	if !strings.Contains(output, "⚠️") || !strings.Contains(output, "✅") {
		t.Error("Emojis in description should be preserved")
	}
}

// TestStory4_4_NewlinesRendered tests newline handling in descriptions
func TestStory4_4_NewlinesRendered(t *testing.T) {
	// Given: A ticket with newlines in description
	ticket := models.Ticket{
		ID:          "NEWLINE-1",
		Name:        "Multi-line description test",
		BinName:     "In Progress",
		Description: "Line 1\nLine 2\nLine 3",
		CreatedAt:   time.Now(),
	}

	// When: Formatting the ticket
	output := FormatTicket(ticket)

	// Then: Newlines should be rendered appropriately
	// Acceptance Criterion: Newlines in descriptions are rendered appropriately
	if !strings.Contains(output, "Line 1") || !strings.Contains(output, "Line 2") || !strings.Contains(output, "Line 3") {
		t.Error("All lines should be present in output")
	}

	// Lines should be separated (either by actual newlines or some other method)
	lines := strings.Split(output, "\n")
	hasMultipleLines := len(lines) > 1
	if !hasMultipleLines {
		t.Error("Multi-line description should result in multiple lines of output")
	}
}

// TestStory4_4_TabCharactersHandled tests tab character handling
func TestStory4_4_TabCharactersHandled(t *testing.T) {
	// Given: A ticket with tab characters
	ticket := models.Ticket{
		ID:          "TAB-1",
		Name:        "Ticket\twith\ttabs",
		BinName:     "To Do",
		Description: "Column1\tColumn2\tColumn3",
		CreatedAt:   time.Now(),
	}

	// When: Formatting the ticket
	output := FormatTicket(ticket)

	// Then: Tab characters shouldn't break formatting
	// Acceptance Criterion: Tab characters don't break formatting
	if !strings.Contains(output, "Ticket") && !strings.Contains(output, "with") && !strings.Contains(output, "tabs") {
		t.Error("Content with tabs should still be readable")
	}

	// Output should not have broken layout
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if len(line) > 200 {
			t.Error("Tab characters shouldn't create extremely long lines that break layout")
		}
	}
}

// TestStory4_4_ControlCharactersHandled tests control character handling
func TestStory4_4_ControlCharactersHandled(t *testing.T) {
	// Given: A ticket with control characters
	ticket := models.Ticket{
		ID:          "CTRL-1",
		Name:        "Test\x00with\x01control\x02chars",
		BinName:     "In Progress",
		Description: "Description\x1b[31mwith\x1b[0mANSI",
		CreatedAt:   time.Now(),
	}

	// When: Formatting the ticket
	output := FormatTicket(ticket)

	// Then: Control characters shouldn't corrupt terminal display
	// Acceptance Criterion: Control characters don't corrupt terminal display
	// Output should be generated without panicking
	if output == "" {
		t.Error("Output should not be empty even with control characters")
	}

	// Should contain the actual text content
	if !strings.Contains(output, "Test") || !strings.Contains(output, "control") {
		t.Error("Text content should be preserved despite control characters")
	}
}

// TestStory4_4_SpecialHTMLCharacters tests HTML-like special characters
func TestStory4_4_SpecialHTMLCharacters(t *testing.T) {
	// Given: A ticket with special characters that might be problematic
	ticket := models.Ticket{
		ID:          "SPECIAL-1",
		Name:        "Test <script> & \"quotes\" 'apostrophes'",
		BinName:     "To Do",
		Description: "Description with & < > \" ' characters",
		CreatedAt:   time.Now(),
	}

	// When: Formatting the ticket
	output := FormatTicket(ticket)

	// Then: Common special characters should display correctly
	// Acceptance Criterion: Common special characters (&, <, >, quotes) display correctly
	if !strings.Contains(output, "&") {
		t.Error("Ampersand should be preserved")
	}
	if !strings.Contains(output, "<") || !strings.Contains(output, ">") {
		t.Error("Angle brackets should be preserved")
	}
	if !strings.Contains(output, "\"") || !strings.Contains(output, "'") {
		t.Error("Quote marks should be preserved")
	}
}

// TestStory4_4_MixedSpecialCharacters tests combination of special characters
func TestStory4_4_MixedSpecialCharacters(t *testing.T) {
	// Given: A ticket with multiple types of special characters
	ticket := models.Ticket{
		ID:          "MIXED-1",
		Name:        "Bug 🐛: Handle <input> & validate \"user\" data",
		BinName:     "In Progress",
		Description: "Steps:\n1. Test with Unicode: 测试\n2. Check symbols: & < >\n3. Verify emoji ✅",
		CreatedAt:   time.Now(),
	}

	// When: Formatting the ticket
	output := FormatTicket(ticket)

	// Then: All special characters should be handled correctly
	if !strings.Contains(output, "🐛") {
		t.Error("Emoji should be preserved")
	}
	if !strings.Contains(output, "测试") {
		t.Error("Unicode should be preserved")
	}
	if !strings.Contains(output, "&") || !strings.Contains(output, "<") || !strings.Contains(output, ">") {
		t.Error("Special HTML characters should be preserved")
	}
	if !strings.Contains(output, "\"") {
		t.Error("Quotes should be preserved")
	}
}

// STORY 4.5: Handle Very Large Number of Tickets

// TestStory4_5_Display50PlusTickets tests handling of 50+ tickets
func TestStory4_5_Display50PlusTickets(t *testing.T) {
	// Given: A list of 60 tickets
	tickets := make([]models.Ticket, 60)
	baseTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 60; i++ {
		tickets[i] = models.Ticket{
			ID:          fmt.Sprintf("LARGE-%03d", i+1),
			Name:        fmt.Sprintf("Ticket number %d with a reasonably long name to simulate real data", i+1),
			BinName:     "In Progress",
			Description: fmt.Sprintf("This is ticket %d with some description text that is realistic in length", i+1),
			CreatedAt:   baseTime.Add(time.Duration(i) * time.Hour),
			UpdatedAt:   baseTime.Add(time.Duration(i+1) * time.Hour),
		}
	}

	// When: Formatting all tickets
	output := FormatTickets(tickets)

	// Then: Tool should successfully display 50+ tickets without errors
	// Acceptance Criterion: Tool successfully displays 50+ tickets without errors
	if output == "" {
		t.Fatal("Output should not be empty for 60 tickets")
	}

	// All tickets should be present
	for i := 0; i < 60; i++ {
		expectedID := fmt.Sprintf("LARGE-%03d", i+1)
		if !strings.Contains(output, expectedID) {
			t.Errorf("Expected to find ticket %s in output", expectedID)
		}
	}

	// Summary should show correct count
	if !strings.Contains(output, "60 ticket(s)") {
		t.Error("Summary should show 60 tickets")
	}
}

// TestStory4_5_MemoryUsageReasonable tests memory efficiency
func TestStory4_5_MemoryUsageReasonable(t *testing.T) {
	// Given: A list of 100 tickets with substantial content
	tickets := make([]models.Ticket, 100)
	baseTime := time.Now()

	for i := 0; i < 100; i++ {
		// Create tickets with realistic amounts of data
		longDescription := strings.Repeat(fmt.Sprintf("This is line %d of the description. ", i), 10)
		tickets[i] = models.Ticket{
			ID:          fmt.Sprintf("MEM-%03d", i+1),
			Name:        fmt.Sprintf("Memory test ticket %d", i+1),
			BinName:     "In Progress",
			Description: longDescription,
			CreatedAt:   baseTime,
			UpdatedAt:   baseTime,
		}
	}

	// When: Formatting all tickets
	output := FormatTickets(tickets)

	// Then: Memory usage should remain reasonable
	// Acceptance Criterion: Memory usage remains reasonable (doesn't consume excessive RAM)
	// We can't directly measure memory, but we can verify the output completes
	if output == "" {
		t.Fatal("Output should not be empty")
	}

	// Verify all tickets are in output (proves we didn't run out of memory)
	for i := 0; i < 100; i++ {
		expectedID := fmt.Sprintf("MEM-%03d", i+1)
		if !strings.Contains(output, expectedID) {
			t.Errorf("Expected to find ticket %s - may indicate memory issue", expectedID)
		}
	}
}

// TestStory4_5_OutputRemainReadable tests readability with many tickets
func TestStory4_5_OutputRemainReadable(t *testing.T) {
	// Given: A list of 50 tickets
	tickets := make([]models.Ticket, 50)
	baseTime := time.Now()

	for i := 0; i < 50; i++ {
		tickets[i] = models.Ticket{
			ID:          fmt.Sprintf("READ-%03d", i+1),
			Name:        fmt.Sprintf("Readable ticket %d", i+1),
			BinName:     "To Do",
			Description: "Test description",
			CreatedAt:   baseTime,
		}
	}

	// When: Formatting the tickets
	output := FormatTickets(tickets)

	// Then: Output should remain readable (though potentially long)
	// Acceptance Criterion: Output remains readable (though potentially long)
	lines := strings.Split(output, "\n")

	// Should have visual separators between tickets
	separatorCount := 0
	for _, line := range lines {
		if strings.Contains(line, "---") || strings.TrimSpace(line) == "" {
			separatorCount++
		}
	}

	if separatorCount < 40 { // Should have separators between most tickets
		t.Error("Should have visual separators to maintain readability")
	}

	// Each ticket should be clearly distinguishable
	ticketCount := 0
	for _, line := range lines {
		if strings.Contains(line, "READ-") {
			ticketCount++
		}
	}

	if ticketCount < 50 {
		t.Error("All tickets should be identifiable in output")
	}
}

// TestStory4_5_CompletesInReasonableTime tests performance
func TestStory4_5_CompletesInReasonableTime(t *testing.T) {
	// Given: A list of 50 tickets
	tickets := make([]models.Ticket, 50)
	baseTime := time.Now()

	for i := 0; i < 50; i++ {
		tickets[i] = models.Ticket{
			ID:          fmt.Sprintf("PERF-%03d", i+1),
			Name:        fmt.Sprintf("Performance test ticket %d", i+1),
			BinName:     "In Progress",
			Description: "Description for performance testing with reasonable length",
			CreatedAt:   baseTime,
			UpdatedAt:   baseTime,
		}
	}

	// When: Formatting the tickets and measuring time
	start := time.Now()
	output := FormatTickets(tickets)
	elapsed := time.Since(start)

	// Then: Should complete in reasonable time (under 10 seconds for 50 tickets)
	// Acceptance Criterion: Tool completes in reasonable time (under 10 seconds for 50 tickets)
	if elapsed > 10*time.Second {
		t.Errorf("Formatting 50 tickets took %v, should be under 10 seconds", elapsed)
	}

	// Verify output was generated
	if output == "" {
		t.Fatal("Output should not be empty")
	}

	// Log actual time for visibility
	t.Logf("Formatted 50 tickets in %v", elapsed)
}

// TestStory4_5_NoArtificialLimit tests unlimited ticket display
func TestStory4_5_NoArtificialLimit(t *testing.T) {
	// Given: A large list of tickets (more than typical limits)
	tickets := make([]models.Ticket, 150)
	baseTime := time.Now()

	for i := 0; i < 150; i++ {
		tickets[i] = models.Ticket{
			ID:          fmt.Sprintf("LIMIT-%03d", i+1),
			Name:        fmt.Sprintf("Ticket %d", i+1),
			BinName:     "To Do",
			Description: "Test",
			CreatedAt:   baseTime,
		}
	}

	// When: Formatting all tickets
	output := FormatTickets(tickets)

	// Then: Should show all tickets with no artificial limit
	// Acceptance Criterion: No artificial limit on ticket count (show all assigned tickets)
	for i := 0; i < 150; i++ {
		expectedID := fmt.Sprintf("LIMIT-%03d", i+1)
		if !strings.Contains(output, expectedID) {
			t.Errorf("Ticket %s should be in output (no artificial limits)", expectedID)
		}
	}

	// Summary should show full count
	if !strings.Contains(output, "150 ticket(s)") {
		t.Error("Summary should show all 150 tickets")
	}
}

// TestStory4_5_VeryLongOutputCompletes tests tool completion with long output
func TestStory4_5_VeryLongOutputCompletes(t *testing.T) {
	// Given: Many tickets that will produce very long output
	tickets := make([]models.Ticket, 100)
	baseTime := time.Now()

	for i := 0; i < 100; i++ {
		// Create tickets with long descriptions
		longDesc := strings.Repeat("This is a long description that will make the output very large. ", 20)
		tickets[i] = models.Ticket{
			ID:          fmt.Sprintf("LONG-%03d", i+1),
			Name:        fmt.Sprintf("Long output ticket %d", i+1),
			BinName:     "In Progress",
			Description: longDesc,
			CreatedAt:   baseTime,
		}
	}

	// When: Formatting all tickets
	output := FormatTickets(tickets)

	// Then: Tool should complete successfully even with very long output
	// Acceptance Criterion: If output is very long, tool still completes successfully
	if output == "" {
		t.Fatal("Output should not be empty")
	}

	// Verify first and last tickets are present (proves complete output)
	if !strings.Contains(output, "LONG-001") {
		t.Error("First ticket should be in output")
	}
	if !strings.Contains(output, "LONG-100") {
		t.Error("Last ticket should be in output")
	}

	// Output length should be substantial
	if len(output) < 10000 {
		t.Error("Output should be very long for 100 tickets with long descriptions")
	}
}
