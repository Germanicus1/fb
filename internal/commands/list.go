// Package commands implements CLI command handlers for Flow Boards operations.
// Each command is responsible for user interaction, orchestrating service calls,
// and formatting output.
package commands

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Germanicus1/fb/config"
	"github.com/Germanicus1/fb/formatter"
	"github.com/Germanicus1/fb/internal/service"
	"github.com/Germanicus1/fb/internal/state"
	"github.com/Germanicus1/fb/models"
)

// Execute runs the main list command to display tickets
func Execute(cfg *config.Config, binFilter string, verbose bool) error {
	apiStart := time.Now()

	ticketService, err := service.NewTicketService(cfg)
	if err != nil {
		return err
	}

	user, err := ticketService.GetCurrentUser(cfg.UserEmail)
	if err != nil {
		return err
	}

	// Convert bin filter name to ID if needed
	binID := ""
	if binFilter != "" {
		binID, err = service.ResolveBinFilter(ticketService.GetClient(), binFilter)
		if err != nil {
			return err
		}
	}

	tickets, err := ticketService.GetUserTicketsFiltered(user.ID, binID, "")
	if err != nil {
		return err
	}

	apiDuration := time.Since(apiStart)

	displayTickets(tickets)

	if verbose {
		fmt.Fprintf(os.Stderr, "API request time: %.3fs\n", apiDuration.Seconds())
	}

	return nil
}

// displayTickets formats and displays tickets to stdout
func displayTickets(tickets []models.Ticket) {
	output := formatTicketsWithCheckoutIndicator(tickets)
	fmt.Print(output)
}

// formatTicketsWithCheckoutIndicator formats tickets and adds indicator for checked-out ticket
func formatTicketsWithCheckoutIndicator(tickets []models.Ticket) string {
	// Load current checkout state
	checkoutState, err := state.LoadCheckout()
	if err != nil {
		// No checkout or error loading - just format normally
		return formatter.FormatTickets(tickets)
	}

	// Format tickets normally
	output := formatter.FormatTickets(tickets)

	// Add indicator to checked-out ticket
	if checkoutState != nil {
		// Find lines containing the checked-out ticket ID
		lines := strings.Split(output, "\n")
		for i, line := range lines {
			if strings.Contains(line, checkoutState.TicketID) {
				// Add indicator to this line
				lines[i] = line + " ← CHECKED OUT"
			}
		}
		output = strings.Join(lines, "\n")
	}

	return output
}
