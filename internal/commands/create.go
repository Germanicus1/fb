package commands

import (
	"fmt"
	"io"
	"strings"

	"github.com/Germanicus1/fb/api"
	"github.com/Germanicus1/fb/internal/service"
	"github.com/Germanicus1/fb/internal/state"
	"github.com/Germanicus1/fb/models"
)

// defaultTicketType is the Ticket type used when none is chosen.
const defaultTicketType = "General"

// CreateDeps are the injected dependencies of the create command: the command
// test seam. Production wiring builds these from config and the real terminal.
// The user's home directory is read from HOME by the state package.
type CreateDeps struct {
	Client     *api.Client
	UserEmail  string
	Stdin      io.Reader
	Stdout     io.Writer
	IsTerminal bool

	// ConfigTicketType is the optional `ticket_type` config key.
	ConfigTicketType string
}

// CreateOptions are the parsed arguments and flags of `fb create`.
type CreateOptions struct {
	Title string
	Bin   string // --bin: name or id
	Type  string // --type: Ticket type name or id
}

// ExecuteCreate creates a Ticket. All validation happens before the single
// create call. Any returned error means exit status 1.
func ExecuteCreate(d CreateDeps, o CreateOptions) error {
	// 1. title
	title := strings.TrimSpace(o.Title)
	if title == "" {
		return fmt.Errorf("title required")
	}

	// 4. bin resolution
	binID, binName, err := resolveCreateBin(d, o.Bin)
	if err != nil {
		return err
	}

	// 5. type resolution
	typeID, err := resolveCreateType(d.Client, firstNonEmpty(o.Type, d.ConfigTicketType, defaultTicketType))
	if err != nil {
		return err
	}

	// 6. assignee resolution (me)
	me, err := d.Client.GetCurrentUser(d.UserEmail)
	if err != nil {
		return fmt.Errorf("failed to get user information: %w", err)
	}

	// 7. create
	id := api.GenerateID()
	if err := d.Client.CreateTicket(models.CreateTicketPayload{
		ID:           id,
		Name:         title,
		TicketTypeID: typeID,
		BinID:        binID,
		AssignedIDs:  []string{me.ID},
	}); err != nil {
		return err
	}

	// 9. confirmation
	fmt.Fprintf(d.Stdout, "✓ Created %s - %s  [%s]\n", id, title, binName)
	return nil
}

// resolveCreateBin returns the bin id and display name: --bin (name or id),
// otherwise the Bin context, otherwise an error naming --bin.
func resolveCreateBin(d CreateDeps, binFlag string) (id, name string, err error) {
	if binFlag == "" {
		bc, loadErr := state.LoadBinContext()
		if loadErr != nil || bc.BinID == "" {
			return "", "", fmt.Errorf("no bin specified: use --bin <name or id>")
		}
		return bc.BinID, bc.BinName, nil
	}

	id, err = service.ResolveBinFilter(d.Client, binFlag)
	if err != nil {
		return "", "", err
	}
	name = binFlag
	if bins, listErr := d.Client.GetBins(); listErr == nil {
		for _, b := range bins {
			if b.ID == id {
				name = b.Name
				break
			}
		}
	}
	return id, name, nil
}

// resolveCreateType resolves a Ticket type by name (case-insensitive) or id.
func resolveCreateType(c *api.Client, want string) (string, error) {
	types, err := c.ListTicketTypes()
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

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
