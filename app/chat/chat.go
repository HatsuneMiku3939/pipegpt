package chat

import (
	"os"

	"github.com/HatsuneMiku3939/pipegpt/pkg/chatgpt"
	"github.com/HatsuneMiku3939/pipegpt/pkg/in"
	"github.com/HatsuneMiku3939/pipegpt/pkg/out"

	"github.com/sashabaranov/go-openai"
)

// New creates a new chat app
func New(client *chatgpt.Client, output *out.Out, history []openai.ChatCompletionMessage) *Chat {
	return &Chat{
		client:  client,
		output:  output,
		History: history,
	}
}

// App is the chat app
type Chat struct {
	client  *chatgpt.Client
	output  *out.Out
	History []openai.ChatCompletionMessage
}

// NewChatHistory creates a new chat history from role, prompt, user input and first response
func CreateChatHistory(role string, prompt string, input string, response interface{}) ([]openai.ChatCompletionMessage, error) {
	return chatgpt.CreateChatHistory(role, prompt, input, response)
}

// Run runs the chat app
func (c *Chat) Run() {
	// display the history
	for _, message := range c.History {
		c.output.Emit(message.Content)
	}

	// it is possible that os.Stdin is pipe,
	// so, we need to use /dev/tty to read user input instead of os.Stdin
	// TODO: support Windows
	tty, err := os.OpenFile("/dev/tty", os.O_RDONLY, os.ModeDevice)
	if err != nil {
		c.output.Emit(err.Error())
		return
	}
	defer tty.Close()
	input := in.New(tty)

	// read user input until met a "/quit" command
	c.output.Emit("> \"/quit\" to quit\n")
	input.ConsumeUntil("> ", "/quit", func(content string, err error) {
		if err != nil {
			c.output.Emit(err.Error())
			return
		}

		// send the user input to the API
		response, err := c.client.Chat(c.History, content)
		if err != nil {
			c.output.Emit(err.Error())
			return
		}

		// display the response
		c.output.Emit(response)

		// Update the history
		c.History = append(c.History, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: content,
		})
		c.History = append(c.History, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: response,
		})
	})
}
