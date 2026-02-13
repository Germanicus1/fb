// Package cli provides the command-line interface framework for Flow Boards CLI.
// It handles flag parsing, command routing, and error formatting.
package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Germanicus1/fb/config"
	"github.com/Germanicus1/fb/internal/commands"
)

// Run is the main entry point for the CLI application
func Run(version string) error {
	// Handle subcommands first (checkout, clear)
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "checkout":
			return handleCheckoutSubcommand()
		case "clear":
			return handleClearSubcommand()
		}
	}

	// Parse flags
	flags, err := parseFlags()
	if err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return fmt.Errorf("error parsing flags: %w", err)
	}

	// Handle version flag
	if flags.ShowVersion {
		fmt.Printf("fb version %s\n", version)
		return nil
	}

	// Handle help flag
	if flags.ShowHelp {
		PrintHelp()
		return nil
	}

	// Handle list-bins flag
	if flags.ListBins {
		cfg, err := loadConfiguration()
		if err != nil {
			return err
		}
		return commands.ExecuteListBins(cfg)
	}

	// Handle list-boards flag
	if flags.ListBoards {
		cfg, err := loadConfiguration()
		if err != nil {
			return err
		}
		return commands.ExecuteListBoards(cfg)
	}

	// Handle quick comment flag
	if flags.QuickComment != "" {
		return commands.ExecuteQuick(flags.QuickComment)
	}

	// Handle show status flag
	if flags.ShowStatus {
		return commands.ExecuteStatus()
	}

	// Handle bare arguments (quick comment without -c flag)
	if len(flags.Args) > 0 && !flags.CommentMode && flags.BinFilter == "" && !flags.ListBins && !flags.ListBoards {
		// Join all arguments as the comment message
		message := strings.Join(flags.Args, " ")
		return commands.ExecuteQuick(message)
	}

	// Handle comment mode
	if flags.CommentMode {
		cfg, err := loadConfiguration()
		if err != nil {
			return err
		}
		return commands.ExecuteInteractive(cfg, flags.BinFilter)
	}

	// Default: run main list command
	startTime := time.Now()

	cfg, err := loadConfiguration()
	if err != nil {
		return err
	}

	if err := commands.Execute(cfg, flags.BinFilter, flags.Verbose); err != nil {
		return err
	}

	if flags.Verbose {
		totalDuration := time.Since(startTime)
		fmt.Fprintf(os.Stderr, "\nPerformance Metrics:\n")
		fmt.Fprintf(os.Stderr, "Total execution time: %.3fs\n", totalDuration.Seconds())
	}

	return nil
}

// handleCheckoutSubcommand handles the checkout subcommand
func handleCheckoutSubcommand() error {
	fs := flag.NewFlagSet("checkout", flag.ExitOnError)
	binFlag := fs.String("bin", "", "Filter tickets by bin name")
	forceFlag := fs.Bool("force", false, "Force replace existing checkout")
	fs.Parse(os.Args[2:])

	args := fs.Args()
	return commands.ExecuteCheckout(args, *binFlag, *forceFlag)
}

// handleClearSubcommand handles the clear subcommand
func handleClearSubcommand() error {
	return commands.ExecuteClear()
}

// loadConfiguration loads and validates the application configuration
func loadConfiguration() (*config.Config, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
