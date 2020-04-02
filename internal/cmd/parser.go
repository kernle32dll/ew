package cmd

import (
	"github.com/kernle32dll/ew/internal"

	"errors"
	"strings"
)

type Command interface {
	Execute() error
}

func ParseCommand(conf internal.Config, args []string) (Command, error) {
	if len(args) == 0 {
		return &ListPathsGroupedCommand{
			config:  conf,
			forTags: conf.GetTagsSorted(),
		}, nil
	}

	if args[0] == "help" {
		return &HelpCommand{}, nil
	}

	if args[0] == "migrate" {
		return &MigrateCommand{
			convertToYaml: len(args) == 2 && args[1] == "--yaml",
		}, nil
	}

	// List all paths
	if args[0] == "paths" {
		if len(args) == 1 || args[1] == "list" {
			return &ListPathsCommand{
				config: conf,
			}, nil
		}
	}

	if args[0] == "tags" {
		if len(args) == 1 || args[1] == "list" {
			return &ListTagsCommand{
				config: conf,
			}, nil
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

	tags, rest := parseTags(conf, args)

	if len(rest) == 0 || (len(rest) == 1 && rest[0] == "list") {
		return &ListPathsCommand{
			config:  conf,
			forTags: tags,
		}, nil
	}

	if (len(rest) == 1 || len(rest) == 2) && rest[0] == "status" {
		if len(rest) == 2 && rest[1] == "-v" {
			return &AnyCommand{
				config:  conf,
				forTags: tags,
				args:    []string{"git", "status", "-sb"},
			}, nil
		}

		return &StatusCommand{
			config:  conf,
			forTags: tags,
		}, nil
	}

	return &AnyCommand{
		config:  conf,
		forTags: tags,
		args:    rest,
	}, nil
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
