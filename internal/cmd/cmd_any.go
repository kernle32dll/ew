package cmd

import (
	"github.com/kernle32dll/ew/internal"

	"github.com/fatih/color"

	"fmt"
	"os/exec"
)

type AnyCommand struct {
	config internal.Config

	forTags []string
	args    []string
}

func (c AnyCommand) Execute() error {
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
		fmt.Println()
		fmt.Fprintln(color.Output, colorPath("in "+path)+":")
		fmt.Println()

		cmd := exec.Command(c.args[0], c.args[1:]...)
		cmd.Dir = path

		cmd.Stdout = color.Output
		cmd.Stderr = color.Output

		if err := cmd.Run(); err != nil {
			color.Red(err.Error())
		}
	}

	return nil
}
