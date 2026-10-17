package cli

import (
	"strings"
	"testing"
)

func TestHelpListsCreateFlags(t *testing.T) {
	help := GetHelpText()
	for _, flag := range []string{"-d <text|->", "--bin <name or id>", "--type <name or id>", "--assign <email>", "--checkout", "--force", "-v "} {
		if !strings.Contains(help, flag) {
			t.Errorf("help text missing create flag %q", flag)
		}
	}
}
