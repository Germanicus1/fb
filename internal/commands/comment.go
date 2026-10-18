package commands

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Germanicus1/fb/config"
	"github.com/Germanicus1/fb/internal/prompt"
	"github.com/Germanicus1/fb/internal/service"
	"github.com/Germanicus1/fb/internal/state"
	"github.com/Germanicus1/fb/models"
)

// ExecuteInteractive enters interactive comment mode to add a comment to a ticket
func ExecuteInteractive(cfg *config.Config, binFilter string) error {
	return ExecuteInteractiveWithOutput(os.Stdout, binFilter, cfg)
}

// ExecuteInteractiveWithOutput enters interactive comment mode with custom output writer (for testing)
func ExecuteInteractiveWithOutput(output io.Writer, binFilter string, cfg *config.Config) error {
	ticketService, err := service.NewTicketService(cfg)
	if err != nil {
		return err
	}

	user, err := ticketService.GetCurrentUser(cfg.UserEmail)
	if err != nil {
		return err
	}

	// Resolve bin filter if provided
	binID := ""
	if binFilter != "" {
		bin, err := service.LookupBin(ticketService.GetClient(), binFilter)
		if err != nil {
			return err
		}
		binID = bin.ID
	}

	// Fetch tickets with optional bin filter
	var tickets []models.Ticket
	if binID != "" {
		tickets, err = ticketService.GetUserTicketsFiltered(user.ID, binID, "")
	} else {
		tickets, err = ticketService.GetUserTickets(user.ID)
	}
	if err != nil {
		return err
	}

	displayTicketsForSelection(output, tickets, binFilter)

	if len(tickets) == 0 {
		return nil
	}

	in := prompt.NewLineReader(os.Stdin, output)
	selectedTicket, err := selectTicketByNumber(in, output, tickets)
	if err != nil {
		return err
	}

	comment, err := enterComment(in, output)
	if err != nil {
		return err
	}

	fmt.Fprintf(output, "Posting comment...\n")

	commentID := service.GenerateCommentID()
	payload := service.BuildCommentPayload(commentID, selectedTicket.ID, comment)

	err = service.PostComment(ticketService.GetClient(), payload)
	if err != nil {
		return err
	}

	displaySuccessConfirmation(output, selectedTicket)

	return nil
}

// ExecuteQuick adds a comment to the checked-out ticket
func ExecuteQuick(comment string) error {
	// Load checkout state
	checkout, err := state.LoadCheckout()
	if err != nil {
		return fmt.Errorf("no ticket checked out. Use 'fb checkout' first")
	}

	// Post comment via API
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	ticketService, err := service.NewTicketService(cfg)
	if err != nil {
		return err
	}

	commentID := service.GenerateCommentID()
	payload := service.BuildCommentPayload(commentID, checkout.TicketID, comment)

	if err := service.PostComment(ticketService.GetClient(), payload); err != nil {
		return err
	}

	fmt.Printf("✓ Comment added to: %s\n", checkout.TicketName)
	return nil
}

// displayTicketsForSelection displays a numbered list of tickets for selection
func displayTicketsForSelection(output io.Writer, tickets []models.Ticket, binFilter string) {
	if len(tickets) == 0 {
		if binFilter != "" {
			fmt.Fprintf(output, "No tickets found in bin '%s'. Cannot add comment.\n", binFilter)
		} else {
			fmt.Fprintf(output, "No tickets assigned to you. Cannot add comment.\n")
		}
		return
	}

	for i, ticket := range tickets {
		fmt.Fprintf(output, "%d. %s - %s [%s]\n", i+1, ticket.ID, ticket.Name, ticket.Status())
	}
}

// selectTicketByNumber prompts the user to select a ticket by number
func selectTicketByNumber(in prompt.LineReader, output io.Writer, tickets []models.Ticket) (*models.Ticket, error) {
	for {
		line, err := in.ReadLine("Enter ticket number to comment on: ")
		userInput := strings.TrimSpace(line)
		if err != nil || userInput == "" {
			fmt.Fprintf(output, "Comment cancelled.\n")
			return nil, fmt.Errorf("operation cancelled")
		}

		var ticketNum int
		_, err = fmt.Sscanf(userInput, "%d", &ticketNum)
		if err != nil || ticketNum < 1 || ticketNum > len(tickets) {
			fmt.Fprintf(output, "Invalid ticket number. Please enter a number between 1 and %d.\n", len(tickets))
			continue
		}

		selectedTicket := &tickets[ticketNum-1]
		fmt.Fprintf(output, "Selected: %s\n", selectedTicket.Name)
		return selectedTicket, nil
	}
}

// enterComment prompts the user to enter a comment
func enterComment(in prompt.LineReader, output io.Writer) (string, error) {
	for {
		line, err := in.ReadLine("Enter comment: ")
		comment := strings.TrimSpace(line)
		if comment == "" && err != nil {
			return "", fmt.Errorf("operation cancelled")
		}
		if comment == "" {
			fmt.Fprintf(output, "Comment cannot be empty. Please enter some text.\n")
			continue
		}

		return comment, nil
	}
}

// displaySuccessConfirmation displays success message after posting a comment
func displaySuccessConfirmation(output io.Writer, ticket *models.Ticket) {
	fmt.Fprintf(output, "Comment added successfully to ticket: %s\n", ticket.Name)
}
