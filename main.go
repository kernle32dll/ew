package main

import (
	"github.com/fatih/color"
	"github.com/kernle32dll/ew/internal"
	"github.com/kernle32dll/ew/internal/cmd"

	"fmt"
	"os"
	"time"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error determinating home dir: %s", err)
		os.Exit(1)
	}

	conf := internal.ParseConfigFromFolder(color.Output, home)

	exec, err := cmd.ParseCommand(color.Output, conf, os.Args[1:])
	if err != nil {
		fmt.Printf("Error while parsing cmd: %s\n", err)
		os.Exit(1)
	}

	start := time.Now()
	defer func() {
		fmt.Printf("\nExecuted in %s\n", time.Since(start))
	}()

	if err := exec.Execute(); err != nil {
		fmt.Printf("Error while executing cmd: %s\n", err)
		os.Exit(1)
	}
}
