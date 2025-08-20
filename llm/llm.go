package llm

import (
	"context"

	"gogurt/types"
)

type LLM interface {
	Generate(ctx context.Context, messages []types.ChatMessage) (*types.ChatMessage, error)
	Stream(ctx context.Context, messages []types.ChatMessage, onToken func(token string) error) (*types.ChatMessage, error)
}