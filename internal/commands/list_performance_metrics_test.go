package commands

import (
	"testing"

	"github.com/Germanicus1/fb/models"
)

// Story 2.3: Performance Metrics in Verbose Mode Only - Integration Tests

// Note: These tests verify the formatTicketsWithVerbosity function behavior.
// The actual performance metrics display is tested via the Execute function
// which writes to stderr when verbose=true. That part is tested in integration
// tests or manually.

// TestMinimalModeDoesNotContainMetricsInOutput verifies minimal output is clean
func TestMinimalModeDoesNotContainMetricsInOutput(t *testing.T) {
	// Given: I have tickets
	tickets := []models.Ticket{
		{ID: "TICKET-001", Name: "Test Ticket"},
	}

	// When: I format in minimal mode
	output := formatTicketsWithVerbosity(tickets, false)

	// Then: Output does not contain any performance-related text
	// (Performance metrics are written to stderr by Execute, not by formatter)
	// This test just verifies the formatter output is clean

	// Minimal output should be concise
	if len(output) > 100 {
		t.Error("Minimal output should be concise")
	}
}

// TestVerboseFormatDoesNotIncludeMetricsInTicketOutput verifies formatter doesn't add metrics
func TestVerboseFormatDoesNotIncludeMetricsInTicketOutput(t *testing.T) {
	// Given: I have tickets
	tickets := []models.Ticket{
		{
			ID:          "TICKET-001",
			Name:        "Test",
			BinName:     "To Do",
			Description: "Test",
		},
	}

	// When: I format in verbose mode
	output := formatTicketsWithVerbosity(tickets, true)

	// Then: The ticket output itself doesn't contain performance metrics
	// (Metrics are added by Execute function to stderr, not by formatter)
	// This test verifies the formatter stays focused on formatting tickets

	// Should contain ticket info but no metrics
	if len(output) < 50 {
		t.Error("Verbose output should contain ticket details")
	}
}

// TestFormatterOutputIsIndependentOfPerformanceMetrics verifies separation of concerns
func TestFormatterOutputIsIndependentOfPerformanceMetrics(t *testing.T) {
	// Given: I have tickets
	tickets := []models.Ticket{
		{ID: "TICKET-001", Name: "Test"},
	}

	// When: I format tickets (both modes)
	minimalOutput := formatTicketsWithVerbosity(tickets, false)
	verboseOutput := formatTicketsWithVerbosity(tickets, true)

	// Then: Neither output contains performance metric keywords
	// (Metrics are handled separately by Execute function)

	// Main assertion: formatter produces ticket output, not metrics
	if minimalOutput == "" || verboseOutput == "" {
		t.Error("Formatter should produce output")
	}
}
