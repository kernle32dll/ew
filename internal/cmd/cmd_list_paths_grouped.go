package cmd

import (
	"github.com/kernle32dll/ew/internal"

	"github.com/fatih/color"

	"fmt"
)

type ListPathsGroupedCommand struct {
	config internal.Config

	forTags []string
}

func (c ListPathsGroupedCommand) Execute() error {
	foundPaths := 0
	for _, tag := range c.forTags {
		paths := c.config.GetPathsOfTagSorted(tag)
		foundPaths += len(paths)

		fmt.Println()
		color.Set(color.FgGreen, color.Bold)
		fmt.Fprintln(color.Output, "@"+tag+":")
		color.Unset()

		for _, path := range paths {
			fmt.Fprintln(color.Output, "    "+colorPath(path))
		}
	}

	if foundPaths == 0 {
		fmt.Println()
		fmt.Fprintln(color.Output, determinateNoPathsErrorMessage(c.forTags))
		fmt.Println()
		return nil
	}

	return nil
}
