package cmd

import (
	"github.com/fatih/color"

	"path/filepath"
	"strings"
)

func colorPath(path string) string {
	prefix, suffix := filepath.Split(path)
	return color.HiBlackString(prefix) + suffix
}

func determinateNoPathsErrorMessage(forTags []string) string {
	errMessage := ""
	if len(forTags) == 0 {
		errMessage = "No paths found"
	} else if len(forTags) == 1 {
		errMessage = "No paths found for tag @" + forTags[0]
	} else {
		errMessage = "No paths found for tags @" + strings.Join(forTags, ", @")
	}
	return color.RedString(errMessage)
}
