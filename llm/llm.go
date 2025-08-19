package llm

import (
	"context"
	"gogurt/types"
)

// Example implementation
type MyLLM struct{}

func (m *MyLLM) Generate(ctx context.Context, messages []types.ChatMessage) (*types.ChatMessage, error) {
	return &types.ChatMessage{Role: types.RoleAssistant, Content: "response"}, nil
}