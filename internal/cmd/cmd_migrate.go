package cmd

import (
	"github.com/kernle32dll/ew/internal/gr"

	"fmt"
	"os"
)

type MigrateCommand struct {
	convertToYaml bool
}

func (c MigrateCommand) Execute() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	conf, err := gr.ParseConfigFromGr(home+"/.grconfig.json", c.convertToYaml)
	if err != nil {
		return err
	}

	migratedPath, err := conf.WriteConfig(home)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully migrated %q to %q\n", home+"/.grconfig.json", migratedPath)

	return nil
}
