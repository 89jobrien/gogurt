package ollama

import (
	"context"
	"fmt"
	// "net/url"
	"os"

	// "gogurt/agent"
	"gogurt/types"

	"github.com/ollama/ollama/api"
)

type Ollama struct {
	client *api.Client
	model  string
}

func New() (types.LLM, error) {
	host := os.Getenv("OLLAMA_HOST")
	if host == "" {
		host = "http://localhost:11434"
	}

	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = "llama3"
	}

	// ollamaURL, err := url.Parse(host)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to parse ollama host: %w", err)
	// }

	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to create ollama client: %w", err)
	}

	return &Ollama{
		client: client,
		model:  model,
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