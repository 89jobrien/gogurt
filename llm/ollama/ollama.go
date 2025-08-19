package ollama

import (
	"context"
	"fmt"

	"gogurt/config"
	"gogurt/types"

	"github.com/ollama/ollama/api"
)

type Ollama struct {
	client *api.Client
	model  string
}

// Update New to accept the config
func New(cfg *config.Config) (types.LLM, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to create ollama client: %w", err)
	}

	return &Ollama{
		client: client,
		model:  cfg.OllamaModel,
	}, nil
}

// Generate generates a response from the Ollama API.
func (o *Ollama) Generate(ctx context.Context, messages []types.ChatMessage) (*types.ChatMessage, error) {
	apiMessages := make([]api.Message, len(messages))
	for i, msg := range messages {
		apiMessages[i] = api.Message{
			Role:    string(msg.Role),
			Content: msg.Content,
		}
	}

	req := &api.ChatRequest{
		Model:    o.model,
		Messages: apiMessages,
	}

	var responseMessage types.ChatMessage
	err := o.client.Chat(ctx, req, func(res api.ChatResponse) error {
		responseMessage = types.ChatMessage{
			Role:    types.Role(res.Message.Role),
			Content: res.Message.Content,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &responseMessage, nil
}