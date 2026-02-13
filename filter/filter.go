package filter

import (
	"strings"

	"github.com/Germanicus1/fb/models"
)

// FilterByBinName filters tickets by bin name or bin ID
// First tries exact match on BinID, then falls back to case-insensitive match on BinName
func FilterByBinName(tickets []models.Ticket, binFilter string) []models.Ticket {
	result := []models.Ticket{}
	lowerBinFilter := strings.ToLower(binFilter)

	for _, ticket := range tickets {
		// Try exact match on bin_id first
		if ticket.BinID == binFilter {
			result = append(result, ticket)
			continue
		}
		// Fall back to case-insensitive match on bin_name
		if strings.ToLower(ticket.BinName) == lowerBinFilter {
			result = append(result, ticket)
		}
	}

	return result
}
