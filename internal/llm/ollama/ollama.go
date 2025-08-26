package ollama

import (
	"context"
	"fmt"
	"strings"

	"gogurt/internal/config"
	"gogurt/internal/llm"
	"gogurt/internal/types"

	"github.com/ollama/ollama/api"
)

type Ollama struct {
	client *api.Client
	model  string
}

// HealthCheck implements types.LLM.
func (o *Ollama) HealthCheck(ctx context.Context) error {
	_, err := o.Generate(ctx, []types.ChatMessage{{Role: "system", Content: "ping"}})
	return err
}

// Metadata implements types.LLM.
func (o *Ollama) Metadata() map[string]any {
	md := make(map[string]any)
	md["model"] = o.model
	return md
}

func New(cfg *config.Config) (llm.LLM, error) {
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

// AGenerate provides an asynchronous Generate.
func (o *Ollama) AGenerate(ctx context.Context, messages []types.ChatMessage) (<-chan *types.ChatMessage, <-chan error) {
	msgCh := make(chan *types.ChatMessage, 1)
	errCh := make(chan error, 1)
	go func() {
		defer close(msgCh)
		defer close(errCh)
		msg, err := o.Generate(ctx, messages)
		if err != nil {
			errCh <- err
			return
		}
		msgCh <- msg
	}()
	return msgCh, errCh
}

// Stream streams model output from Ollama, forwarding chunks to onToken.
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
		if responseRole == "" && res.Message.Role != "" {
			responseRole = types.Role(res.Message.Role)
		}

		token := res.Message.Content
		if token != "" {
			if err := onToken(token); err != nil {
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

// AStream provides an asynchronous streaming interface returning tokens.
func (o *Ollama) AStream(ctx context.Context, messages []types.ChatMessage) (<-chan string, <-chan error) {
	tokenCh := make(chan string)
	errCh := make(chan error, 1)
	go func() {
		defer close(tokenCh)
		defer close(errCh)

		_, err := o.Stream(ctx, messages, func(token string) error {
			tokenCh <- token
			return nil
		})
		if err != nil {
			errCh <- err
		}
	}()
	return tokenCh, errCh
}