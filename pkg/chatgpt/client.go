package chatgpt

import (
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// Client is a client for ChatGPT client
type Client struct {
	client  *openai.Client
	timeout time.Duration
	model   string
}

// NewClient creates a new GPTClient
func NewClient(apiKey string, model string, timeout time.Duration) *Client {
	// create client
	client := openai.NewClient(apiKey)

	return &Client{
		client:  client,
		timeout: timeout,
		model:   model,
	}
}

// NewAzureOpenAIClient creates a new GPTClient
func NewAzureOpenAIClient(apiKey string, endpoint string, model string, modelMapping map[string]string, timeout time.Duration) *Client {
	// create client
	config := openai.DefaultAzureConfig(apiKey, endpoint)
	config.APIVersion = "2023-07-01-preview"
	config.AzureModelMapperFunc = func(model string) string {
		if val, ok := modelMapping[model]; ok {
			return val
		}
		return model
	}

	client := openai.NewClientWithConfig(config)
	return &Client{
		client:  client,
		timeout: timeout,
		model:   model,
	}
}
