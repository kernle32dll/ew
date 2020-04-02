package cmd

import (
	"github.com/kernle32dll/ew/internal"

	"github.com/fatih/color"

	"fmt"
)

type ListPathsCommand struct {
	config internal.Config

	forTags []string
}

func (c ListPathsCommand) Execute() error {
	var paths []string
	if len(c.forTags) > 0 {
		paths = c.config.GetPathsOfTagsSorted(c.forTags...)
	} else {
		paths = c.config.GetPathsOfTagsSorted(c.config.GetTagsSorted()...)
	}

	if len(paths) == 0 {
		fmt.Println()
		fmt.Fprintln(color.Output, determinateNoPathsErrorMessage(c.forTags))
		fmt.Println()
		return nil
	}

	for _, path := range paths {
		fmt.Fprintln(color.Output, colorPath(path))
	}

	return nil
}
