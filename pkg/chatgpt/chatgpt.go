package chatgpt

import (
	"context"
	"encoding/json"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

// Chat question to chatgpt with given prompt and user input
func (gpt *Client) Question(role string, prompt string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gpt.timeout)
	defer cancel()

	// create chat completion
	resp, err := gpt.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: gpt.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: role,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("%s\n---\n%s", prompt, input),
				},
			},
		},
	)
	if err != nil {
		return "", err
	}

	// return first choice
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	return resp.Choices[0].Message.Content, nil
}

// FunctionCall question to OpenAI in function calling format with given prompt and user input, and function definitions
func (gpt *Client) FunctionCall(role string, prompt string, input string, funcs []openai.FunctionDefinition) (map[string]interface{}, error) {
	// create chat completion
	ctx, cancel := context.WithTimeout(context.Background(), gpt.timeout)
	defer cancel()

	// create chat completion
	resp, err := gpt.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: gpt.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: role,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("%s\n---\n%s", prompt, input),
				},
			},
			Functions: funcs,
		},
	)
	if err != nil {
		return nil, err
	}

	// return first choice if function call arguments are valid
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned")
	}

	if resp.Choices[0].Message.FunctionCall == nil {
		return nil, fmt.Errorf("no function call returned")
	}

	var args map[string]interface{}
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.FunctionCall.Arguments), &args); err != nil {
		return nil, err
	}

	return args, nil
}
