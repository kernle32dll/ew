package cmd

import (
	"github.com/kernle32dll/ew/internal"

	"fmt"
	"io"
)

// ListPathsCommand lists all paths of the given tags,
// in sorted order (but order is tag agnostic).
type ListPathsCommand struct {
	output io.Writer
	config internal.Config

	forTags []string
}

// NewListPathsCommand creates a new ListPathsCommand.
func NewListPathsCommand(
	output io.Writer,
	config internal.Config,
	forTags []string,
) *ListPathsCommand {
	return &ListPathsCommand{
		output:  output,
		config:  config,
		forTags: forTags,
	}
}

func (c ListPathsCommand) Execute() error {
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
		fmt.Fprintln(c.output, colorPath(path))
	}

	return nil
}
