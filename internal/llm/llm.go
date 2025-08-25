package llm

import (
	"context"

	"gogurt/internal/types"
)

type LLM interface {
	Generate(ctx context.Context, messages []types.ChatMessage) (*types.ChatMessage, error)
	Stream(ctx context.Context, messages []types.ChatMessage, onToken func(token string) error) (*types.ChatMessage, error)
	HealthCheck(ctx context.Context) error
	Metadata() map[string]any
}
