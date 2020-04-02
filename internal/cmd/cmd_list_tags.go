package cmd

import (
	"github.com/kernle32dll/ew/internal"

	"github.com/fatih/color"

	"fmt"
)

type ListTagsCommand struct {
	config internal.Config
}

func (c ListTagsCommand) Execute() error {
	tags := c.config.GetTagsSorted()

	if len(tags) == 0 {
		fmt.Println()
		fmt.Fprintln(color.Output, color.RedString("No tags found"))
		fmt.Println()
		return nil
	}

	for _, tag := range tags {
		fmt.Fprintln(color.Output, "@"+tag)
	}

	return nil
}
