package cmd

import (
	"github.com/fatih/color"
	"github.com/kernle32dll/ew/internal"

	"fmt"
	"io"
	"os/exec"
	"runtime"
)

// AnyCommand is the "do-it-all" command, which allows to
// execute a user-supplied command for all paths of the given
// tags. Paths are in sorted order, but are sorted tag agnostic.
type AnyCommand struct {
	output io.Writer
	config internal.Config

	forTags  []string
	args     []string
	parallel bool
}

// NewAnyCommand creates a new AnyCommand.
func NewAnyCommand(
	output io.Writer,
	config internal.Config,
	forTags []string,
	args []string,
	parallel bool,
) *AnyCommand {
	return &AnyCommand{
		output:   output,
		config:   config,
		forTags:  forTags,
		args:     args,
		parallel: parallel,
	}
}

func (c AnyCommand) Execute() error {
	var paths []string
	if len(c.forTags) > 0 {
		paths = c.config.GetPathsOfTagsSorted(c.forTags...)
	} else {
		paths = c.config.GetPathsOfTagsSorted(c.config.GetTagsSorted()...)
	}

	if len(paths) == 0 {
		fmt.Fprintln(c.output)
		fmt.Fprintln(c.output, determinateNoPathsErrorMessage(c.forTags))
		fmt.Fprintln(c.output)
		return nil
	}

	if c.parallel {
		c.executeParallel(paths)
		return nil
	}

	for _, path := range paths {
		fmt.Fprintln(c.output)
		fmt.Fprintln(c.output, colorPath("in "+path)+":")
		fmt.Fprintln(c.output)

		cmd := exec.Command(c.args[0], c.args[1:]...)
		cmd.Dir = path

		cmd.Stdout = c.output
		cmd.Stderr = c.output

		if err := cmd.Run(); err != nil {
			fmt.Fprintln(c.output, color.RedString(err.Error()))
		}
	}

	return nil
}

func (c AnyCommand) executeParallel(paths []string) {
	parallelFactor := runtime.NumCPU()
	in := make(chan int, parallelFactor)
	out := make([]chan *io.PipeReader, len(paths))
	for i := range paths {
		out[i] = make(chan *io.PipeReader, 1)
	}
	go func() {
		for i := range paths {
			in <- i
		}
	}()

	for i := 0; i < parallelFactor; i++ {
		go func() {
			for {
				select {
				case job, more := <-in:
					if !more {
						return
					}

					reader, writer := io.Pipe()
					out[job] <- reader

					cmd := exec.Command(c.args[0], c.args[1:]...)
					cmd.Dir = paths[job]

					cmd.Stdout = writer
					cmd.Stderr = writer

					if err := cmd.Run(); err != nil {
						fmt.Fprintln(writer, color.RedString(err.Error()))
					}

					if err := writer.Close(); err != nil {
						panic(err)
					}
				}
			}
		}()
	}

	for i, path := range paths {
		fmt.Fprintln(c.output)
		fmt.Fprintln(c.output, colorPath("in "+path)+":")
		fmt.Fprintln(c.output)

		io.Copy(c.output, <-out[i])
	}
}
