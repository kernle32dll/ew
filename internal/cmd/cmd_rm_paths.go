package cmd

import (
	"github.com/kernle32dll/ew/internal"

	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// RmPathsCommand removes one or multiple paths from a tag.
type RmPathsCommand struct {
	output io.Writer
	config internal.Config
	args   []string
}

// NewRmPathsCommand creates a new RmPathsCommand.
func NewRmPathsCommand(
	output io.Writer,
	config internal.Config,
	args []string,
) *RmPathsCommand {
	return &RmPathsCommand{
		output: output,
		config: config,
		args:   args,
	}
}

func (c RmPathsCommand) Execute() error {
	if len(c.args) == 0 {
		return errors.New("missing @tag to remove path from")
	}

	// Current path from tag
	if len(c.args) == 1 {
		return c.rmCurrentPath()
	}

	// Multiple paths
	return c.rmPaths()
}

func (c RmPathsCommand) rmPaths() error {
	tag := c.args[len(c.args)-1]
	if !strings.HasPrefix(tag, "@") {
		return fmt.Errorf("expected @tag for last argument, got %q", tag)
	}
	tag = strings.TrimPrefix(tag, "@")

	paths := c.args[:1]

	c.config.RemovePathsFromTag(tag, paths...)

	if _, err := c.config.ReWriteConfig(); err != nil {
		return err
	}

	fmt.Fprintf(c.output, "Removed paths from tag @%s:\n", tag)
	for _, path := range paths {
		fmt.Fprintf(c.output, "  %s\n", path)
	}
	return nil
}

func (c RmPathsCommand) rmCurrentPath() error {
	tag := c.args[0]
	if !strings.HasPrefix(tag, "@") {
		return fmt.Errorf("expected @tag, got %q", tag)
	}
	tag = strings.TrimPrefix(tag, "@")

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	c.config.RemovePathsFromTag(tag, path)

	if _, err := c.config.ReWriteConfig(); err != nil {
		return err
	}

	fmt.Fprintf(c.output, "Removing path %s from tag @%s\n", path, tag)
	return nil
}
