package cmd

import (
	"github.com/kernle32dll/ew/internal"

	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// AddPathsCommand adds one or multiple paths to a tag.
type AddPathsCommand struct {
	output io.Writer
	config internal.Config
	args   []string
}

// NewAddPathsCommand creates a new AddPathsCommand.
func NewAddPathsCommand(
	output io.Writer,
	config internal.Config,
	args []string,
) *AddPathsCommand {
	return &AddPathsCommand{
		output: output,
		config: config,
		args:   args,
	}
}

func (c AddPathsCommand) Execute() error {
	if len(c.args) == 0 {
		return errors.New("missing @tag to add path to")
	}

	// Current path to tag
	if len(c.args) == 1 {
		return c.addCurrentPath()
	}

	// Multiple paths
	return c.addPaths()
}

func (c AddPathsCommand) addPaths() error {
	tag := c.args[len(c.args)-1]
	if !strings.HasPrefix(tag, "@") {
		return fmt.Errorf("expected @tag for last argument, got %q", tag)
	}
	tag = strings.TrimPrefix(tag, "@")

	paths := c.args[:1]

	c.config.AddPathsToTag(tag, paths...)

	if _, err := c.config.ReWriteConfig(); err != nil {
		return err
	}

	fmt.Fprintf(c.output, "Added paths to tag @%s:\n", tag)
	for _, path := range paths {
		fmt.Fprintf(c.output, "  %s\n", path)
	}
	return nil
}

func (c AddPathsCommand) addCurrentPath() error {
	tag := c.args[0]
	if !strings.HasPrefix(tag, "@") {
		return fmt.Errorf("expected @tag, got %q", tag)
	}
	tag = strings.TrimPrefix(tag, "@")

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	c.config.AddPathsToTag(tag, path)

	if _, err := c.config.ReWriteConfig(); err != nil {
		return err
	}

	fmt.Fprintf(c.output, "Added path %s to tag @%s\n", path, tag)
	return nil
}
