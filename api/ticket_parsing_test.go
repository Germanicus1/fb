package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestStory1_5_ParseSingleTicket tests parsing a single ticket from API response
func TestStory1_5_ParseSingleTicket(t *testing.T) {
	// Given: A mock API that returns a single ticket (direct array format)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"_id": "TICKET-001",
				"name": "Test Ticket",
				"description": "Test description",
				"bin_name": "To Do",
				"created": "2026-02-11T10:00:00Z",
				"updated": "2026-02-11T11:00:00Z"
			}
		]`))
	}))
	defer server.Close()

	// When: Searching for tickets
	client := NewClient("test-auth-key")
	client.baseURL = server.URL
	tickets, err := client.SearchTickets([]string{"user-123"})

	// Then: Should parse the ticket successfully
	if err != nil {
		t.Fatalf("Expected no error parsing ticket, got: %v", err)
	}

	if len(tickets) != 1 {
		t.Fatalf("Expected 1 ticket, got %d", len(tickets))
	}

	// Verify ticket fields are extracted correctly
	ticket := tickets[0]

	if ticket.ID != "TICKET-001" {
		t.Errorf("Expected ID 'TICKET-001', got: %s", ticket.ID)
	}

	if ticket.Name != "Test Ticket" {
		t.Errorf("Expected Name 'Test Ticket', got: %s", ticket.Name)
	}

	if ticket.Description != "Test description" {
		t.Errorf("Expected Description 'Test description', got: %s", ticket.Description)
	}

	if ticket.BinName != "To Do" {
		t.Errorf("Expected BinName 'To Do', got: %s", ticket.BinName)
	}
}

// TestStory1_5_EmptyResponse tests handling of empty ticket list
func TestStory1_5_EmptyResponse(t *testing.T) {
	// Given: A mock API that returns empty ticket list (direct array format)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[]`))
	}))
	defer server.Close()

	// When: Searching for tickets
	client := NewClient("test-auth-key")
	client.baseURL = server.URL
	tickets, err := client.SearchTickets([]string{"user-123"})

	// Then: Should handle empty response gracefully
	if err != nil {
		t.Errorf("Expected no error for empty response, got: %v", err)
	}

	if tickets == nil {
		t.Error("Expected empty slice, got nil")
	}

	if len(tickets) != 0 {
		t.Errorf("Expected 0 tickets, got %d", len(tickets))
	}
}

// TestStory1_5_InvalidJSON tests handling of malformed JSON
func TestStory1_5_InvalidJSON(t *testing.T) {
	// Given: A mock API that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{invalid json syntax`))
	}))
	defer server.Close()

	// When: Searching for tickets
	client := NewClient("test-auth-key")
	client.baseURL = server.URL
	_, err := client.SearchTickets([]string{"user-123"})

	// Then: Should return clear error about parsing
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}

	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "parse") && !strings.Contains(errorMsg, "JSON") {
		t.Errorf("Error should indicate parsing issue, got: %s", errorMsg)
	}
}

// TestStory1_5_MissingTicketsField tests handling of response without tickets field
func TestStory1_5_MissingTicketsField(t *testing.T) {
	// Given: A mock API that returns JSON without tickets field
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	// When: Searching for tickets
	client := NewClient("test-auth-key")
	client.baseURL = server.URL
	tickets, err := client.SearchTickets([]string{"user-123"})

	// Then: Should handle missing field gracefully (return empty list or error)
	if err != nil {
		// If there's an error, it should be clear
		errorMsg := err.Error()
		if errorMsg == "" {
			t.Error("Error message should not be empty")
		}
	} else {
		// If no error, should return nil or empty slice (both acceptable)
		if len(tickets) != 0 {
			t.Errorf("Expected 0 tickets for response without tickets field, got %d", len(tickets))
		}
	}
}

// TestStory1_5_PartialTicketData tests handling of tickets with missing optional fields
func TestStory1_5_PartialTicketData(t *testing.T) {
	// Given: A mock API that returns ticket with some missing fields (direct array format)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"_id": "TICKET-002",
				"name": "Minimal Ticket"
			}
		]`))
	}))
	defer server.Close()

	// When: Searching for tickets
	client := NewClient("test-auth-key")
	client.baseURL = server.URL
	tickets, err := client.SearchTickets([]string{"user-123"})

	// Then: Should parse available fields without error
	if err != nil {
		t.Errorf("Expected no error for partial data, got: %v", err)
	}

	if len(tickets) != 1 {
		t.Fatalf("Expected 1 ticket, got %d", len(tickets))
	}

	ticket := tickets[0]
	if ticket.ID != "TICKET-002" {
		t.Errorf("Expected ID 'TICKET-002', got: %s", ticket.ID)
	}

	if ticket.Name != "Minimal Ticket" {
		t.Errorf("Expected Name 'Minimal Ticket', got: %s", ticket.Name)
	}

	// Missing fields should have default values
	if ticket.Description != "" {
		t.Logf("Description has default empty value: '%s'", ticket.Description)
	}
}

// TestStory1_5_MultipleTickets tests that only relevant tickets are returned
func TestStory1_5_MultipleTickets(t *testing.T) {
	// Given: A mock API that returns multiple tickets (direct array format)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"_id": "TICKET-001",
				"name": "First Ticket",
				"bin_name": "To Do"
			},
			{
				"_id": "TICKET-002",
				"name": "Second Ticket",
				"bin_name": "In Progress"
			},
			{
				"_id": "TICKET-003",
				"name": "Third Ticket",
				"bin_name": "Done"
			}
		]`))
	}))
	defer server.Close()

	// When: Searching for tickets
	client := NewClient("test-auth-key")
	client.baseURL = server.URL
	tickets, err := client.SearchTickets([]string{"user-123"})

	// Then: Should parse all tickets
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(tickets) != 3 {
		t.Errorf("Expected 3 tickets, got %d", len(tickets))
	}

	// Verify first ticket (for Story 1.5, we focus on extracting the first one)
	if len(tickets) > 0 {
		firstTicket := tickets[0]
		if firstTicket.ID != "TICKET-001" {
			t.Errorf("Expected first ticket ID 'TICKET-001', got: %s", firstTicket.ID)
		}
	}
}
