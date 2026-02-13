package cli

import (
	"flag"
	"os"
)

// Flags represents all CLI flags
type Flags struct {
	ShowVersion  bool
	ShowHelp     bool
	BinFilter    string
	ListBins     bool
	ListBoards   bool
	CommentMode  bool
	QuickComment string
	ShowStatus   bool
	Verbose      bool
	Args         []string
}

// parseFlags parses command line flags and returns a Flags struct
func parseFlags() (*Flags, error) {
	flags := &Flags{}

	// Create a new flag set to avoid conflicts in tests
	fs := flag.NewFlagSet("fb", flag.ContinueOnError)
	fs.BoolVar(&flags.ShowVersion, "version", false, "Display version information")
	fs.BoolVar(&flags.ShowHelp, "help", false, "Display help message")
	fs.StringVar(&flags.BinFilter, "bin", "", "Filter tickets by bin name")
	fs.BoolVar(&flags.ListBins, "list-bins", false, "List all available bins")
	fs.BoolVar(&flags.ListBoards, "list-boards", false, "List all available boards")
	fs.BoolVar(&flags.CommentMode, "comment", false, "Add a comment to a ticket")
	fs.StringVar(&flags.QuickComment, "c", "", "Quick comment on checked-out ticket")
	fs.BoolVar(&flags.ShowStatus, "o", false, "View current checkout status")
	fs.BoolVar(&flags.Verbose, "verbose", false, "Enable verbose output")
	fs.BoolVar(&flags.Verbose, "v", false, "Enable verbose output (short flag)")
	fs.BoolVar(&flags.Verbose, "debug", false, "Enable debug output")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	flags.Args = fs.Args()
	return flags, nil
}
