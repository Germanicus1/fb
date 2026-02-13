package cli

import "fmt"

// PrintHelp prints the help message to stdout
func PrintHelp() {
	fmt.Print(GetHelpText())
}

// GetHelpText returns the help text as a string
func GetHelpText() string {
	return `fb - Flow Boards Ticket Viewer

A command-line tool for viewing your assigned Flow Boards tickets directly in the terminal.

Usage:
  fb                        Display all tickets assigned to you
  fb --bin "In Progress"    Display tickets in a specific bin
  fb --comment              Add a comment to a ticket (interactive)
  fb checkout --bin "Bin"   Check out a ticket to work on
  fb checkout TICKET-ID     Check out a specific ticket by ID
  fb -c "message"           Quick comment on checked-out ticket
  fb -o                     View currently checked-out ticket
  fb clear                  Clear checked-out ticket
  fb --version              Display version information
  fb --help                 Display this help message

Flags:
  --help                    Show this help message
  --version                 Show version information
  --bin <id or name>        Filter tickets by bin ID or bin name
  --comment                 Add a comment to a ticket (interactive)
  -c <message>              Quick comment on checked-out ticket
  -o                        View current checkout status
  --verbose                 Enable verbose output with performance metrics

Checkout Workflow:
  1. Check out a ticket:    fb checkout --bin "In Progress"
  2. Add quick comments:    fb -c "Started work"
                            fb -c "Fixed the bug"
  3. View checkout:         fb -o
  4. Clear when done:       fb clear

Examples:
  fb --bin "In Progress"           Show only tickets in the "In Progress" bin
  fb --bin kX41z9DVe               Show only tickets in the bin with ID "kX41z9DVe..."
  fb --comment                     Add a comment to a ticket (interactive)
  fb --comment --bin "In Progress" Add a comment to a ticket in the "In Progress" bin

  fb checkout --bin "Doing"        Check out a ticket from "Doing" bin
  fb checkout yL4rjYNU5PMlu7K8B    Check out specific ticket by ID
  fb -c "Making progress"          Quick comment on checked-out ticket
  fb -o                            Show which ticket is checked out
  fb clear                         Clear the checked-out ticket

Configuration:
  The tool reads configuration from ~/.fb/config.yaml

  Required configuration fields:
    auth_key:    Your Flow Boards API authentication key
    org_id:      Your organization identifier
    user_email:  Your email address for filtering tickets

Example configuration file (~/.fb/config.yaml):
  auth_key: your-api-key-here
  org_id: your-org-id
  user_email: you@example.com

`
}
