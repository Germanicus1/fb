package service

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/Germanicus1/fb/api"
	"github.com/Germanicus1/fb/models"
)

var errBinNotFound = errors.New("bin not found")

// LookupBin finds a bin given by name or id in the bin list: an exact id
// match, else a case-insensitive name match. A ref not in the list is
// "bin not found", even if it looks like an id.
func LookupBin(client *api.Client, ref string) (models.Bin, error) {
	bins, err := client.GetBins()
	if err != nil {
		return models.Bin{}, err
	}
	for _, b := range bins {
		if b.ID == ref {
			return b, nil
		}
	}
	for _, b := range bins {
		if strings.EqualFold(b.Name, ref) {
			return b, nil
		}
	}
	return models.Bin{}, fmt.Errorf("failed to find bin '%s': %w: %s", ref, errBinNotFound, ref)
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
