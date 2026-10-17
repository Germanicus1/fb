package service

import (
	"errors"
	"fmt"
	"net/http"
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

// ResolveTicketType resolves an enabled Ticket type by name
// (case-insensitive) or id.
func ResolveTicketType(client *api.Client, want string) (string, error) {
	types, err := client.ListTicketTypes()
	if err != nil {
		return "", err
	}
	for _, t := range types {
		if t.ID == want || strings.EqualFold(t.Name, want) {
			if t.Disabled {
				return "", fmt.Errorf("ticket type disabled: %s", t.Name)
			}
			return t.ID, nil
		}
	}
	return "", fmt.Errorf("ticket type not found: %s", want)
}

// ResolveAssignees returns the Assignee user ids: the given emails, or the
// user with email me when none are given. Every unknown (404) email is
// reported together.
func ResolveAssignees(client *api.Client, me string, emails []string) ([]string, error) {
	if len(emails) == 0 {
		u, err := client.GetCurrentUser(me)
		if err != nil {
			return nil, fmt.Errorf("failed to get user information: %w", err)
		}
		return []string{u.ID}, nil
	}
	var ids, unknown []string
	for _, e := range emails {
		u, err := client.GetCurrentUser(e)
		if err != nil {
			var apiErr *api.APIError
			if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
				unknown = append(unknown, e)
				continue
			}
			return nil, fmt.Errorf("failed to look up %s: %w", e, err)
		}
		ids = append(ids, u.ID)
	}
	if len(unknown) > 0 {
		return nil, fmt.Errorf("unknown assignee email(s): %s", strings.Join(unknown, ", "))
	}
	return ids, nil
}

// CreateTicket creates a Ticket under a fresh id (p.ID is ignored) and returns
// that id. On a transport failure (any error that is not an API response) the
// request may still have reached the server, so the ticket is looked up once:
// if it exists the create succeeded. API errors are returned unchanged.
func CreateTicket(client *api.Client, p models.CreateTicketPayload) (string, error) {
	p.ID = api.GenerateID()
	err := client.CreateTicket(p)
	var apiErr *api.APIError
	if err != nil && !errors.As(err, &apiErr) {
		if _, getErr := client.GetTicket(p.ID); getErr == nil {
			return p.ID, nil
		}
	}
	return p.ID, err
}
