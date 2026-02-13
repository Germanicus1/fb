package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Germanicus1/fb/internal/state"
)

// TestViewCheckoutStatus tests Story 1.4: View Currently Checked-Out Ticket
// Given I have a ticket checked out
// When I run the checkout status command
// Then I see the ticket name and ID
// And I see the bin name where the ticket is located
// And I see how long ago the checkout occurred
func TestViewCheckoutStatus(t *testing.T) {
	// Setup temporary directories with checkout state
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".fb", "config.yaml")
	checkoutPath := filepath.Join(tempDir, ".fb", "checkout.json")

	if err := os.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	// Setup config
	configContent := `auth_key: test-key
org_id: test-org
user_email: test@example.com
`
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Setup checkout state - checked out 2 hours ago
	twoHoursAgo := time.Now().Add(-2 * time.Hour).Unix()
	checkoutState := state.CheckoutState{
		TicketID:     "TICKET-001",
		TicketName:   "Fix login bug",
		BinID:        "bin-doing",
		BinName:      "Doing",
		CheckedOutAt: fmt.Sprintf("%d", twoHoursAgo),
	}
	checkoutData, _ := json.Marshal(checkoutState)
	if err := os.WriteFile(checkoutPath, checkoutData, 0600); err != nil {
		t.Fatalf("Failed to write checkout state: %v", err)
	}

	// Override home directory
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	var output bytes.Buffer

	// Run checkout status command
	err := runCheckoutStatus(&output)

	if err != nil {
		t.Fatalf("Expected status command to succeed, got error: %v", err)
	}

	outputStr := output.String()

	// Verify ticket name is displayed
	if !strings.Contains(outputStr, "Fix login bug") {
		t.Errorf("Expected output to show ticket name 'Fix login bug', got: %s", outputStr)
	}

	// Verify ticket ID is displayed
	if !strings.Contains(outputStr, "TICKET-001") {
		t.Errorf("Expected output to show ticket ID 'TICKET-001', got: %s", outputStr)
	}

	// Verify bin name is displayed
	if !strings.Contains(outputStr, "Doing") {
		t.Errorf("Expected output to show bin name 'Doing', got: %s", outputStr)
	}

	// Verify time since checkout is shown (human-readable)
	if !strings.Contains(outputStr, "ago") && !strings.Contains(outputStr, "hour") {
		t.Errorf("Expected output to show human-readable time, got: %s", outputStr)
	}
}

// TestViewCheckoutStatusFormatting tests the display formatting
// Given I have a ticket checked out
// When I view the status
// Then the information is formatted clearly and readable
func TestViewCheckoutStatusFormatting(t *testing.T) {
	tempDir := t.TempDir()
	checkoutPath := filepath.Join(tempDir, ".fb", "checkout.json")

	if err := os.MkdirAll(filepath.Dir(checkoutPath), 0700); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Setup checkout state
	checkoutState := state.CheckoutState{
		TicketID:     "TICKET-123",
		TicketName:   "Implement feature X",
		BinID:        "bin-review",
		BinName:      "In Review",
		CheckedOutAt: fmt.Sprintf("%d", time.Now().Add(-30*time.Minute).Unix()),
	}
	checkoutData, _ := json.Marshal(checkoutState)
	if err := os.WriteFile(checkoutPath, checkoutData, 0600); err != nil {
		t.Fatalf("Failed to write checkout state: %v", err)
	}

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	var output bytes.Buffer
	err := runCheckoutStatus(&output)

	if err != nil {
		t.Fatalf("Expected status to succeed, got: %v", err)
	}

	outputStr := output.String()

	// Verify output has structure (multiple lines with key information)
	lines := strings.Split(strings.TrimSpace(outputStr), "\n")
	if len(lines) < 2 {
		t.Errorf("Expected multi-line formatted output, got: %s", outputStr)
	}

	// Verify key-value pairs or labeled information
	hasLabels := strings.Contains(outputStr, ":") || strings.Contains(outputStr, "Ticket") || strings.Contains(outputStr, "Bin")
	if !hasLabels {
		t.Errorf("Expected labeled output format, got: %s", outputStr)
	}
}

// TestCheckoutStatusCompletesQuickly tests the performance requirement
// Given I have a ticket checked out
// When I run the status command
// Then it completes in under 500ms
func TestCheckoutStatusCompletesQuickly(t *testing.T) {
	tempDir := t.TempDir()
	checkoutPath := filepath.Join(tempDir, ".fb", "checkout.json")

	if err := os.MkdirAll(filepath.Dir(checkoutPath), 0700); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	checkoutState := state.CheckoutState{
		TicketID:     "TICKET-001",
		TicketName:   "Test",
		BinID:        "bin-1",
		BinName:      "Bin 1",
		CheckedOutAt: fmt.Sprintf("%d", time.Now().Unix()),
	}
	checkoutData, _ := json.Marshal(checkoutState)
	if err := os.WriteFile(checkoutPath, checkoutData, 0600); err != nil {
		t.Fatalf("Failed to write checkout state: %v", err)
	}

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	start := time.Now()
	var output bytes.Buffer
	err := runCheckoutStatus(&output)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Status command failed: %v", err)
	}

	if duration > 500*time.Millisecond {
		t.Errorf("Status command took %v, expected under 500ms", duration)
	}
}

// Helper function implementation for Green phase
func runCheckoutStatus(output *bytes.Buffer) error {
	// Load checkout state
	state, err := state.LoadCheckout()
	if err != nil {
		if os.IsNotExist(err) {
			// No checkout exists - display friendly message (not an error)
			fmt.Fprintf(output, "No ticket currently checked out.\n")
			fmt.Fprintf(output, "Use checkout command to select a ticket.\n")
			return nil
		}
		return err
	}

	// Display checkout information
	fmt.Fprintf(output, "Currently checked out:\n")
	fmt.Fprintf(output, "Ticket: %s - %s\n", state.TicketID, state.TicketName)
	fmt.Fprintf(output, "Bin: %s\n", state.BinName)

	// Calculate and display time since checkout
	timeSince := formatTimeSince(state.CheckedOutAt)
	fmt.Fprintf(output, "Checked out: %s\n", timeSince)

	return nil
}

// formatTimeSince converts a Unix timestamp to human-readable "X time ago" format
func formatTimeSince(timestampStr string) string {
	timestamp, err := parseTimestamp(timestampStr)
	if err != nil {
		return "unknown time"
	}

	duration := time.Since(time.Unix(timestamp, 0))
	return humanizeDuration(duration)
}

// parseTimestamp parses a timestamp string to int64
func parseTimestamp(timestampStr string) (int64, error) {
	var timestamp int64
	_, err := fmt.Sscanf(timestampStr, "%d", &timestamp)
	return timestamp, err
}

// humanizeDuration converts a duration to human-readable format
func humanizeDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 0 {
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}

	if minutes > 0 {
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	}

	return "just now"
}
