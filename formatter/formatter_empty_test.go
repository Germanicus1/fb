package formatter

import (
	"strings"
	"testing"

	"github.com/Germanicus1/fb/models"
)

// Story 1.3: Empty Ticket List Handling - Acceptance Tests

// TestStory1_3_EmptyListShowsClearMessage verifies clear message for no tickets
// Acceptance Criterion 1: When I run `fb` and have no assigned tickets, I see: "No tickets assigned to you."
func TestStory1_3_EmptyListShowsClearMessage(t *testing.T) {
	// Given: I have no assigned tickets
	tickets := []models.Ticket{}

	// When: I format the empty list in minimal mode
	output := FormatTicketsMinimal(tickets)

	// Then: I see a clear message
	if output != "No tickets assigned to you." {
		t.Errorf("Expected 'No tickets assigned to you.', got: %s", output)
	}
}

// TestStory1_3_EmptyListVerboseModeSameMessage verifies consistency across modes
// Acceptance Criterion 2: The message is identical whether using minimal or verbose mode
func TestStory1_3_EmptyListVerboseModeSameMessage(t *testing.T) {
	// Given: I have no assigned tickets
	tickets := []models.Ticket{}

	// When: I format in both minimal and verbose mode
	minimalOutput := FormatTicketsMinimal(tickets)
	verboseOutput := FormatTickets(tickets)

	// Then: Both show the same message
	if minimalOutput != verboseOutput {
		t.Errorf("Empty list message should be identical in both modes.\nMinimal: %s\nVerbose: %s", minimalOutput, verboseOutput)
	}

	// And: Both show the expected message
	expectedMessage := "No tickets assigned to you."
	if minimalOutput != expectedMessage {
		t.Errorf("Expected '%s', got: %s", expectedMessage, minimalOutput)
	}
}

// TestStory1_3_NoErrorMessageOrStackTrace verifies no error-like output
// Acceptance Criterion 3: No error message or stack trace appears
func TestStory1_3_NoErrorMessageOrStackTrace(t *testing.T) {
	// Given: I have no assigned tickets
	tickets := []models.Ticket{}

	// When: I format the empty list
	output := FormatTicketsMinimal(tickets)

	// Then: No error-like text appears
	lowerOutput := strings.ToLower(output)

	if strings.Contains(lowerOutput, "error") {
		t.Error("Output should not contain 'error'")
	}
	if strings.Contains(lowerOutput, "failed") {
		t.Error("Output should not contain 'failed'")
	}
	if strings.Contains(lowerOutput, "exception") {
		t.Error("Output should not contain 'exception'")
	}
	if strings.Contains(lowerOutput, "panic") {
		t.Error("Output should not contain 'panic'")
	}
	if strings.Contains(output, "stack trace") {
		t.Error("Output should not contain 'stack trace'")
	}
}

// TestStory1_3_ExitStatusSuccess verifies function completes successfully
// Acceptance Criterion 4: The command exits with success status (exit code 0)
// Note: This test verifies the formatter doesn't panic or error
func TestStory1_3_ExitStatusSuccess(t *testing.T) {
	// Given: I have no assigned tickets
	tickets := []models.Ticket{}

	// When: I format the empty list
	// Then: Function should complete without panic (implicit test)
	output := FormatTicketsMinimal(tickets)

	// And: Should return a string (not error)
	if output == "" {
		t.Error("Output should not be empty string")
	}
}

// TestStory1_3_NoAdditionalOutputBeyondMessage verifies clean output
// Acceptance Criterion 5: No additional output beyond the "no tickets" message appears
func TestStory1_3_NoAdditionalOutputBeyondMessage(t *testing.T) {
	// Given: I have no assigned tickets
	tickets := []models.Ticket{}

	// When: I format the empty list
	output := FormatTicketsMinimal(tickets)

	// Then: Output is exactly the message with no additional content
	expectedMessage := "No tickets assigned to you."
	if output != expectedMessage {
		t.Errorf("Output should be exactly '%s' with no additional content, got: %s", expectedMessage, output)
	}

	// And: No header line like "Found 0 ticket(s)"
	if strings.Contains(output, "Found 0") {
		t.Error("Empty list should not show 'Found 0 ticket(s)' header")
	}

	// And: No blank lines before or after
	if strings.HasPrefix(output, "\n") || strings.HasSuffix(output, "\n") {
		t.Error("Output should not have leading or trailing newlines")
	}
}

// TestStory1_3_MessageIsUserFriendly verifies friendly tone
func TestStory1_3_MessageIsUserFriendly(t *testing.T) {
	// Given: I have no assigned tickets
	tickets := []models.Ticket{}

	// When: I format the empty list
	output := FormatTicketsMinimal(tickets)

	// Then: Message should be friendly and clear
	// Should use "you" to make it personal
	if !strings.Contains(output, "you") {
		t.Error("Message should use 'you' to be personal and friendly")
	}

	// Should explicitly state "No tickets" (not just "None" or "Empty")
	if !strings.Contains(output, "No tickets") {
		t.Error("Message should explicitly state 'No tickets'")
	}
}

