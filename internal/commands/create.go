package commands

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Germanicus1/fb/api"
	"github.com/Germanicus1/fb/formatter"
	"github.com/Germanicus1/fb/internal/prompt"
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
	Stderr     io.Writer // warnings; nil discards
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
	in := prompt.NewLineReader(d.Stdin, d.Stdout)
	interactive := o.Title == ""

	// 1. a title is required unless we can prompt for it
	if interactive && !d.IsTerminal {
		return fmt.Errorf("title required")
	}

	// 2. Checkout precondition (before any prompt)
	if o.Checkout && !o.Force {
		if existing, err := state.LoadCheckout(); err == nil {
			return fmt.Errorf("ticket already checked out: %s\nUse 'fb clear' or 'fb checkout --force'", existing.TicketName)
		}
	}

	// 3. interactive title and description
	if interactive {
		if err := promptTitleAndDescription(in, d.Stdout, &o); err != nil {
			return err
		}
	}

	// 4. resolve everything the create call needs
	p, bin, err := buildCreatePayload(d, in, o, interactive)
	if err != nil {
		return err
	}

	// 5. create
	id, err := service.CreateTicket(d.Client, p)
	if err != nil {
		var apiErr *api.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusForbidden {
			return fmt.Errorf("bin not found or no access: %s", bin.Name)
		}
		return err
	}

	// 6. confirmation
	fmt.Fprintf(d.Stdout, "✓ Created %s - %s  [%s]\n", id, p.Name, bin.Name)
	if o.Verbose {
		printCreatedTicket(d, id)
	}

	// 7. optional Checkout and Bin context save (after the ticket exists)
	if o.Checkout {
		return saveCreatedCheckout(id, p.Name, bin)
	}
	return nil
}

// promptTitleAndDescription prompts for the title and, unless -d was given,
// the description.
func promptTitleAndDescription(in prompt.LineReader, out io.Writer, o *CreateOptions) error {
	t, err := promptTitle(in, out)
	if err != nil {
		return err
	}
	o.Title = t
	if o.Description == "" {
		if o.Description, err = promptDescription(in, out); err != nil {
			return err
		}
	}
	return nil
}

// buildCreatePayload validates the title and resolves bin, Ticket type,
// Assignees and description. The payload's ID is left for the service to set.
func buildCreatePayload(d CreateDeps, in prompt.LineReader, o CreateOptions, interactive bool) (models.CreateTicketPayload, models.Bin, error) {
	var p models.CreateTicketPayload
	title := strings.TrimSpace(o.Title)
	if title == "" {
		return p, models.Bin{}, fmt.Errorf("title required")
	}

	bin, err := resolveCreateBin(d, in, interactive, o.Bin)
	if err != nil {
		return p, models.Bin{}, err
	}

	typeID, err := service.ResolveTicketType(d.Client, firstNonEmpty(o.Type, d.ConfigTicketType, defaultTicketType))
	if err != nil {
		return p, models.Bin{}, err
	}

	assignedIDs, err := service.ResolveAssignees(d.Client, d.UserEmail, o.Assign)
	if err != nil {
		return p, models.Bin{}, err
	}

	description := o.Description
	if description == "-" {
		b, readErr := io.ReadAll(d.Stdin)
		if readErr != nil {
			return p, models.Bin{}, fmt.Errorf("failed to read description from stdin: %w", readErr)
		}
		description = string(b)
	}

	return models.CreateTicketPayload{
		Name:         title,
		TicketTypeID: typeID,
		BinID:        bin.ID,
		Description:  description,
		AssignedIDs:  assignedIDs,
	}, bin, nil
}

// printCreatedTicket prints the created ticket (-v). A failed fetch is only a
// warning: the ticket was created.
func printCreatedTicket(d CreateDeps, id string) {
	t, err := d.Client.GetTicket(id)
	if err != nil {
		if d.Stderr != nil {
			fmt.Fprintf(d.Stderr, "Warning: could not fetch created ticket: %v\n", err)
		}
		return
	}
	fmt.Fprint(d.Stdout, "\n"+formatter.FormatTicket(*t))
}

// saveCreatedCheckout makes the new ticket the Checkout and its bin the Bin
// context.
func saveCreatedCheckout(id, title string, bin models.Bin) error {
	if err := state.SaveCheckout(&state.CheckoutState{
		TicketID:     id,
		TicketName:   title,
		BinID:        bin.ID,
		BinName:      bin.Name,
		CheckedOutAt: time.Now().Format(time.RFC3339),
	}); err != nil {
		return fmt.Errorf("ticket created but saving checkout failed: %w", err)
	}
	if err := state.SaveBinContext(bin.ID, bin.Name); err != nil {
		return fmt.Errorf("ticket created but saving bin context failed: %w", err)
	}
	return nil
}

