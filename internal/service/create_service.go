package service

import (
	"fmt"
	"strings"

	"github.com/Germanicus1/fb/api"
	"github.com/Germanicus1/fb/models"
)

// ResolveBin resolves a bin given by name or id against the bin list: an
// exact id match, else a case-insensitive name match. An unlisted ref that
// looks like an id is returned as-is (the server decides; a 403 on create means
// unknown or inaccessible). Anything else is "bin not found".
func ResolveBin(client *api.Client, ref string) (models.Bin, error) {
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
	if IsBinID(ref) {
		return models.Bin{ID: ref, Name: ref}, nil
	}
	return models.Bin{}, fmt.Errorf("failed to find bin '%s': bin not found: %s", ref, ref)
}
