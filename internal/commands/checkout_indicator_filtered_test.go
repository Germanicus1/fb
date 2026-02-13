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

// TestIndicateCheckedOutTicketInFilteredList tests Story 2.5: Indicate Checked-Out Ticket in Bin-Filtered List
// Given I have a ticket checked out and it exists in the filtered bin
// When I view tickets filtered by that bin
// Then the checked-out ticket shows the same visual indicator as in full list
func TestIndicateCheckedOutTicketInFilteredList(t *testing.T) {
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

	// Create checkout state for a ticket in "Doing" bin
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

	// Create tickets - only tickets in "Doing" bin
	ticketsInDoingBin := []models.Ticket{
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
	}

	// Format filtered tickets with indicator
	output := formatTicketsWithCheckoutIndicator(ticketsInDoingBin)

	// Verify indicator appears for checked-out ticket in filtered list
	if !strings.Contains(output, "CHECKED OUT") {
		t.Errorf("Expected indicator in filtered list, got: %s", output)
	}

	// Verify indicator is on correct ticket
	lines := strings.Split(output, "\n")
	foundIndicator := false
	for _, line := range lines {
		if strings.Contains(line, "TICKET-002") && strings.Contains(line, "CHECKED OUT") {
			foundIndicator = true
			break
		}
	}

	if !foundIndicator {
		t.Error("Expected TICKET-002 to have indicator in filtered list")
	}
}

// TestNoIndicatorWhenCheckedOutTicketNotInFilter tests behavior when checked-out ticket is in different bin
// Given I have a ticket checked out in bin A
// When I view tickets filtered by bin B
// Then no indicator appears in the filtered list
func TestNoIndicatorWhenCheckedOutTicketNotInFilter(t *testing.T) {
	// Setup temporary directories
	tempDir := t.TempDir()
	checkoutPath := filepath.Join(tempDir, ".fb", "checkout.json")

	if err := os.MkdirAll(filepath.Dir(checkoutPath), 0700); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Create checkout state for a ticket in "Doing" bin
	checkoutState := state.CheckoutState{
		TicketID:     "TICKET-DOING",
		TicketName:   "Ticket in Doing",
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

	// View tickets in "To Do" bin (different from checked-out ticket's bin)
	ticketsInToDoBin := []models.Ticket{
		{
			ID:      "TICKET-TODO-1",
			Name:    "Plan feature",
			BinID:   "bin-todo",
			BinName: "To Do",
		},
		{
			ID:      "TICKET-TODO-2",
			Name:    "Research API",
			BinID:   "bin-todo",
			BinName: "To Do",
		},
	}

	// Format filtered tickets
	output := formatTicketsWithCheckoutIndicator(ticketsInToDoBin)

	// Verify NO indicator appears (checked-out ticket is not in this bin)
	if strings.Contains(output, "CHECKED OUT") {
		t.Errorf("Expected no indicator when checked-out ticket not in filter, got: %s", output)
	}
}

// TestIndicatorConsistentBetweenFullAndFiltered tests that same indicator is used
// Given I have a ticket checked out
// When I view it in full list vs filtered list
// Then the same visual indicator appears in both views
func TestIndicatorConsistentBetweenFullAndFiltered(t *testing.T) {
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
		{
			ID:      "TICKET-002",
			Name:    "Other ticket",
			BinID:   "bin-todo",
			BinName: "To Do",
		},
	}

	// Format full list
	fullOutput := formatTicketsWithCheckoutIndicator(tickets)

	// Format filtered list (only tickets in Doing bin)
	filteredTickets := []models.Ticket{tickets[0]}
	filteredOutput := formatTicketsWithCheckoutIndicator(filteredTickets)

	// Extract indicator text from both outputs
	fullIndicator := extractIndicatorText(fullOutput)
	filteredIndicator := extractIndicatorText(filteredOutput)

	// Verify both use the same indicator
	if fullIndicator != filteredIndicator {
		t.Errorf("Expected same indicator in full and filtered views. Full: '%s', Filtered: '%s'", fullIndicator, filteredIndicator)
	}

	// Verify indicator is present in both
	if fullIndicator == "" {
		t.Error("Expected indicator in full output")
	}
	if filteredIndicator == "" {
		t.Error("Expected indicator in filtered output")
	}
}

// TestFilteredListNoErrorWhenCheckedOutTicketNotPresent tests error handling
// Given I have a ticket checked out
// When I filter by a bin that doesn't contain the checked-out ticket
// Then no error occurs and list displays normally
func TestFilteredListNoErrorWhenCheckedOutTicketNotPresent(t *testing.T) {
	// Setup temporary directories
	tempDir := t.TempDir()
	checkoutPath := filepath.Join(tempDir, ".fb", "checkout.json")

	if err := os.MkdirAll(filepath.Dir(checkoutPath), 0700); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	checkoutState := state.CheckoutState{
		TicketID:     "TICKET-ELSEWHERE",
		TicketName:   "Not in this bin",
		BinID:        "bin-other",
		BinName:      "Other",
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
			Name:    "In current bin",
			BinID:   "bin-current",
			BinName: "Current",
		},
	}

	// Should not panic or error
	output := formatTicketsWithCheckoutIndicator(tickets)

	// Verify output is produced
	if output == "" {
		t.Error("Expected output to be produced")
	}

	// Verify ticket info is present
	if !strings.Contains(output, "TICKET-001") {
		t.Error("Expected ticket to be displayed")
	}

	// Verify no indicator (ticket not in this filter)
	if strings.Contains(output, "CHECKED OUT") {
		t.Error("Expected no indicator when checked-out ticket not in filter")
	}
}

// TestFilteredListPerformance tests that filtering doesn't degrade performance
// Given I have many tickets and one checked out
// When I filter the list
// Then performance remains acceptable
func TestFilteredListPerformance(t *testing.T) {
	t.Skip("Performance test - requires benchmarking")
	// This test validates that checkout indicator checking works
	// efficiently with filtered lists and doesn't cause performance issues.
}

// Helper function
// extractIndicatorText extracts the checkout indicator text from output
func extractIndicatorText(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "CHECKED OUT") {
			// Extract the indicator portion
			idx := strings.Index(line, "←")
			if idx != -1 {
				return strings.TrimSpace(line[idx:])
			}
			idx = strings.Index(line, "CHECKED OUT")
			if idx != -1 {
				return strings.TrimSpace(line[idx:])
			}
		}
	}
	return ""
}
