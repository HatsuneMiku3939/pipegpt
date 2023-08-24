> This README is generated from pipegpt itself. modified by hand.
>
> cat cmd/pipegpt/main.go cmd/pipegpt/cmd/subcmd.go cmd/pipegpt/cmd/root.go | pipegpt -p "Give me a detailed and structured README.md for these code. README.md must contain description within 300 words with 2~3 paragraphs. README.md must contain installation steps using 'go install'. README.md must contain 2~3 simple usage example. README.md must contain 2~3 advanced usage example using subcommands. README.md must contain detailed description of config file and env vars. README.md must contain advanced usage example using config file. README.md must contain license information (MIT). "

# pipegpt

## Description
pipegpt is a Robust Command Line Interface (CLI) tool that you can leverage to interact with the openAI Chatgpt model. Simple yet potent, this tool, written in Go, waits for your input from stdin and outputs the result on stdout. The indicated prompt serves as the first piece of input for chatgpt.

The CLI tool's prime function includes creating subcommands based on configurations and executing the root command. It carries the potency to produce two types of subcommands - 'generic' and 'function-call' and organise them under the main command. The application is encased in a sophisticated error handling mechanism helping you detect and rectify issues aptly.

## Installation
To install pipegpt, follow these steps:

Use the 'go install' command to install pipegpt:

```
$ go install github.com/HatsuneMiku3939/pipegpt/cmd/pipegpt@latest
```

## Simple Usage Examples
Once you've installed pipegpt, you can interact with it from the command line like so:

1. For staged code review:

```
$ git diff --staged | pipegpt -p "code review for this change"
```

2. For converting JSON to YAML:

```
$ cat sample.json | pipegpt -p "convert JSON to YAML"
```

## Advanced Usage Examples

1. For defining a custom role and a prompt:

```
$ echo "What is Linux?"|  pipegpt --role "act like you're linus torvalds." --prompt "answer the question."
```

2. For creating a subcommand:

You can create subcommands by defining roles and prompts in a configuration file. For example, you can create a subcommand for code review like so:

```
review:
  role: Act like you're professional IT engineer.
  prompt: code review for this change
```

Finally, you can run the subcommand like so:

```
$ git diff --staged | pipegpt review
```

3. For creating a function-call subcommand:

You can create a function-call subcommand by defining a function in a configuration file. For example, you can create a subcommand for generating a bash command like so:

```
shell:
  role: Act like you're professional IT engineer.
  prompt: write a bash command for the following task.
  function-call:
    - name: command
      description: bash command to execute
      parameters:
        type: object
        properties:
          command:
            type: string
            description: bash command to execute
        required:
          - command
```

Define wrapper bash function like so:

```
function shell() {
  input="$1"

  cmd=$(echo "$input" | pipegpt shell)
  if [ $? -ne 0 ]; then
    return
  fi

  cmd=$(echo "$cmd" | jq -r '.command')
  if [ $? -ne 0 ]; then
    return
  fi

  printf "Run '%s' ? [Y/n] " "$cmd"
  read -r yn

  if [ "$yn" == "y" ] || [ "$yn" == "Y" ] || ["$yn" == ""]; then
    eval "$cmd"
  fi
}
```

Finally, you can run the subcommand like so:

```
$ shell "find all go files in the current directory"

Run 'find . -name "*.go"' ? [Y/n] y
./app/generic/generic.go
./app/function/function.go
./tools/tools.go
./pkg/chatgpt/client.go
./pkg/chatgpt/chatgpt.go
./pkg/stdin/stdin.go
./cmd/pipegpt/main.go
./cmd/pipegpt/cmd/subcmd.go
./cmd/pipegpt/cmd/root.go
```

## Config Files and Environment Variables

Config file can be defined using the `--config` option. If no file is specified, the tool defaults to reading `$HOME/.pipegpt.yaml` or `./.pipegpt.yaml`.

The environment variables are prefixed with 'pipegpt'. Here are some of them:

- `PIPEGPT_API_KEY`: The OpenAI API key
- `PIPEGPT_API_MODEL`: The OpenAI API model used
- `PIPEGPT_API_TIMEOUT`: The timeout value for the OpenAI API request
- `PIPEGPT_DEFAULT_ROLE`: The default role of the AI assistant

If you create a subcommand, you can override the default role by defining a role in the configuration file. For example:

```
review:
  role: Act like you're professional IT engineer.
  prompt: code review for this change
```

Note that, role and prompt of your subcommand also can be configured using environment variables. For example:

```
$ git diff --staged | PIPEGPT_REVIEW_ROLE="Act like you're professional IT engineer." PIPEGPT_REVIEW_PROMPT="code review for this change" pipegpt review
```

Detailed description of config file and env vars can be found from help message. (including your subcommands)

```
$ pipegpt --help
```

You can find a sample config file [here](./.pipegpt.sample.yaml).

## License Information

pipegpt is licensed under the MIT License. The terms of the license can be found in the [LICENSE](LICENSE) file.
