package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/HatsuneMiku3939/pipegpt/app/generic"
	"github.com/HatsuneMiku3939/pipegpt/pkg/stdin"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CreateSubcommand creates a subcommand
func CreateSubcommand(name string, definition map[string]interface{}) error {
	// detect subcommand definition type
	definitionNames := []string{}
	for k := range definition {
		definitionNames = append(definitionNames, k)
	}

	// create subcommand
	switch determineSubcommandType(definitionNames) {
	case "generic":
		return createGenericSubcommand(name, definition)
	case "function-call":
		return createFunctionCallCommand(name, definition)
	}

	return fmt.Errorf("invalid subcommand definition: %v", definition)
}

// determineSubcommandType determines subcommand type
func determineSubcommandType(definitionNames []string) string {
	// TODO: refactor this
	switch {
	case contains(definitionNames, "role") && contains(definitionNames, "prompt") && len(definitionNames) == 2:
		return "generic"
	case contains(definitionNames, "function") && contains(definitionNames, "role") && contains(definitionNames, "prompt") && len(definitionNames) == 3:
		return "function-call"
	}

	return ""
}

// contains checks if a string slice contains a string
func contains(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}

	return false
}

// createGenericSubcommand creates a generic subcommand
func createGenericSubcommand(name string, definition map[string]interface{}) error {
	// this is not necessary for generic subcommand
	_ = definition

	// create subcommand
	subcmd := &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Ask a question with predefined role and prompt for %s task", name),
		Run: func(cmd *cobra.Command, args []string) {
			prompt := viper.GetString(fmt.Sprintf("%s.prompt", name))
			role := viper.GetString(fmt.Sprintf("%s.role", name))
			input := stdin.ConsumeStdin()

			client, err := createClient()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			result, err := generic.New(client).Run(role, prompt, input)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println(result)
		},
	}

	// add flags
	subcmd.Flags().StringP("role", "r", "",
		fmt.Sprintf("role for the AI assistant, you can also set it with PIPEGPT_%s_ROLE environment variable or config file, name", strings.ToUpper(name)),
	)
	subcmd.Flags().StringP("prompt", "p", "",
		fmt.Sprintf("prompt for the AI assistant, you can also set it with PIPEGPT_%s_PROMPT environment variable or config file, name", strings.ToUpper(name)),
	)

	// bind flags to viper
	if err := viper.BindPFlag(fmt.Sprintf("%s.role", name), subcmd.Flags().Lookup("role")); err != nil {
		return err
	}
	if err := viper.BindPFlag(fmt.Sprintf("%s.prompt", name), subcmd.Flags().Lookup("prompt")); err != nil {
		return err
	}

	// add to root command
	RootCmd.AddCommand(subcmd)
	return nil
}

// createFunctionCallCommand creates a function call subcommand
func createFunctionCallCommand(name string, definition map[string]interface{}) error {
	// TODO: implement this
	return createGenericSubcommand(name, definition)
}
