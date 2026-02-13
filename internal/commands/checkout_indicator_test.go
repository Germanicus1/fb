package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Germanicus1/fb/internal/state"
	"github.com/Germanicus1/fb/models"
)

// TestIndicateCheckedOutTicketInMainList tests Story 2.4: Indicate Checked-Out Ticket in List Display
// Given I have a ticket checked out
// When I view all tickets
// Then the checked-out ticket is marked with a visual indicator
// And the indicator is easily distinguished from other tickets
// And all other ticket information remains the same
func TestIndicateCheckedOutTicketInMainList(t *testing.T) {
	// Setup temporary directories
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".fb", "config.yaml")
	checkoutPath := filepath.Join(tempDir, ".fb", "checkout.json")

	if err := os.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	configContent := `auth_key: test-key
org_id: test-org
user_email: test@example.com
`
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Create checkout state for TICKET-002
	checkoutState := state.CheckoutState{
		TicketID:     "TICKET-002",
		TicketName:   "Add user profile",
		BinID:        "bin-doing",
		BinName:      "Doing",
		CheckedOutAt: "1234567890",
	}
	checkoutData, _ := json.Marshal(checkoutState)
	if err := os.WriteFile(checkoutPath, checkoutData, 0600); err != nil {
		t.Fatalf("Failed to write checkout state: %v", err)
	}

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create test tickets
	tickets := []models.Ticket{
		{
			ID:      "TICKET-001",
			Name:    "Fix login bug",
			BinID:   "bin-doing",
			BinName: "Doing",
		},
		{
			ID:      "TICKET-002",
			Name:    "Add user profile",
			BinID:   "bin-doing",
			BinName: "Doing",
		},
		{
			ID:      "TICKET-003",
			Name:    "Update documentation",
			BinID:   "bin-todo",
			BinName: "To Do",
		},
	}

	// Format tickets with checkout indicator
	output := formatTicketsWithCheckoutIndicator(tickets)

	// Verify visual indicator appears for checked-out ticket
	if !strings.Contains(output, "← CHECKED OUT") && !strings.Contains(output, "CHECKED OUT") {
		t.Errorf("Expected visual indicator for checked-out ticket, got: %s", output)
	}

	// Verify indicator appears on correct ticket (TICKET-002)
	lines := strings.Split(output, "\n")
	var ticket002Line string
	for _, line := range lines {
		if strings.Contains(line, "TICKET-002") || strings.Contains(line, "Add user profile") {
			ticket002Line = line
			break
		}
	}

	if ticket002Line == "" {
		t.Fatal("Could not find TICKET-002 in output")
	}

	if !strings.Contains(ticket002Line, "CHECKED OUT") {
		t.Errorf("Expected TICKET-002 line to have indicator, got: %s", ticket002Line)
	}

	// Verify other tickets do NOT have indicator
	for _, line := range lines {
		if (strings.Contains(line, "TICKET-001") || strings.Contains(line, "Fix login bug")) &&
		   strings.Contains(line, "CHECKED OUT") {
			t.Errorf("TICKET-001 should not have indicator, got: %s", line)
		}
		if (strings.Contains(line, "TICKET-003") || strings.Contains(line, "Update documentation")) &&
		   strings.Contains(line, "CHECKED OUT") {
			t.Errorf("TICKET-003 should not have indicator, got: %s", line)
		}
	}
}

// TestIndicatorIsVisuallySeparated tests that indicator is easily distinguished
// Given I view tickets with a checkout indicator
// When I read the output
// Then the indicator is clearly separated from ticket information
func TestIndicatorIsVisuallySeparated(t *testing.T) {
	// Setup temporary directories
	tempDir := t.TempDir()
	checkoutPath := filepath.Join(tempDir, ".fb", "checkout.json")

	if err := os.MkdirAll(filepath.Dir(checkoutPath), 0700); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	checkoutState := state.CheckoutState{
		TicketID:     "TICKET-001",
		TicketName:   "Test ticket",
		BinID:        "bin-doing",
		BinName:      "Doing",
		CheckedOutAt: "1234567890",
	}
	checkoutData, _ := json.Marshal(checkoutState)
	if err := os.WriteFile(checkoutPath, checkoutData, 0600); err != nil {
		t.Fatalf("Failed to write checkout state: %v", err)
	}

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	tickets := []models.Ticket{
		{
			ID:      "TICKET-001",
			Name:    "Test ticket",
			BinID:   "bin-doing",
			BinName: "Doing",
		},
	}

	output := formatTicketsWithCheckoutIndicator(tickets)

	// Verify indicator is right-aligned or clearly separated
	// The indicator should not interfere with reading ticket information
	if !strings.Contains(output, "CHECKED OUT") {
		t.Error("Expected indicator in output")
	}

	// Verify all ticket info is still present
	if !strings.Contains(output, "TICKET-001") {
		t.Error("Expected ticket ID to be present")
	}
	if !strings.Contains(output, "Test ticket") {
		t.Error("Expected ticket name to be present")
	}
}

