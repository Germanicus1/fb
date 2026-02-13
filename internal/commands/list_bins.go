package commands

import (
	"fmt"

	"github.com/Germanicus1/fb/config"
	"github.com/Germanicus1/fb/internal/service"
	"github.com/Germanicus1/fb/models"
)

// ExecuteListBins lists all available bins
func ExecuteListBins(cfg *config.Config) error {
	ticketService, err := service.NewTicketService(cfg)
	if err != nil {
		return err
	}

	bins, err := ticketService.GetBins()
	if err != nil {
		return err
	}

	output := formatBinList(bins)
	fmt.Print(output)
	return nil
}

// formatBinList formats a list of bins for display
func formatBinList(bins []models.Bin) string {
	if len(bins) == 0 {
		return "No bins found.\n"
	}

	output := "Available Bins:\n\n"
	for _, bin := range bins {
		output += fmt.Sprintf("  %s - %s\n", bin.ID, bin.Name)
	}
	return output
}