// resolveCreateBin returns the bin: --bin (name or id); else, in interactive
// mode, the picker with the Bin context as default; else the Bin context; else
// the picker at a terminal; else an error naming --bin.
func resolveCreateBin(d CreateDeps, in prompt.LineReader, interactive bool, binFlag string) (models.Bin, error) {
	if binFlag != "" {
		return service.ResolveBin(d.Client, binFlag)
	}
	bc, loadErr := state.LoadBinContext()
	hasCtx := loadErr == nil && bc.BinID != ""
	if hasCtx && !interactive {
		return models.Bin{ID: bc.BinID, Name: bc.BinName}, nil
	}
	if !d.IsTerminal {
		return models.Bin{}, fmt.Errorf("no bin specified: use --bin <name or id>")
	}
	var def *models.Bin
	if hasCtx {
		def = &models.Bin{ID: bc.BinID, Name: bc.BinName}
	}
	bins, err := d.Client.GetBins()
	if err != nil {
		return models.Bin{}, err
	}
	return pickBin(in, d.Stdout, bins, def)
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

// readLine shows p and reads one line. eof is true when input ended (or failed)
// while reading it; line may still hold a final partial line. err is set only
// when the user pressed Ctrl-C.
func readLine(in prompt.LineReader, p string) (line string, eof bool, err error) {
	line, readErr := in.ReadLine(p)
	if errors.Is(readErr, prompt.ErrInterrupted) {
		return "", true, readErr
	}
	return line, readErr != nil, nil
}

// promptTitle asks for the title until it is non-empty. EOF is "title required".
func promptTitle(in prompt.LineReader, out io.Writer) (string, error) {
	for {
		line, eof, err := readLine(in, "Title: ")
		if err != nil {
			return "", err
		}
		if t := strings.TrimSpace(line); t != "" {
			return t, nil
		}
		if eof {
			fmt.Fprintln(out)
			return "", fmt.Errorf("title required")
		}
	}
}

// promptDescription reads lines until an empty line or EOF; an empty first
// line means no description.
func promptDescription(in prompt.LineReader, out io.Writer) (string, error) {
	fmt.Fprintln(out, "Description (optional; finish with an empty line):")
	var lines []string
	for {
		line, eof, err := readLine(in, "")
		if err != nil {
			return "", err
		}
		if line == "" {
			break
		}
		lines = append(lines, line)
		if eof {
			break
		}
	}
	return strings.Join(lines, "\n"), nil
}

const maxPickList = 20

var errNoBinSelected = errors.New("no bin selected")

// pickBin prompts for a bin until one is chosen. Empty input picks def (if
// any); otherwise input filters bins by case-insensitive substring.
func pickBin(in prompt.LineReader, out io.Writer, bins []models.Bin, def *models.Bin) (models.Bin, error) {
	for {
		p := "Bin: "
		if def != nil {
			p = fmt.Sprintf("Bin [%s]: ", def.Name)
		}
		line, eof, err := readLine(in, p)
		if err != nil {
			return models.Bin{}, err
		}
		input := strings.TrimSpace(line)
		if input == "" {
			if def != nil {
				return *def, nil
			}
			if eof {
				fmt.Fprintln(out)
				return models.Bin{}, errNoBinSelected
			}
			continue
		}
		var matches []models.Bin
		lower := strings.ToLower(input)
		for _, b := range bins {
			if strings.Contains(strings.ToLower(b.Name), lower) {
				matches = append(matches, b)
			}
		}
		switch {
		case len(matches) == 0:
			fmt.Fprintf(out, "No bin matches %q.\n", input)
		case len(matches) == 1:
			fmt.Fprintf(out, "→ %s\n", matches[0].Name)
			return matches[0], nil
		case len(matches) > maxPickList:
			fmt.Fprintf(out, "%d bins match %q — type more.\n", len(matches), input)
		default:
			if b, ok, nErr := numberedPick(in, out, matches); nErr != nil {
				return models.Bin{}, nErr
			} else if ok {
				return b, nil
			}
		}
		if eof {
			return models.Bin{}, errNoBinSelected
		}
	}
}

// numberedPick lists matches and reads a number. An invalid number re-prompts;
// an empty line returns ok=false so the caller goes back to filtering.
func numberedPick(in prompt.LineReader, out io.Writer, matches []models.Bin) (models.Bin, bool, error) {
	for i, b := range matches {
		fmt.Fprintf(out, "  %d) %s\n", i+1, b.Name)
	}
	for {
		line, eof, err := readLine(in, fmt.Sprintf("Number (1-%d): ", len(matches)))
		if err != nil {
			return models.Bin{}, false, err
		}
		input := strings.TrimSpace(line)
		if input == "" {
			if eof {
				return models.Bin{}, false, errNoBinSelected
			}
			return models.Bin{}, false, nil
		}
		if n, convErr := strconv.Atoi(input); convErr == nil && n >= 1 && n <= len(matches) {
			fmt.Fprintf(out, "→ %s\n", matches[n-1].Name)
			return matches[n-1], true, nil
		}
		fmt.Fprintln(out, "Invalid number.")
		if eof {
			return models.Bin{}, false, errNoBinSelected
		}
	}
}
