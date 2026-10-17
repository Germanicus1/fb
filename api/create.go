package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/Germanicus1/fb/models"
)

// NewClientWithBaseURL creates a client aimed at an explicit base URL,
// skipping REST-prefix discovery.
func NewClientWithBaseURL(authKey, baseURL string) *Client {
	c := NewClient(authKey)
	c.baseURL = baseURL
	return c
}

// CreateTicket creates a ticket. The API answers 204 with an empty body.
func (c *Client) CreateTicket(payload models.CreateTicketPayload) error {
	if err := c.requireBaseURL(); err != nil {
		return err
	}
	if payload.AssignedIDs == nil {
		payload.AssignedIDs = []string{}
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal ticket payload: %w", err)
	}

	path := "/tickets/" + url.PathEscape(payload.ID)
	if _, err := c.doRequest("POST", path, strings.NewReader(string(jsonData))); err != nil {
		return fmt.Errorf("failed to create ticket: %w", err)
	}
	return nil
}

// ListTicketTypes retrieves all ticket types.
func (c *Client) ListTicketTypes() ([]models.TicketType, error) {
	if err := c.requireBaseURL(); err != nil {
		return nil, err
	}
	resp, err := c.doRequest(httpMethodGET, "/ticket-types", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket types: %w", err)
	}
	var types []models.TicketType
	if err := json.Unmarshal(resp, &types); err != nil {
		return nil, fmt.Errorf("failed to parse ticket types response: %w", err)
	}
	return types, nil
}

// GetTicket retrieves a single ticket by id.
func (c *Client) GetTicket(id string) (*models.Ticket, error) {
	if err := c.requireBaseURL(); err != nil {
		return nil, err
	}
	resp, err := c.doRequest(httpMethodGET, "/tickets/"+url.PathEscape(id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}
	var ticket models.Ticket
	if err := json.Unmarshal(resp, &ticket); err != nil {
		return nil, fmt.Errorf("failed to parse ticket response: %w", err)
	}
	return &ticket, nil
}
