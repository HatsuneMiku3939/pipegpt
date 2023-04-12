> This README is generated from pipegpt itself. modified by hand.
>
> cat cmd/pipegpt/cmd/root.go | pipegpt -p "give me a README.md" > README.md

# pipegpt

pipegpt is a simple command-line tool for interacting with OpenAI GPT-3 using chatgpt.

## Installation

You can install pipegpt by running `go install github.com/HatsuneMiku3939/pipegpt/cmd/pipegpt@latest`.

## Usage

To use pipegpt, you must set the `OPENAI_API_KEY` environment variable to your OpenAI API key. You can get an API key from the [OpenAI website](https://beta.openai.com/docs/).

You can use pipegpt to question chatgpt with your pre-defined prompt or interactively (in cli manner). pipegpt always waits for your input from stdin and outputs the result to stdout. The provided prompt is used as the first input to chatgpt.

Here is a simple example:

```
# code review in staged changes
git diff --staged | pipegpt -p "code review for this change"

# commit with suggested commit message
git diff --staged | pipegpt -p "suggest conventional commit messages for these changes" | git commit -F - -e

# convert JSON to YAML
cat sample.json | pipegpt -p "convert JSON to YAML"
```

## Flags

- `-p` or `--prompt`: set the prompt for chatgpt. This flag is required.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

