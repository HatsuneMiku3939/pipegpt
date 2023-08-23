package generic

import (
	"github.com/HatsuneMiku3939/pipegpt/pkg/chatgpt"
)

// New creates a new generic question app
func New(client *chatgpt.Client) *App {
	return &App{
		client: client,
	}
}

// App is the default app
type App struct {
	client *chatgpt.Client
}

// Run runs the app
func (a *App) Run(role string, prompt string, input string) (string, error) {
	return a.client.Question(role, prompt, input)
}
