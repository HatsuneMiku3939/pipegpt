package chatgpt

import (
	"context"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

const (
	chatTimeout = 30 * time.Second
)

// Chat question to chatgpt with given prompt and user input
func Chat(apikey string, prompt string, input string) (string, error) {
	// join prompt and input
	content := prompt + "\n" + input

	// create client
	client := openai.NewClient(apikey)

	// create chat completion
	ctx, cancel := context.WithTimeout(context.Background(), chatTimeout)
	defer cancel()

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
