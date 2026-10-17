package commands

import (
	"fmt"
	"io"
	"strings"
	"time"

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
}

// CreateOptions are the parsed arguments and flags of `fb create`.
type CreateOptions struct {
	Title string
	Bin   string // --bin: name or id

	Checkout bool // --checkout: make the new ticket the Checkout
	Force    bool // --force: replace an existing Checkout
}

// ExecuteCreate creates a Ticket. All validation happens before the single
// create call. Any returned error means exit status 1.
func ExecuteCreate(d CreateDeps, o CreateOptions) error {
	// 1. title
	title := strings.TrimSpace(o.Title)
	if title == "" {
		return fmt.Errorf("title required")
	}

	// 3. Checkout precondition
	if o.Checkout && !o.Force {
		if existing, err := state.LoadCheckout(); err == nil {
			return fmt.Errorf("ticket already checked out: %s\nUse 'fb clear' or 'fb checkout --force'", existing.TicketName)
		}
	}

	// 4. bin resolution
	binID, binName, err := resolveCreateBin(d, o.Bin)
	if err != nil {
		return err
	}

	// 5. type resolution
	typeID, err := resolveCreateType(d.Client, defaultTicketType)
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

	// 8. optional Checkout and Bin context save (after the ticket exists)
	if o.Checkout {
		if err := state.SaveCheckout(&state.CheckoutState{
			TicketID:     id,
			TicketName:   title,
			BinID:        binID,
			BinName:      binName,
			CheckedOutAt: time.Now(),
		}); err != nil {
			return fmt.Errorf("ticket created but saving checkout failed: %w", err)
		}
		if err := state.SaveBinContext(binID, binName); err != nil {
			return fmt.Errorf("ticket created but saving bin context failed: %w", err)
		}
	}
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
			return t.ID, nil
		}
	}
	return "", fmt.Errorf("ticket type not found: %s", want)
}
