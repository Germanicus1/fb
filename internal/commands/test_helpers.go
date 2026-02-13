package commands

import (
	"io"
	"os"
	"path/filepath"

	"github.com/Germanicus1/fb/internal/state"
)

// Test helper functions

// runQuickComment is a test helper for quick comment functionality
func runQuickComment(output io.Writer, comment string) error {
	return ExecuteQuick(comment)
}

// getCheckoutFilePathForRead returns the path to the checkout state file for testing
func getCheckoutFilePathForRead() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".fb", "checkout.json"), nil
}

// loadCheckoutStateTest loads checkout state for testing
func loadCheckoutStateTest() (*state.CheckoutState, error) {
	return state.LoadCheckout()
}
