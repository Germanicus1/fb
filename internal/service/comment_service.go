package service

import (
	"fmt"

	"github.com/Germanicus1/fb/api"
	"github.com/Germanicus1/fb/models"
)

// GenerateCommentID generates a unique comment ID using cryptographically secure randomness
func GenerateCommentID() string {
	return api.GenerateID()
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