// TestNoIndicatorWhenNoCheckout tests that no indicator appears without checkout
// Given I have no ticket checked out
// When I view tickets
// Then no checkout indicator appears
func TestNoIndicatorWhenNoCheckout(t *testing.T) {
	// Setup temporary directories - no checkout file created
	tempDir := t.TempDir()

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	tickets := []models.Ticket{
		{
			ID:      "TICKET-001",
			Name:    "Test ticket",
			BinID:   "bin-doing",
			BinName: "Doing",
		},
		{
			ID:      "TICKET-002",
			Name:    "Another ticket",
			BinID:   "bin-doing",
			BinName: "Doing",
		},
	}

	output := formatTicketsWithCheckoutIndicator(tickets)

	// Verify no indicator appears
	if strings.Contains(output, "CHECKED OUT") {
		t.Errorf("Expected no indicator when no checkout exists, got: %s", output)
	}
}

// TestIndicatorDoesNotAffectOtherTickets tests that only checked-out ticket is marked
// Given I have multiple tickets and one checked out
// When I view tickets
// Then only the checked-out ticket has the indicator
func TestIndicatorDoesNotAffectOtherTickets(t *testing.T) {
	// Setup temporary directories
	tempDir := t.TempDir()
	checkoutPath := filepath.Join(tempDir, ".fb", "checkout.json")

	if err := os.MkdirAll(filepath.Dir(checkoutPath), 0700); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	checkoutState := state.CheckoutState{
		TicketID:     "TICKET-002",
		TicketName:   "Middle ticket",
		BinID:        "bin-doing",
		BinName:      "Doing",
		CheckedOutAt: "1234567890",
	}
	checkoutData, _ := json.Marshal(checkoutState)
	if err := os.WriteFile(checkoutPath, checkoutData, 0600); err != nil {
		t.Fatalf("Failed to write checkout state: %v", err)
	}

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	tickets := []models.Ticket{
		{
			ID:      "TICKET-001",
			Name:    "First ticket",
			BinID:   "bin-doing",
			BinName: "Doing",
		},
		{
			ID:      "TICKET-002",
			Name:    "Middle ticket",
			BinID:   "bin-doing",
			BinName: "Doing",
		},
		{
			ID:      "TICKET-003",
			Name:    "Last ticket",
			BinID:   "bin-doing",
			BinName: "Doing",
		},
	}

	output := formatTicketsWithCheckoutIndicator(tickets)
	lines := strings.Split(output, "\n")

	indicatorCount := 0
	for _, line := range lines {
		if strings.Contains(line, "CHECKED OUT") {
			indicatorCount++
			// Verify it's on the correct ticket
			if !strings.Contains(line, "TICKET-002") && !strings.Contains(line, "Middle ticket") {
				t.Errorf("Indicator appears on wrong ticket: %s", line)
			}
		}
	}

	// Should have exactly one indicator
	if indicatorCount != 1 {
		t.Errorf("Expected exactly 1 indicator, got %d", indicatorCount)
	}
}

// TestIndicatorPerformance tests that indicator doesn't degrade performance
// Given I have many tickets
// When I format them with indicator
// Then performance remains acceptable
func TestIndicatorPerformance(t *testing.T) {
	t.Skip("Performance test - requires benchmarking")
	// This test validates that adding checkout indicator checking
	// doesn't significantly impact list rendering performance.
	// Acceptable performance: <100ms for 100 tickets
}

// Helper function implementation
// formatTicketsWithCheckoutIndicator is now implemented in main.go
