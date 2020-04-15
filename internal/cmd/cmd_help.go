package cmd

import (
	"github.com/fatih/color"

	"fmt"
	"io"
	"strings"
)

// HelpCommand prints out the CLI help.
type HelpCommand struct {
	output io.Writer
}

// NewHelpCommand creates a new HelpCommand.
func NewHelpCommand(
	output io.Writer,
) *HelpCommand {
	return &HelpCommand{
		output: output,
	}
}

func (c HelpCommand) Execute() error {
	fmt.Fprintln(c.output, "EW - (run things) e(very)w(here)")
	fmt.Fprintln(c.output, "Original author: Björn Gerdau")

	fmt.Fprint(c.output, "Version: ")
	NewVersionCommand(c.output).Execute()

	fmt.Fprintln(c.output)
	fmt.Fprintln(c.output, "Available commands:")

	nyi := color.RedString(" (NOT YET IMPLEMENTED)")
	cmds := [][]string{
		{"ew", "list all paths, grouped by their tags"},
		{"ew help", "displays this help"},
		{"ew --help", ""},
		{"ew version", "displays version and build time (if applicable)"},
		{"ew --version", ""},
		{"ew migrate", "migrate from mixu/gr config, and keep json format"},
		{"ew migrate --yaml", "migrate from mixu/gr config, and use new yaml format"},
		{"ew paths", "list all paths (alias for ew paths list)"},
		{"ew paths list", "list all paths"},
		{"ew tags", "list all tags (alias for ew tags list)"},
		{"ew tags list", "list all tags"},
		{"ew tags add @some-tag", "add current directory to tag \"some-tag\"" + nyi},
		{"ew tags add \\some\\path @some-tag", "add \\some\\path to tag \"some-tag\"" + nyi},
		{"ew tags rm @some-tag", "add current directory to tag \"some-tag\"" + nyi},
		{"ew tags rm \\some\\path @some-tag", "add \\some\\path to tag \"some-tag\"" + nyi},
		{"ew status", "show quick git status for all paths"},
		{"ew @tag1 status", "show quick git status for all paths of tag1 (supports multiple tags)"},
		{"ew @tag1 some-cmd", "executes some-cmd in all paths of tag1 (supports multiple tags)"},
		{"ew --parallel @tag1 some-cmd", "parallel executes some-cmd in all paths of tag1 (supports multiple tags)"},
		{"ew -p @tag1 some-cmd", "(note: You will probably loose colored output!)"},
	}

	longestHelp := 0
	for _, cmd := range cmds {
		if le := len(cmd[0]); le > longestHelp {
			longestHelp = le
		}
	}

	longestHelp += 5

	for _, cmd := range cmds {
		filler := strings.Repeat(" ", longestHelp-len(cmd[0]))
		fmt.Fprintln(c.output, "  "+cmd[0]+filler+cmd[1])
	}

	return nil
}
