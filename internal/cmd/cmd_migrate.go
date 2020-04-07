package cmd

import (
	"github.com/kernle32dll/ew/internal"
	"github.com/kernle32dll/ew/internal/gr"

	"fmt"
	"io"
	"os"
	"path/filepath"
)

// MigrateCommand reads an existing gr config, and persists an
// migrated copy as a new ew config.
type MigrateCommand struct {
	output        io.Writer
	convertToYaml bool
}

// NewMigrateCommand creates a new MigrateCommand.
func NewMigrateCommand(
	output io.Writer,
	convertToYaml bool,
) *MigrateCommand {
	return &MigrateCommand{
		output:        output,
		convertToYaml: convertToYaml,
	}
}

func (c MigrateCommand) Execute() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	conf, err := gr.ParseConfigFromGr(filepath.Join(home, ".grconfig.json"))
	if err != nil {
		return err
	}

	if c.convertToYaml {
		conf.Source = internal.YamlSrc
	}

	migratedPath, err := conf.WriteConfig(home)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully migrated %q to %q\n", home+"/.grconfig.json", migratedPath)

	return nil
}