// TestStory1_3_EmptyAfterFiltering verifies empty result after filtering
func TestStory1_3_EmptyAfterFiltering(t *testing.T) {
	// Given: Filter results in zero tickets (simulating fb --bin "NonExistent")
	filteredTickets := []models.Ticket{}

	// When: I format the empty filtered result
	output := FormatTicketsMinimal(filteredTickets)

	// Then: Shows same clear message
	if output != "No tickets assigned to you." {
		t.Errorf("Filtered empty list should show same message, got: %s", output)
	}
}

// TestStory1_3_NilTicketListHandled verifies nil slice handling
func TestStory1_3_NilTicketListHandled(t *testing.T) {
	// Given: I have a nil ticket slice (edge case)
	var tickets []models.Ticket = nil

	// When: I format the nil list
	output := FormatTicketsMinimal(tickets)

	// Then: Shows same clear message (no panic)
	if output != "No tickets assigned to you." {
		t.Errorf("Nil ticket list should show same message, got: %s", output)
	}
}

// TestStory1_3_OutputIsConcise verifies message is short and scannable
func TestStory1_3_OutputIsConcise(t *testing.T) {
	// Given: I have no assigned tickets
	tickets := []models.Ticket{}

	// When: I format the empty list
	output := FormatTicketsMinimal(tickets)

	// Then: Message should be concise (single line, under 50 characters)
	if strings.Contains(output, "\n") {
		t.Error("Message should be single line")
	}

	if len(output) > 50 {
		t.Errorf("Message should be concise (under 50 chars), got %d chars", len(output))
	}
}

// TestStory1_3_ConsistentWithNonEmptyFormat verifies format consistency
func TestStory1_3_ConsistentWithNonEmptyFormat(t *testing.T) {
	// Given: I have empty and non-empty ticket lists
	emptyTickets := []models.Ticket{}
	nonEmptyTickets := []models.Ticket{
		{ID: "TICKET-001", Name: "Test Ticket"},
	}

	// When: I format both lists
	emptyOutput := FormatTicketsMinimal(emptyTickets)
	nonEmptyOutput := FormatTicketsMinimal(nonEmptyTickets)

	// Then: Empty output should be plain text (like non-empty)
	// No special formatting that would be inconsistent
	// Both should be plain text with no control characters
	for _, r := range emptyOutput {
		if r < 32 && r != '\n' && r != '\t' {
			t.Errorf("Empty output should not contain control characters, found: %U", r)
		}
	}

	// Both should be printable
	if emptyOutput == "" {
		t.Error("Empty output should not be empty string")
	}
	if nonEmptyOutput == "" {
		t.Error("Non-empty output should not be empty string")
	}
}

// TestStory1_3_DistinguishFromZeroFilteredResults verifies clear meaning
func TestStory1_3_DistinguishFromZeroFilteredResults(t *testing.T) {
	// Given: I have no tickets (could be no tickets assigned OR no tickets match filter)
	tickets := []models.Ticket{}

	// When: I format the empty list
	output := FormatTicketsMinimal(tickets)

	// Then: Message should indicate "assigned to you" to clarify scope
	// This distinguishes "no tickets exist" from "no tickets match filter"
	if !strings.Contains(output, "assigned to you") {
		t.Error("Message should clarify tickets 'assigned to you' to distinguish from filter results")
	}
}

// TestStory1_3_VerboseModeAlsoHandlesEmpty verifies verbose mode works too
func TestStory1_3_VerboseModeAlsoHandlesEmpty(t *testing.T) {
	// Given: I have no assigned tickets
	tickets := []models.Ticket{}

	// When: I format in verbose mode
	output := FormatTickets(tickets)

	// Then: Shows the same clear message
	if output != "No tickets assigned to you." {
		t.Errorf("Verbose mode should show same message for empty list, got: %s", output)
	}

	// And: No error-like output
	if strings.Contains(strings.ToLower(output), "error") {
		t.Error("Verbose mode should not show error for empty list")
	}
}

// TestStory1_3_EmptyListMessageIsAccessible verifies message clarity
func TestStory1_3_EmptyListMessageIsAccessible(t *testing.T) {
	// Given: I have no assigned tickets
	tickets := []models.Ticket{}

	// When: I format the empty list
	output := FormatTicketsMinimal(tickets)

	// Then: Message should be clear English with no technical jargon
	// Should not use terms like "null", "empty array", "zero records"
	lowerOutput := strings.ToLower(output)

	if strings.Contains(lowerOutput, "null") {
		t.Error("Message should avoid technical term 'null'")
	}
	if strings.Contains(lowerOutput, "array") {
		t.Error("Message should avoid technical term 'array'")
	}
	if strings.Contains(lowerOutput, "empty") && !strings.Contains(lowerOutput, "no tickets") {
		t.Error("Message should say 'No tickets' rather than just 'empty'")
	}
}
