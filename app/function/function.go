package function

import (
	"encoding/json"

	"github.com/HatsuneMiku3939/pipegpt/pkg/chatgpt"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

// New creates a new function call app
func New(client *chatgpt.Client) *App {
	return &App{
		client: client,
	}
}

// App is the genefic app
type App struct {
	client *chatgpt.Client
}

// Run runs the app
func (a *App) Run(role string, prompt string, input string, funcs []openai.FunctionDefinition) (map[string]interface{}, error) {
	res, err := a.client.FunctionCall(role, prompt, input, funcs)
	if err != nil {
		return map[string]interface{}{}, err
	}

	return res, nil
}

// ToFunctionSchema converts a yaml definition to a JSONSchema for function call
func ToFunctionSchema(v interface{}) (openai.FunctionDefinition, error) {
	// serialize the yaml definition to json
	rawJSON, err := json.Marshal(v)
	if err != nil {
		return openai.FunctionDefinition{}, err
	}

	// deserialize the json to a FunctionDefinition struct
	var definition definition
	if err := json.Unmarshal(rawJSON, &definition); err != nil {
		return openai.FunctionDefinition{}, err
	}

	// convert to JSONSchema(openai.FunctionDefinition)
	functionDefinition := openai.FunctionDefinition{
		Name:        definition.Name,
		Description: definition.Description,
		Parameters:  definition.Parameters,
	}
	return functionDefinition, nil
}

// definition is a definition of a function call
type definition struct {
	Name        string                `json:"name"`
	Description string                `json:"description,omitempty"`
	Parameters  jsonschema.Definition `json:"parameters"`
}
