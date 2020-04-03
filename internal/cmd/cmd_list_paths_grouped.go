package cmd

import (
	"github.com/kernle32dll/ew/internal"

	"github.com/fatih/color"

	"fmt"
	"io"
)

type ListPathsGroupedCommand struct {
	output io.Writer
	config internal.Config

	forTags []string
}

func NewListPathsGroupedCommand(
	output io.Writer,
	config internal.Config,
	forTags []string,
) *ListPathsGroupedCommand {
	return &ListPathsGroupedCommand{
		output:  output,
		config:  config,
		forTags: forTags,
	}
}

func (c ListPathsGroupedCommand) Execute() error {
	foundPaths := 0
	for _, tag := range c.forTags {
		paths := c.config.GetPathsOfTagSorted(tag)
		foundPaths += len(paths)

		fmt.Fprintln(c.output)
		color.Set(color.FgGreen, color.Bold)
		fmt.Fprintln(c.output, "@"+tag+":")
		color.Unset()

		for _, path := range paths {
			fmt.Fprintln(c.output, "    "+colorPath(path))
		}
	}

	if foundPaths == 0 {
		fmt.Fprintln(c.output)
		fmt.Fprintln(c.output, determinateNoPathsErrorMessage(c.forTags))
		fmt.Fprintln(c.output)
		return nil
	}

	return nil
}
