package commands

import (
	"strings"
	"testing"

	"github.com/Germanicus1/fb/models"
)

// Integration Tests for Minimal/Verbose List Output

// TestListCommand_DefaultUsesMinimalFormat verifies default output is minimal
func TestListCommand_DefaultUsesMinimalFormat(t *testing.T) {
	// Given: I have tickets
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "Test Ticket",
			BinName:     "In Progress",
			Description: "This should not appear in minimal mode",
		},
	}

	// When: I format with verbose=false (default)
	output := formatTicketsWithCheckoutIndicator(tickets, false)

	// Then: Output uses minimal format
	if !strings.Contains(output, "[TICKET-001] Test Ticket") {
		t.Error("Should show ticket in minimal format")
	}

	// And: Does not contain verbose details
	if strings.Contains(output, "Status:") {
		t.Error("Minimal mode should not show 'Status:'")
	}
	if strings.Contains(output, "Description:") {
		t.Error("Minimal mode should not show 'Description:'")
	}
	if strings.Contains(output, "This should not appear") {
		t.Error("Minimal mode should not show description text")
	}
}

// TestListCommand_VerboseFlagShowsDetails verifies --verbose shows detailed output
func TestListCommand_VerboseFlagShowsDetails(t *testing.T) {
	// Given: I have tickets
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "Test Ticket",
			BinName:     "In Progress",
			Description: "Detailed description",
		},
	}

	// When: I format with verbose=true
	output := formatTicketsWithVerbosity(tickets, true)

	// Then: Output uses verbose format
	if !strings.Contains(output, "[TICKET-001] Test Ticket") {
		t.Error("Should show ticket header")
	}

	// And: Contains verbose details
	if !strings.Contains(output, "Status:") {
		t.Error("Verbose mode should show 'Status:'")
	}
	if !strings.Contains(output, "Description:") {
		t.Error("Verbose mode should show 'Description:'")
	}
	if !strings.Contains(output, "Detailed description") {
		t.Error("Verbose mode should show description text")
	}
}

// TestListCommand_MinimalVsVerboseOutputDiffers verifies modes are different
func TestListCommand_MinimalVsVerboseOutputDiffers(t *testing.T) {
	// Given: I have tickets with descriptions
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "Test",
			BinName:     "To Do",
			Description: "Some description",
		},
	}

	// When: I format in both modes
	minimalOutput := formatTicketsWithVerbosity(tickets, false)
	verboseOutput := formatTicketsWithVerbosity(tickets, true)

	// Then: Outputs should be different
	if minimalOutput == verboseOutput {
		t.Error("Minimal and verbose output should be different")
	}

	// And: Verbose should be longer
	minimalLines := len(strings.Split(minimalOutput, "\n"))
	verboseLines := len(strings.Split(verboseOutput, "\n"))

	if verboseLines <= minimalLines {
		t.Errorf("Verbose output (%d lines) should be longer than minimal (%d lines)", verboseLines, minimalLines)
	}
}

// TestListCommand_CheckoutIndicatorInMinimalMode verifies checkout works in minimal
func TestListCommand_CheckoutIndicatorInMinimalMode(t *testing.T) {
	// Given: I have tickets (checkout indicator is tested separately)
	tickets := []models.Ticket{
		{ID: "TICKET-001", Name: "First Ticket"},
		{ID: "TICKET-002", Name: "Second Ticket"},
	}

	// When: I format in minimal mode
	output := formatTicketsWithVerbosity(tickets, false)

	// Then: Output is in minimal format
	if !strings.Contains(output, "[TICKET-001]") {
		t.Error("Should show tickets in minimal format")
	}

	// Note: Actual checkout indicator integration is tested in checkout_indicator_test.go
	// This test just verifies the basic formatting works
}

// TestListCommand_EmptyListInMinimalMode verifies empty list handling
func TestListCommand_EmptyListInMinimalMode(t *testing.T) {
	// Given: I have no tickets
	tickets := []models.Ticket{}

	// When: I format in minimal mode
	output := formatTicketsWithVerbosity(tickets, false)

	// Then: Shows clear message
	if !strings.Contains(output, "No tickets assigned to you.") {
		t.Errorf("Should show no tickets message, got: %s", output)
	}
}

// TestListCommand_EmptyListInVerboseMode verifies empty list in verbose
func TestListCommand_EmptyListInVerboseMode(t *testing.T) {
	// Given: I have no tickets
	tickets := []models.Ticket{}

	// When: I format in verbose mode
	output := formatTicketsWithVerbosity(tickets, true)

	// Then: Shows same message as minimal
	if !strings.Contains(output, "No tickets assigned to you.") {
		t.Errorf("Should show no tickets message, got: %s", output)
	}
}

// TestListCommand_ManyTicketsMinimalMode verifies minimal works with many tickets
func TestListCommand_ManyTicketsMinimalMode(t *testing.T) {
	// Given: I have 20 tickets
	tickets := make([]models.Ticket, 20)
	for i := 0; i < 20; i++ {
		tickets[i] = models.Ticket{
			ID:          "TICKET-" + string(rune('A'+i)),
			Name:        "Ticket " + string(rune('A'+i)),
			BinName:     "To Do",
			Description: "Description",
		}
	}

	// When: I format in minimal mode
	output := formatTicketsWithVerbosity(tickets, false)

	// Then: Output is compact (approximately 22 lines)
	lines := strings.Split(output, "\n")
	lineCount := len(lines) - 1

	if lineCount < 21 || lineCount > 23 {
		t.Errorf("Expected approximately 22 lines for 20 tickets in minimal mode, got %d", lineCount)
	}

	// And: No description text appears
	if strings.Contains(output, "Description:") {
		t.Error("Minimal mode should not show description labels")
	}
}

// TestListCommand_ManyTicketsVerboseMode verifies verbose works with many tickets
func TestListCommand_ManyTicketsVerboseMode(t *testing.T) {
	// Given: I have 20 tickets
	tickets := make([]models.Ticket, 20)
	for i := 0; i < 20; i++ {
		tickets[i] = models.Ticket{
			ID:          "TICKET-" + string(rune('A'+i)),
			Name:        "Ticket " + string(rune('A'+i)),
			BinName:     "To Do",
			Description: "Description",
		}
	}

	// When: I format in verbose mode
	output := formatTicketsWithVerbosity(tickets, true)

	// Then: Output is detailed (more than 60 lines for 20 tickets with details)
	lines := strings.Split(output, "\n")
	lineCount := len(lines) - 1

	if lineCount < 60 {
		t.Errorf("Expected more than 60 lines for 20 tickets in verbose mode, got %d", lineCount)
	}

	// And: Description labels appear for all tickets
	descCount := strings.Count(output, "Description:")
	if descCount != 20 {
		t.Errorf("Expected 'Description:' to appear 20 times, got %d", descCount)
	}
}
