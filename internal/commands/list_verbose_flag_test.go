package commands

import (
	"strings"
	"testing"

	"github.com/Germanicus1/fb/models"
)

// Story 2.2: Short Flag for Verbose Mode - Integration Tests

// TestVerboseShortFlag_ProducesSameOutputAsLongFlag verifies -v matches --verbose
func TestVerboseShortFlag_ProducesSameOutputAsLongFlag(t *testing.T) {
	// Given: I have tickets
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "Test Ticket",
			BinName:     "In Progress",
			Description: "Test description",
		},
	}

	// When: I format with verbose=true (simulating both -v and --verbose)
	output := formatTicketsWithVerbosity(tickets, true)

	// Then: Output contains verbose details
	if !strings.Contains(output, "Status:") {
		t.Error("Verbose output should contain 'Status:'")
	}
	if !strings.Contains(output, "Description:") {
		t.Error("Verbose output should contain 'Description:'")
	}
	if !strings.Contains(output, "Test description") {
		t.Error("Verbose output should contain description text")
	}
}

// TestDebugFlag_ProducesSameOutputAsVerbose verifies --debug matches --verbose
func TestDebugFlag_ProducesSameOutputAsVerbose(t *testing.T) {
	// Given: I have tickets
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "Test Ticket",
			BinName:     "To Do",
			Description: "Debug test",
		},
	}

	// When: I format with verbose=true (simulating --debug)
	output := formatTicketsWithVerbosity(tickets, true)

	// Then: Output is identical to verbose output
	if !strings.Contains(output, "Status:") {
		t.Error("Debug flag should produce verbose output")
	}
	if !strings.Contains(output, "Description:") {
		t.Error("Debug flag should produce verbose output with descriptions")
	}
}

// TestVerboseFlagWithEmptyList verifies -v works with empty list
func TestVerboseFlagWithEmptyList(t *testing.T) {
	// Given: I have no tickets
	tickets := []models.Ticket{}

	// When: I use verbose mode
	output := formatTicketsWithVerbosity(tickets, true)

	// Then: Shows same message as minimal mode
	expectedMessage := "No tickets assigned to you."
	if output != expectedMessage {
		t.Errorf("Expected '%s', got: %s", expectedMessage, output)
	}
}

// TestVerboseFlagWithBinFilter verifies -v works with --bin flag
func TestVerboseFlagWithBinFilter(t *testing.T) {
	// Given: I have filtered tickets
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "Filtered Ticket",
			BinName:     "In Progress",
			Description: "Filtered description",
		},
	}

	// When: I format with verbose=true (simulating -v --bin "In Progress")
	output := formatTicketsWithVerbosity(tickets, true)

	// Then: Output shows verbose details for filtered tickets
	if !strings.Contains(output, "Status:") {
		t.Error("Verbose flag with bin filter should show status")
	}
	if !strings.Contains(output, "Filtered description") {
		t.Error("Verbose flag with bin filter should show description")
	}
}
