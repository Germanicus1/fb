// Package service provides business logic and API orchestration for Flow Boards operations.
// Services encapsulate API interactions and provide clean interfaces for commands to use.
package service

import (
	"fmt"

	"github.com/Germanicus1/fb/api"
	"github.com/Germanicus1/fb/config"
	"github.com/Germanicus1/fb/models"
)

// TicketService handles ticket-related operations
type TicketService struct {
	client *api.Client
	cfg    *config.Config
}

// NewTicketService creates a new ticket service with an initialized API client
func NewTicketService(cfg *config.Config) (*TicketService, error) {
	client := api.NewClient(cfg.AuthKey)

	if err := client.DiscoverRestPrefix(cfg.OrgID); err != nil {
		return nil, fmt.Errorf("failed to discover API endpoint: %w", err)
	}

	return &TicketService{
		client: client,
		cfg:    cfg,
	}, nil
}

// GetClient returns the underlying API client
func (s *TicketService) GetClient() *api.Client {
	return s.client
}

// GetCurrentUser retrieves the current user information by email
func (s *TicketService) GetCurrentUser(email string) (*models.User, error) {
	user, err := s.client.GetCurrentUser(email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user information: %w", err)
	}
	return user, nil
}

// GetUserTickets retrieves all tickets assigned to the specified user
func (s *TicketService) GetUserTickets(userID string) ([]models.Ticket, error) {
	tickets, err := s.client.SearchTickets([]string{userID})
	if err != nil {
		return nil, fmt.Errorf("failed to search tickets: %w", err)
	}
	return tickets, nil
}

// GetUserTicketsFiltered retrieves tickets with server-side filtering
func (s *TicketService) GetUserTicketsFiltered(userID, binID, boardID string) ([]models.Ticket, error) {
	tickets, err := s.client.SearchTicketsWithFilters([]string{userID}, binID, boardID)
	if err != nil {
		return nil, fmt.Errorf("failed to search tickets: %w", err)
	}
	return tickets, nil
}

// GetBins retrieves all bins
func (s *TicketService) GetBins() ([]models.Bin, error) {
	bins, err := s.client.GetBins()
	if err != nil {
		return nil, fmt.Errorf("failed to get bins: %w", err)
	}
	return bins, nil
}

// GetBoards retrieves all boards
func (s *TicketService) GetBoards() ([]models.Board, error) {
	boards, err := s.client.GetBoards()
	if err != nil {
		return nil, fmt.Errorf("failed to get boards: %w", err)
	}
	return boards, nil
}
