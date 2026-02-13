package formatter

import (
	"fmt"
	"strings"

	"github.com/Germanicus1/fb/models"
)

const (
	maxDescriptionLength     = 200
	maxLineWidth             = 80
	fieldIndent              = "  "     // 2 spaces for field labels
	descriptionIndent        = "    "   // 4 spaces for wrapped lines
	emptyDescriptionPlaceholder = "(none)" // Placeholder for empty descriptions
)

// FormatTicket formats a single ticket for display in the terminal
func FormatTicket(ticket models.Ticket) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Ticket ID: %s\n", ticket.ID))
	builder.WriteString(fmt.Sprintf("Ticket Name: %s\n", ticket.Name))
	builder.WriteString(fmt.Sprintf("Status: %s\n", ticket.Status()))

	if ticket.HasDescription() {
		builder.WriteString(fmt.Sprintf("Description: %s\n", ticket.Description))
	}

	return builder.String()
}

// FormatTickets formats tickets for display in the terminal
func FormatTickets(tickets []models.Ticket) string {
	if len(tickets) == 0 {
		return "No tickets assigned to you."
	}

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Found %d ticket(s) assigned to you:\n\n", len(tickets)))

	for i, ticket := range tickets {
		if i > 0 {
			builder.WriteString("\n")
		}

		formatTicketHeader(&builder, ticket)
		formatTicketStatus(&builder, ticket)
		formatTicketDates(&builder, ticket)
		formatTicketDescription(&builder, ticket)
	}

	return builder.String()
}

// formatTicketHeader writes the ticket ID and name to the builder.
func formatTicketHeader(builder *strings.Builder, ticket models.Ticket) {
	writeField(builder, "[%s] %s", ticket.ID, ticket.Name)
}

// formatTicketStatus writes the ticket status to the builder.
func formatTicketStatus(builder *strings.Builder, ticket models.Ticket) {
	writeIndentedField(builder, "Status", ticket.Status())
}

// writeField writes a formatted field to the builder.
func writeField(builder *strings.Builder, format string, args ...interface{}) {
	builder.WriteString(fmt.Sprintf(format+"\n", args...))
}

// writeIndentedField writes an indented labeled field to the builder.
func writeIndentedField(builder *strings.Builder, label, value string) {
	writeField(builder, "  %s: %s", label, value)
}

// formatTicketDates writes the created, updated, and due dates to the builder.
func formatTicketDates(builder *strings.Builder, ticket models.Ticket) {
	writeDateField(builder, "Created", ticket.FormattedCreatedDate())
	writeDateField(builder, "Updated", ticket.FormattedUpdatedDate())
	writeDateField(builder, "Due", ticket.FormattedDueDate())
}

// writeDateField writes a labeled date field to the builder if the date is present.
func writeDateField(builder *strings.Builder, label, date string) {
	if date != "" {
		writeIndentedField(builder, label, date)
	}
}

// formatTicketDescription writes the ticket description to the builder.
// Long descriptions are word-wrapped to multiple lines.
// Empty descriptions are shown as "(none)".
func formatTicketDescription(builder *strings.Builder, ticket models.Ticket) {
	description := prepareDescription(ticket.Description)
	descriptionLabel := fieldIndent + "Description: "

	// Handle empty descriptions by showing placeholder
	if description == "" {
		builder.WriteString(fmt.Sprintf("%s%s\n", descriptionLabel, emptyDescriptionPlaceholder))
		return
	}

	// Calculate available width for description text (account for label and indent)
	availableWidth := maxLineWidth - len(descriptionLabel)

	// Wrap the description text to fit within available width
	wrappedLines := wrapText(description, availableWidth)

	if len(wrappedLines) == 0 {
		return
	}

	// Write first line with label
	builder.WriteString(fmt.Sprintf("%s%s\n", descriptionLabel, wrappedLines[0]))

	// Write continuation lines with additional indentation
	for i := 1; i < len(wrappedLines); i++ {
		builder.WriteString(fmt.Sprintf("%s%s\n", descriptionIndent, wrappedLines[i]))
	}
}

// prepareDescription prepares a description for display by trimming, truncating, and normalizing.
func prepareDescription(description string) string {
	description = strings.TrimSpace(description)
	if description == "" {
		return ""
	}
	description = truncateDescription(description)
	return normalizeWhitespace(description)
}

// truncateDescription truncates long descriptions with an ellipsis.
func truncateDescription(description string) string {
	if len(description) > maxDescriptionLength {
		return description[:maxDescriptionLength] + "..."
	}
	return description
}

// normalizeWhitespace replaces newlines with spaces for compact display.
func normalizeWhitespace(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", "")
	return s
}

// wrapText wraps text to the specified width, respecting word boundaries.
// Returns a slice of lines, each no longer than maxWidth characters.
// Very long words (URLs, code) that exceed maxWidth are placed on their own line.
func wrapText(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		maxWidth = 80
	}

	// If text fits on one line, return it as-is
	if len(text) <= maxWidth {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{text}
	}

	var lines []string
	currentLine := ""

	for _, word := range words {
		if currentLine == "" {
			// First word on the line - add it regardless of length
			currentLine = word
			continue
		}

		// Check if adding this word would exceed the line width
		proposedLine := currentLine + " " + word
		if len(proposedLine) <= maxWidth {
			currentLine = proposedLine
		} else {
			// Line would be too long - save current line and start new one
			lines = append(lines, currentLine)
			currentLine = word
		}
	}

	// Add the last line
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}
