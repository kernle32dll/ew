package cmd

import (
	"github.com/kernle32dll/ew/internal"

	"errors"
	"io"
	"strings"
)

// Command is an abstraction of an executable command
type Command interface {
	// Execute executes the command, and returns an error if any occurred.
	// Note, errors are always assumed fatal - non-technical problems don't warrant errors.
	Execute() error
}

// ParseCommand parses the given arguments for a command to apply, and configures it with
// the given output stream and tag config.
func ParseCommand(output io.Writer, conf internal.Config, args []string) (Command, error) {
	if len(args) == 0 {
		return NewListPathsGroupedCommand(output, conf, conf.GetTagsSorted()), nil
	}

	// Tag independent commands
	switch args[0] {
	case "help", "--help":
		return NewHelpCommand(output), nil
	case "version", "--version":
		return NewVersionCommand(output), nil
	case "migrate":
		return NewMigrateCommand(output, len(args) == 2 && args[1] == "--yaml"), nil
	case "paths":
		if len(args) == 1 || args[1] == "list" {
			return NewListPathsCommand(output, conf, nil), nil
		}
	case "tags":
		if len(args) == 1 || args[1] == "list" {
			return NewListTagsCommand(output, conf), nil
		}

		if args[1] == "add" {
			//&ListPathsCommand{config: conf}.Execute()
			return nil, errors.New("tag adding not implemented yet")
		}
		if args[1] == "rm" {
			//&ListPathsCommand{config: conf}.Execute()
			return nil, errors.New("tag adding not implemented yet")
		}
	}

	return parseTagSpecificCommands(output, conf, args)
}

func parseTagSpecificCommands(output io.Writer, conf internal.Config, args []string) (Command, error) {
	parallelFlag := false
	if args[0] == "-p" || args[0] == "--parallel" {
		parallelFlag = true
		args = args[1:]
	}

	tags, rest := parseTags(conf, args)

	if len(rest) == 0 || (len(rest) == 1 && rest[0] == "list") {
		return NewListPathsCommand(output, conf, tags), nil
	}

	if (len(rest) == 1 || len(rest) == 2) && rest[0] == "status" {
		return NewStatusCommand(output, conf, tags), nil
	}

	return NewAnyCommand(output, conf, tags, rest, parallelFlag), nil
}

func parseTags(config internal.Config, args []string) ([]string, []string) {
	tags := make(map[string]struct{}, len(args))

	spliceIndex := 0
	for _, arg := range args {
		if strings.HasPrefix(arg, "@") {
			tags[strings.TrimPrefix(arg, "@")] = struct{}{}
			spliceIndex++
		} else {
			break
		}
	}

	tagList := make([]string, len(tags))
	i := 0
	for tag := range tags {
		tagList[i] = tag
		i++
	}

	// If we don't have any tags, assume all tags
	if len(tagList) == 0 {
		tagList = config.GetTagsSorted()
	}

	return tagList, args[spliceIndex:]
}
