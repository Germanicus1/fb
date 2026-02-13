package commands

import (
	"fmt"

	"github.com/Germanicus1/fb/config"
	"github.com/Germanicus1/fb/internal/service"
	"github.com/Germanicus1/fb/models"
)

// ExecuteListBoards lists all available boards
func ExecuteListBoards(cfg *config.Config) error {
	ticketService, err := service.NewTicketService(cfg)
	if err != nil {
		return err
	}

	boards, err := ticketService.GetBoards()
	if err != nil {
		return err
	}

	output := formatBoardList(boards)
	fmt.Print(output)
	return nil
}

// formatBoardList formats a list of boards for display
func formatBoardList(boards []models.Board) string {
	if len(boards) == 0 {
		return "No boards found.\n"
	}

	output := "Available Boards:\n\n"
	for _, board := range boards {
		output += fmt.Sprintf("  %s - %s\n", board.ID, board.Name)
	}
	return output
}
