package commands

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Germanicus1/fb/api"
	"github.com/Germanicus1/fb/formatter"
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

	Checkout bool // --checkout: make the new ticket the Checkout
	Force    bool // --force: replace an existing Checkout
	Verbose  bool // -v: also print the created ticket
	// Description is -d: text sent verbatim, or "-" to read all of stdin.
	Description string
	Assign      []string // --assign emails (repeatable); replaces me
	Type        string   // --type: Ticket type name or id
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
	typeID, err := resolveCreateType(d.Client, firstNonEmpty(o.Type, d.ConfigTicketType, defaultTicketType))
	if err != nil {
		return err
	}

	// 6. assignee resolution (me unless --assign)
	assignedIDs, err := resolveCreateAssignees(d, o.Assign)
	if err != nil {
		return err
	}

	description := o.Description
	if description == "-" {
		b, readErr := io.ReadAll(d.Stdin)
		if readErr != nil {
			return fmt.Errorf("failed to read description from stdin: %w", readErr)
		}
		description = string(b)
	}

	// 7. create
	id := api.GenerateID()
	if err := d.Client.CreateTicket(models.CreateTicketPayload{
		ID:           id,
		Name:         title,
		TicketTypeID: typeID,
		BinID:        binID,
		Description:  description,
		AssignedIDs:  assignedIDs,
	}); err != nil {
		var apiErr *api.APIError
		var netErr *url.Error
		switch {
		case errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusForbidden:
			return fmt.Errorf("bin not found or no access: %s", binRef(o.Bin, binID))
		case errors.As(err, &netErr):
			// The request may have reached the server: check before reporting.
			if _, getErr := d.Client.GetTicket(id); getErr != nil {
				return err
			}
		default:
			return err
		}
	}

	// 9. confirmation
	fmt.Fprintf(d.Stdout, "✓ Created %s - %s  [%s]\n", id, title, binName)

	if o.Verbose {
		if t, err := d.Client.GetTicket(id); err == nil {
			fmt.Fprint(d.Stdout, "\n"+formatter.FormatTicket(*t))
		}
	}

	// 8. optional Checkout and Bin context save (after the ticket exists)
	if o.Checkout {
		if err := state.SaveCheckout(&state.CheckoutState{
			TicketID:     id,
			TicketName:   title,
			BinID:        binID,
			BinName:      binName,
			CheckedOutAt: time.Now().Format(time.RFC3339),
		}); err != nil {
			return fmt.Errorf("ticket created but saving checkout failed: %w", err)
		}
		if err := state.SaveBinContext(binID, binName); err != nil {
			return fmt.Errorf("ticket created but saving bin context failed: %w", err)
		}
	}
	return nil
}

// binRef is how the user referred to the bin, falling back to its id.
func binRef(flag, id string) string {
	if flag != "" {
		return flag
	}
	return id
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

// resolveCreateAssignees returns the Assignee user ids: the --assign emails
// (replacing me), or me when none given. Every unknown (404) email is
// reported together.
func resolveCreateAssignees(d CreateDeps, emails []string) ([]string, error) {
	if len(emails) == 0 {
		me, err := d.Client.GetCurrentUser(d.UserEmail)
		if err != nil {
			return nil, fmt.Errorf("failed to get user information: %w", err)
		}
		return []string{me.ID}, nil
	}
	var ids, unknown []string
	for _, e := range emails {
		u, err := d.Client.GetCurrentUser(e)
		if err != nil {
			if strings.Contains(err.Error(), "status 404") {
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
