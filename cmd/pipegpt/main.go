package main

import (
	"fmt"
	"os"

	"github.com/HatsuneMiku3939/pipegpt/cmd/pipegpt/cmd"
	"github.com/spf13/viper"
)

func main() {
	if err := createSubcommand(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// createSubcommand creates subcommands from configuration
func createSubcommand() error {
	var config map[string]interface{} = viper.AllSettings()

	// create subcommands
	for name, subcmd := range config {
		// skip api configuration or default configuration
		if name == "api" || name == "default" {
			continue
		}

		// otherwise, create subcommand
		definition, ok := subcmd.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid subcommand definition: %s", name)
		}

		if err := cmd.CreateSubcommand(name, definition); err != nil {
			return err
		}
	}

	return nil
}
