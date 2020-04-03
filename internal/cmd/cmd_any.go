package cmd

import (
	"github.com/kernle32dll/ew/internal"

	"github.com/fatih/color"

	"fmt"
	"io"
	"os/exec"
)

// AnyCommand is the "do-it-all" command, which allows to
// execute a user-supplied command for all paths of the given
// tags. Paths are in sorted order, but are sorted tag agnostic.
type AnyCommand struct {
	output io.Writer
	config internal.Config

	forTags []string
	args    []string
}

// NewAnyCommand creates a new AnyCommand.
func NewAnyCommand(
	output io.Writer,
	config internal.Config,
	forTags []string,
	args []string,
) *AnyCommand {
	return &AnyCommand{
		output:  output,
		config:  config,
		forTags: forTags,
		args:    args,
	}
}

func (c AnyCommand) Execute() error {
	var paths []string
	if len(c.forTags) > 0 {
		paths = c.config.GetPathsOfTagsSorted(c.forTags...)
	} else {
		paths = c.config.GetPathsOfTagsSorted(c.config.GetTagsSorted()...)
	}

	if len(paths) == 0 {
		fmt.Fprintln(c.output)
		fmt.Fprintln(c.output, determinateNoPathsErrorMessage(c.forTags))
		fmt.Fprintln(c.output)
		return nil
	}

	for _, path := range paths {
		fmt.Fprintln(c.output)
		fmt.Fprintln(c.output, colorPath("in "+path)+":")
		fmt.Fprintln(c.output)

		cmd := exec.Command(c.args[0], c.args[1:]...)
		cmd.Dir = path

		cmd.Stdout = c.output
		cmd.Stderr = c.output

		if err := cmd.Run(); err != nil {
			fmt.Fprintln(c.output, color.RedString(err.Error()))
		}
	}

	return nil
}
