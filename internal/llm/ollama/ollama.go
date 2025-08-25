package ollama

import (
	"context"
	"fmt"
	"strings"

	"gogurt/internal/config"
	"gogurt/internal/types"

	"github.com/ollama/ollama/api"
)

type Ollama struct {
	client *api.Client
	model  string
}

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

	var responseContent strings.Builder
	var responseRole types.Role

	err := o.client.Chat(ctx, req, func(res api.ChatResponse) error {
		responseContent.WriteString(res.Message.Content)

		if responseRole == "" {
			responseRole = types.Role(res.Message.Role)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if responseRole == "" {
		responseRole = types.RoleAssistant
	}

	return &types.ChatMessage{
		Role:    responseRole,
		Content: responseContent.String(),
	}, nil
}

// Stream streams model output from Ollama, forwarding chunks to onToken
func (o *Ollama) Stream(ctx context.Context, messages []types.ChatMessage, onToken func(token string) error) (*types.ChatMessage, error) {
	apiMessages := make([]api.Message, len(messages))
	for i, msg := range messages {
		apiMessages[i] = api.Message{
			Role:    string(msg.Role),
			Content: msg.Content,
		}
	}

	stream := true
	req := &api.ChatRequest{
		Model:    o.model,
		Messages: apiMessages,
		Stream:   &stream,
	}

	var responseContent strings.Builder
	var responseRole types.Role

	err := o.client.Chat(ctx, req, func(res api.ChatResponse) error {
		// Determine role if provided
		if responseRole == "" && res.Message.Role != "" {
			responseRole = types.Role(res.Message.Role)
		}

		token := res.Message.Content
		if token != "" {
			// Forward to the provided callback.
			if err := onToken(token); err != nil {
				// Returning an error from this callback will stop the stream and bubble up.
				return err
			}
			responseContent.WriteString(token)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if responseRole == "" {
		responseRole = types.RoleAssistant
	}

	return &types.ChatMessage{
		Role:    responseRole,
		Content: responseContent.String(),
	}, nil
}