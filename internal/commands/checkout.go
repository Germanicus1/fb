package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Germanicus1/fb/config"
	"github.com/Germanicus1/fb/internal/service"
	"github.com/Germanicus1/fb/internal/state"
	"github.com/Germanicus1/fb/models"
)

// ExecuteCheckout handles the checkout command with optional bin filter and ticket ID
func ExecuteCheckout(args []string, binFlag string, forceFlag bool) error {
	if len(args) > 0 {
		// Direct checkout by ticket ID
		return ExecuteDirectCheckout(args[0])
	}

	// Checkout with bin filter or use last bin context
	if binFlag != "" {
		return ExecuteBinCheckout(binFlag, forceFlag)
	}

	// No arguments - use last bin context
	return ExecuteCheckoutWithLastBin()
}

// ExecuteBinCheckout checks out a ticket from a specific bin
func ExecuteBinCheckout(binName string, force bool) error {
	// Check for existing checkout
	if !force {
		if existing, err := state.LoadCheckout(); err == nil {
			return fmt.Errorf("ticket already checked out: %s\nUse 'fb clear' or 'fb checkout --force'", existing.TicketName)
		}
	}

	// Load config and initialize API
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	ticketService, err := service.NewTicketService(cfg)
	if err != nil {
		return err
	}

	// Get user
	user, err := ticketService.GetCurrentUser(cfg.UserEmail)
	if err != nil {
		return err
	}

	// Resolve bin name to ID
	binID, err := service.ResolveBinFilter(ticketService.GetClient(), binName)
	if err != nil {
		return err
	}

	// Fetch tickets in this bin
	tickets, err := ticketService.GetUserTicketsFiltered(user.ID, binID, "")
	if err != nil {
		return err
	}

	if len(tickets) == 0 {
		fmt.Printf("No tickets found in bin '%s'\n", binName)
		return nil
	}

	// Display tickets
	fmt.Printf("Tickets in '%s' bin:\n\n", binName)
	for i, ticket := range tickets {
		fmt.Printf("%d. [%s] %s\n", i+1, ticket.ID, ticket.Name)
	}

	// Get selection
	fmt.Print("\nEnter ticket number to checkout: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return fmt.Errorf("cancelled")
	}

	var selection int
	if _, err := fmt.Sscanf(input, "%d", &selection); err != nil || selection < 1 || selection > len(tickets) {
		return fmt.Errorf("invalid selection")
	}

	selectedTicket := tickets[selection-1]

	// Save checkout state
	checkout := state.CheckoutState{
		TicketID:     selectedTicket.ID,
		TicketName:   selectedTicket.Name,
		BinID:        selectedTicket.BinID,
		BinName:      selectedTicket.BinName,
		CheckedOutAt: time.Now().Format(time.RFC3339),
	}

	if err := state.SaveCheckout(&checkout); err != nil {
		return err
	}

	// Save bin context
	if err := state.SaveBinContext(binID, binName); err != nil {
		return err
	}

	fmt.Printf("\n✓ Checked out: %s\n", selectedTicket.Name)
	return nil
}

// ExecuteDirectCheckout checks out a ticket by ID
func ExecuteDirectCheckout(ticketID string) error {
	// Check for existing checkout
	if existing, err := state.LoadCheckout(); err == nil {
		return fmt.Errorf("ticket already checked out: %s\nUse 'fb clear' first", existing.TicketName)
	}

	// Load config and initialize API
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	ticketService, err := service.NewTicketService(cfg)
	if err != nil {
		return err
	}

	// Get user to verify ticket is assigned
	user, err := ticketService.GetCurrentUser(cfg.UserEmail)
	if err != nil {
		return err
	}

	// Fetch all user tickets and find the one with matching ID
	tickets, err := ticketService.GetUserTickets(user.ID)
	if err != nil {
		return err
	}

	var selectedTicket *models.Ticket
	for _, ticket := range tickets {
		if ticket.ID == ticketID {
			selectedTicket = &ticket
			break
		}
	}

	if selectedTicket == nil {
		return fmt.Errorf("ticket %s not found or not assigned to you", ticketID)
	}

	// Save checkout state
	checkout := state.CheckoutState{
		TicketID:     selectedTicket.ID,
		TicketName:   selectedTicket.Name,
		BinID:        selectedTicket.BinID,
		BinName:      selectedTicket.BinName,
		CheckedOutAt: time.Now().Format(time.RFC3339),
	}

	if err := state.SaveCheckout(&checkout); err != nil {
		return err
	}

	fmt.Printf("✓ Checked out: %s\n", selectedTicket.Name)
	return nil
}

// ExecuteCheckoutWithLastBin checks out using the last used bin context
func ExecuteCheckoutWithLastBin() error {
	binContext, err := state.LoadBinContext()
	if err != nil {
		return fmt.Errorf("no bin context found. Use 'fb checkout --bin \"Bin Name\"' first")
	}

	return ExecuteBinCheckout(binContext.BinName, false)
}

// ExecuteClear clears the current checkout state
func ExecuteClear() error {
	if err := state.ClearCheckout(); err != nil {
		return err
	}
	fmt.Println("✓ Checkout cleared")
	return nil
}
