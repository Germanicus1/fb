package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SaveCheckout saves the checkout state to ~/.fb/checkout.json
func SaveCheckout(checkout *CheckoutState) error {
	homeDir, _ := os.UserHomeDir()
	fbDir := filepath.Join(homeDir, ".fb")
	os.MkdirAll(fbDir, 0700)

	data, err := json.MarshalIndent(checkout, "", "  ")
	if err != nil {
		return err
	}

	checkoutPath := filepath.Join(fbDir, "checkout.json")
	return os.WriteFile(checkoutPath, data, 0600)
}

// ClearCheckout removes the checkout state file
func ClearCheckout() error {
	checkoutPath := getCheckoutFilePath()
	if err := os.Remove(checkoutPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to clear checkout: %w", err)
	}
	return nil
}

// LoadCheckout loads the checkout state from ~/.fb/checkout.json
func LoadCheckout() (*CheckoutState, error) {
	checkoutPath := getCheckoutFilePath()
	data, err := os.ReadFile(checkoutPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no checkout file found")
		}
		return nil, fmt.Errorf("failed to read checkout file: %w", err)
	}

	var checkout CheckoutState
	if err := json.Unmarshal(data, &checkout); err != nil {
		return nil, fmt.Errorf("failed to parse checkout file: %w", err)
	}

	return &checkout, nil
}

// getCheckoutFilePath returns the path to the checkout state file
func getCheckoutFilePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".fb", "checkout.json")
}
