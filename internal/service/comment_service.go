package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/Germanicus1/fb/api"
	"github.com/Germanicus1/fb/models"
)

// GenerateCommentID generates a unique comment ID using cryptographically secure randomness
func GenerateCommentID() string {
	// Generate 13 random bytes (will produce ~17 chars when base64 encoded)
	b := make([]byte, 13)
	if _, err := rand.Read(b); err != nil {
		// Fallback to timestamp-based ID if crypto/rand fails (extremely rare)
		return fmt.Sprintf("comment-%d", time.Now().UnixNano())
	}

	// Encode to base64 URL-safe format and remove padding
	id := base64.URLEncoding.EncodeToString(b)
	id = strings.TrimRight(id, "=")

	return id
}

// BuildCommentPayload creates a comment payload for API submission
func BuildCommentPayload(commentID, ticketID, comment string) models.CommentPayload {
	return models.CommentPayload{
		ID:       commentID,
		TicketID: ticketID,
		Comment:  comment,
	}
}

// PostComment posts a comment to a ticket
func PostComment(client *api.Client, payload models.CommentPayload) error {
	if err := client.PostComment(payload); err != nil {
		return fmt.Errorf("failed to post comment: %w", err)
	}
	return nil
}
