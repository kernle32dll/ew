package cmd

import (
	"github.com/kernle32dll/ew/internal"

	"github.com/fatih/color"

	"fmt"
	"io"
)

// ListTagsCommand lists all tags in sorted order.
type ListTagsCommand struct {
	output io.Writer
	config internal.Config
}

// NewListTagsCommand creates a new ListTagsCommand.
func NewListTagsCommand(
	output io.Writer,
	config internal.Config,
) *ListTagsCommand {
	return &ListTagsCommand{
		output: output,
		config: config,
	}
}

func (c ListTagsCommand) Execute() error {
	tags := c.config.GetTagsSorted()

	if len(tags) == 0 {
		fmt.Fprintln(c.output)
		fmt.Fprintln(c.output, color.RedString("No tags found"))
		fmt.Fprintln(c.output)
		return nil
	}

	for _, tag := range tags {
		fmt.Fprintln(c.output, "@"+tag)
	}

	return nil
}
