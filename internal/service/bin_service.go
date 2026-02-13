package service

import (
	"fmt"
	"unicode"

	"github.com/Germanicus1/fb/api"
)

// ResolveBinFilter converts a bin name to a bin ID.
// If the input is already a bin ID (alphanumeric only), it returns it unchanged.
// Otherwise, it performs a case-insensitive lookup to find the matching bin ID.
func ResolveBinFilter(client *api.Client, binFilter string) (string, error) {
	if IsBinID(binFilter) {
		return binFilter, nil
	}

	binID, err := client.LookupBinIDByName(binFilter)
	if err != nil {
		return "", fmt.Errorf("failed to find bin '%s': %w", binFilter, err)
	}
	return binID, nil
}

// IsBinID determines if a string is a bin ID based on its format.
// Bin IDs are alphanumeric strings without spaces or special characters.
// Bin names typically contain spaces, dots, or special characters like '+'.
func IsBinID(s string) bool {
	if s == "" {
		return false
	}

	for _, c := range s {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return false
		}
	}

	return true
}
