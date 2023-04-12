package cmd

import (
	"fmt"
	"os"

	"github.com/HatsuneMiku3939/pipegpt/pkg/chatgpt"
	"github.com/HatsuneMiku3939/pipegpt/pkg/stdin"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pipegpt",
	Short: "A simple CLI tools to question chatgpt",
	Long: `pipegpt is a simple CLI tools to question chatgpt.

You can use it to question chatgpt with your pre-defined prompt or interactively (in cli manner).
pipegpt always wait for your input from stdin and output the result to stdout.
Provided prompt is used as the first input to chatgpt.

Simple Example:
# code review in staged changes
git diff --staged | pipegpt -p "code review for this change"

# commit with suggested commit message
git diff --staged  | pipegpt -p "suggest conventional commit messages for these changes" | git commit -F - -e

# convert JSON to YAML
cat sample.json | pipegpt -p "convert JSON to YAML"
`,
	Run: func(cmd *cobra.Command, args []string) {
		prompt, err := cmd.Flags().GetString("prompt")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		apikey := os.Getenv("OPENAI_API_KEY")
		userInput := stdin.ConsumeStdin()

		result, err := chatgpt.Chat(apikey, prompt, userInput)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(result)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	apikey := os.Getenv("OPENAI_API_KEY")
	if apikey == "" {
		fmt.Println("Please set OPENAI_API_KEY environment variable")
		os.Exit(1)
	}

	rootCmd.Flags().StringP("prompt", "p", "", "prompt for chatgpt")
	if err := rootCmd.MarkFlagRequired("prompt"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
