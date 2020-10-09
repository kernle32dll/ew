package cmd

import (
	"fmt"
	"io"
)

var (
	buildBranch string
	buildDate   string
)

// VersionCommand prints out the CLI help.
type VersionCommand struct {
	output io.Writer
}

// NewVersionCommand creates a new VersionCommand.
func NewVersionCommand(
	output io.Writer,
) *VersionCommand {
	return &VersionCommand{
		output: output,
	}
}

func (c VersionCommand) Execute() error {
	if buildBranch == "" {
		buildBranch = "master"
	}

	if buildDate != "" {
		fmt.Fprintf(c.output, "%s built %s\n", buildBranch, buildDate)
	} else {
		fmt.Fprintln(c.output, buildBranch)
	}

	return nil
}
