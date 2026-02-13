package models

import "time"

const (
	unknownStatus = "Unknown"
	dateFormat    = "2006-01-02"
)

// User represents a Flow Boards user
type User struct {
	ID    string `json:"_id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Bin represents a Flow Boards bin
type Bin struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

// Board represents a Flow Boards board
type Board struct {
	ID   string   `json:"_id"`
	Name string   `json:"name"`
	Bins []string `json:"bins"`
}

// Ticket represents a Flow Boards ticket
type Ticket struct {
	ID          string    `json:"_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BinID       string    `json:"bin_id"`
	BinName     string    `json:"bin_name"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
	DueDate     time.Time `json:"dueDate,omitempty"`
	AssignedIDs []string  `json:"assigned_ids"`
}

// Status returns the status of the ticket.
// It prioritizes BinName over BinID and returns "Unknown" if neither is set.
// This follows the Information Expert principle - the ticket knows its own status.
func (t Ticket) Status() string {
	if t.BinName != "" {
		return t.BinName
	}
	if t.BinID != "" {
		return t.BinID
	}
	return unknownStatus
}

// HasDescription returns true if the ticket has a non-empty description.
func (t Ticket) HasDescription() bool {
	return t.Description != ""
}

// FormattedCreatedDate returns the creation date in YYYY-MM-DD format.
// Returns empty string if the date is zero.
func (t Ticket) FormattedCreatedDate() string {
	return formatDate(t.CreatedAt)
}

// FormattedUpdatedDate returns the update date in YYYY-MM-DD format.
// Returns empty string if the date is zero.
func (t Ticket) FormattedUpdatedDate() string {
	return formatDate(t.UpdatedAt)
}

// FormattedDueDate returns the due date in YYYY-MM-DD format.
// Returns empty string if the date is zero.
func (t Ticket) FormattedDueDate() string {
	return formatDate(t.DueDate)
}

// formatDate converts a time.Time to YYYY-MM-DD format.
// Returns empty string if the date is zero.
func formatDate(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.Format(dateFormat)
}

// RestPrefixResponse represents the response from the REST directory endpoint
type RestPrefixResponse struct {
	RestPrefix string `json:"restUrlPrefix"`
}

// TicketSearchResponse represents the response from the ticket search endpoint
type TicketSearchResponse struct {
	Tickets []Ticket `json:"tickets"`
}

// UserSearchResponse represents the response from the user search endpoint
type UserSearchResponse struct {
	Users []User `json:"users"`
}

// HasUsers returns true if the response contains at least one user.
func (r UserSearchResponse) HasUsers() bool {
	return len(r.Users) > 0
}

// FirstUser returns the first user from the search response.
// Returns nil if no users are found.
func (r UserSearchResponse) FirstUser() *User {
	if !r.HasUsers() {
		return nil
	}
	return &r.Users[0]
}

// BinsResponse represents the paginated response from the bins endpoint
type BinsResponse struct {
	Results   []Bin  `json:"results"`
	PageToken string `json:"page-token,omitempty"`
}

// BoardsResponse represents the paginated response from the boards endpoint
type BoardsResponse struct {
	Results   []Board `json:"results"`
	PageToken string  `json:"page-token,omitempty"`
}

// CommentPayload represents the data structure for posting a comment
type CommentPayload struct {
	ID       string `json:"_id"`
	TicketID string `json:"ticket_id"`
	Comment  string `json:"comment"`
}
