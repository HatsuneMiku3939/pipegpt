package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/HatsuneMiku3939/pipegpt/app/chat"
	"github.com/HatsuneMiku3939/pipegpt/app/generic"
	"github.com/HatsuneMiku3939/pipegpt/pkg/chatgpt"
	"github.com/HatsuneMiku3939/pipegpt/pkg/in"
	"github.com/HatsuneMiku3939/pipegpt/pkg/out"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// defaultConfigFileName is used when config file is not specified
const defaultConfigFileName = ".pipegpt"

// defaultRole is used when role is not specified in config file
const defaultRole = "Act like you are professional IT engineer to help solve user's business problem in enterprise IT tech company."

// configFile config file given by flag
var configFile string

var RootCmd = &cobra.Command{
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

		chatEnabled, err := cmd.Flags().GetBool("chat")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		role := viper.GetString("default.role")
		i := in.New(os.Stdin)
		input, err := i.Consume()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

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

		// print result with markdown formatter
		o := out.New(os.Stdout, out.MarkdownFormatter)
		if !chatEnabled {
			o.Emit(result)
			return
		}

		// interactive chat mode
		history, err := chat.CreateChatHistory(role, prompt, input, result)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		chat.New(client, o, history).Run()
	},
}

// initFlag is function to initialize flag of RootCmd
func initFlag() {
	// define config flag of rootCmd
	RootCmd.PersistentFlags().StringVar(&configFile, "config", "", fmt.Sprintf("config file (default is $HOME/%s.yaml or ./%s.yaml", defaultConfigFileName, defaultConfigFileName))

	// define prompt flag of rootCmd
	RootCmd.PersistentFlags().StringP("key", "k", "", "OpenAI API key, you can also set it with PIPEGPT_API_KEY environment variable or config file")
	RootCmd.PersistentFlags().StringP("model", "m", "gpt-4", "OpenAI API model, you can also set it with PIPEGPT_API_MODEL environment variable or config file")
	RootCmd.PersistentFlags().StringP("timeout", "t", "240s", "Timeout of OpenAI API request, you can also set it with PIPEGPT_API_TIMEOUT environment variable or config file")
	RootCmd.PersistentFlags().StringP("endpoint", "e", "", "Endpoint of Azure OpenAI API, you can also set it with PIPEGPT_API_ENDPOINT environment variable or config file")
	RootCmd.PersistentFlags().String("conversion", "", "comma separated list of model conversion table of Azure OpenAI API. ex) 'gpt-4=foo-gpt-4, gpt-3=bar-gpt-3'")
	RootCmd.PersistentFlags().BoolP("chat", "c", false, "interactive chat mode enable (default is false)")
	RootCmd.Flags().StringP("role", "r", defaultRole, "role of the AI assistant, you can also set it with PIPEGPT_DEFAULT_ROLE environment variable or config file")
	RootCmd.Flags().StringP("prompt", "p", "", "prompt to use for the AI assistant")
	if err := RootCmd.MarkFlagRequired("prompt"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// bind flag to viper
	if err := viper.BindPFlag("api.key", RootCmd.PersistentFlags().Lookup("key")); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := viper.BindPFlag("api.model", RootCmd.PersistentFlags().Lookup("model")); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := viper.BindPFlag("api.timeout", RootCmd.PersistentFlags().Lookup("timeout")); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := viper.BindPFlag("api.endpoint", RootCmd.PersistentFlags().Lookup("endpoint")); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := viper.BindPFlag("api.conversion", RootCmd.PersistentFlags().Lookup("conversion")); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := viper.BindPFlag("default.role", RootCmd.Flags().Lookup("role")); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// initViper is function to initialize viper
func initViper() {
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		// get home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(defaultConfigFileName)
	}

	// read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("pipegpt")
	viper.AutomaticEnv()

	// If a config file is found
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("can't read config:", err)
		os.Exit(1)
	}
}

func init() {
	initFlag()
	initViper()
}

// createClient is function to create chatgpt client
func createClient() (*chatgpt.Client, error) {
	// if endpoint is set, create azure openai client
	if viper.GetString("api.endpoint") != "" {
		return createAzureOpenAIClient()
	}

	// otherwise, create openai client
	return createOpenAIClient()
}

// createOpenAIClient is function to create openai client
func createOpenAIClient() (*chatgpt.Client, error) {
	// get arguments from viper
	key := viper.GetString("api.key")
	timeout, err := time.ParseDuration(viper.GetString("api.timeout"))
	if err != nil {
		return nil, err
	}
	model := viper.GetString("api.model")

	// create client
	client := chatgpt.NewClient(key, model, timeout)
	return client, nil
}

// createAzureOpenAIClient is function to create azure openai client
func createAzureOpenAIClient() (*chatgpt.Client, error) {
	// get arguments from viper
	key := viper.GetString("api.key")
	endpoint := viper.GetString("api.endpoint")
	timeout, err := time.ParseDuration(viper.GetString("api.timeout"))
	if err != nil {
		return nil, err
	}
	model := viper.GetString("api.model")
	rawModelMap := viper.GetString("api.conversion")

	// create model map
	const requiredConversionTokenCount = 2
	modelMap := map[string]string{}
	for _, v := range strings.Split(rawModelMap, ",") {
		v = strings.TrimSpace(v)
		kv := strings.Split(v, "=")
		if len(kv) != requiredConversionTokenCount {
			return nil, fmt.Errorf("'api.conversion' must be a key-value pair separated by '='")
		}
		modelMap[kv[0]] = kv[1]
	}

	// create client
	client := chatgpt.NewAzureOpenAIClient(key, endpoint, model, modelMap, timeout)
	return client, nil
}
