package cmd

import (
	"github.com/kernle32dll/ew/internal"

	"github.com/fatih/color"

	"bytes"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"sync"
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
	fmt.Fprintf(c.output, "Executing in parallel by %d, please wait...\n\n", parallelFactor)

	wg := &sync.WaitGroup{}
	in := make(chan int, parallelFactor)
	out := make([]*bytes.Buffer, len(paths))
	for i := 0; i < parallelFactor; i++ {
		go func() {
			for {
				select {
				case job, more := <-in:
					if !more {
						return
					}

					buf := &bytes.Buffer{}

					cmd := exec.Command(c.args[0], c.args[1:]...)
					cmd.Dir = paths[job]

					cmd.Stdout = buf
					cmd.Stderr = buf

					if err := cmd.Run(); err != nil {
						fmt.Fprintln(buf, color.RedString(err.Error()))
					}

					out[job] = buf
					wg.Done()
				}
			}
		}()
	}

	wg.Add(len(paths))
	for i := range paths {
		in <- i
	}

	wg.Wait()
	close(in)

	for i, path := range paths {
		fmt.Fprintln(c.output)
		fmt.Fprintln(c.output, colorPath("in "+path)+":")
		fmt.Fprintln(c.output)

		out[i].WriteTo(c.output)
	}
}
