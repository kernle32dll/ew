package cmd

import (
	"github.com/kernle32dll/ew/internal"

	"github.com/acarl005/stripansi"
	"github.com/fatih/color"

	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

const padLength = 3

// StatusCommand prints out an opinionated git short-summary
// for all paths of the given tags. Paths are in sorted order,
// but are sorted tag agnostic.
type StatusCommand struct {
	output io.Writer
	config internal.Config

	forTags []string
}

// NewStatusCommand creates a new StatusCommand.
func NewStatusCommand(
	output io.Writer,
	config internal.Config,
	forTags []string,
) *StatusCommand {
	return &StatusCommand{
		output:  output,
		config:  config,
		forTags: forTags,
	}
}

func (c StatusCommand) Execute() error {
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

	resultMatrix := make([][]string, len(paths))

	buf := &bytes.Buffer{}

	for i, path := range paths {
		resultMatrix[i] = make([]string, 4)
		resultMatrix[i][0] = colorPath(path)
		resultMatrix[i][3] = "@" + strings.Join(c.config.GetTagsOfPathSorted(path), " @")

		buf.Reset()

		cmd := exec.Command("git", "status", "--branch", "--porcelain")
		cmd.Dir = path
		cmd.Stdout = buf
		cmd.Stderr = buf

		if err := cmd.Run(); err != nil {
			fmt.Fprintln(c.output, color.RedString(err.Error()))
		}

		lines := strings.Split(
			strings.Replace(buf.String(), "\r\n", "\n", -1),
			"\n",
		)

		for i, line := range lines {
			lines[i] = strings.TrimSpace(line)
		}

		branchLine := strings.Split(strings.TrimSpace(strings.TrimPrefix(lines[0], "##")), "...")
		branchInfo := branchLine[0]

		if len(branchLine) > 1 {
			additionalBranchInfo := strings.SplitN(branchLine[1], " ", 2)
			if len(additionalBranchInfo) > 1 {
				branchInfo += " " + additionalBranchInfo[1]
			}
		} else {
			branchInfo += " " + color.RedString("(remote missing)")
		}

		resultMatrix[i][1] = branchInfo

		// -1 for branch info, -1 for empty line after branch info
		modifiedCount := len(lines) - 2

		if modifiedCount == 0 {
			resultMatrix[i][2] = color.GreenString("Clean")
		} else {
			resultMatrix[i][2] = color.RedString("%d modified", modifiedCount)
		}
	}

	resultRows, err := padMatrix(resultMatrix)
	if err != nil {
		return err
	}

	for _, row := range resultRows {
		fmt.Fprintln(c.output, row)
	}

	return nil
}

func padMatrix(matrix [][]string) ([]string, error) {
	newMatrix := make([]string, len(matrix))

	colLength := len(matrix[0])

	lengths := make([]int, colLength)
	for i, row := range matrix {
		if len(row) != colLength {
			return nil, fmt.Errorf("row %d has unexpected length %d (expected %d)", i, len(row), colLength)
		}

		for j := range matrix[i] {
			le := characterLength(row[j])
			if le > lengths[j] {
				lengths[j] = le
			}
		}
	}

	buf := bytes.Buffer{}
	for i, row := range matrix {
		buf.Reset()

		for j, col := range row {
			buf.WriteString(col)

			filler := strings.Repeat(" ", lengths[j]-characterLength(col)+padLength)
			buf.WriteString(filler)
		}

		newMatrix[i] = buf.String()
	}

	return newMatrix, nil
}

func characterLength(val string) int {
	return len(stripansi.Strip(val))
}
