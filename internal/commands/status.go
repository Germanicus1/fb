package commands

import (
	"fmt"
	"time"

	"github.com/Germanicus1/fb/internal/state"
)

// ExecuteStatus displays the currently checked-out ticket
func ExecuteStatus() error {
	checkout, err := state.LoadCheckout()
	if err != nil {
		fmt.Println("No ticket currently checked out")
		fmt.Println("Use 'fb checkout --bin \"Bin Name\"' to check out a ticket")
		return nil
	}

	fmt.Println("Currently checked out:")
	fmt.Printf("  Ticket: [%s] %s\n", checkout.TicketID, checkout.TicketName)
	if checkout.BinName != "" {
		fmt.Printf("  Bin: %s\n", checkout.BinName)
	}

	// Show time since checkout
	checkedOutTime, err := time.Parse(time.RFC3339, checkout.CheckedOutAt)
	if err == nil {
		duration := time.Since(checkedOutTime)
		fmt.Printf("  Checked out: %s ago\n", formatDuration(duration))
	}

	return nil
}

// formatDuration formats a duration into a human-readable string
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return "less than a minute"
	}
	if d < time.Hour {
		mins := int(d.Minutes())
		if mins == 1 {
			return "1 minute"
		}
		return fmt.Sprintf("%d minutes", mins)
	}
	if d < 24*time.Hour {
		hours := int(d.Hours())
		if hours == 1 {
			return "1 hour"
		}
		return fmt.Sprintf("%d hours", hours)
	}
	days := int(d.Hours() / 24)
	if days == 1 {
		return "1 day"
	}
	return fmt.Sprintf("%d days", days)
}
