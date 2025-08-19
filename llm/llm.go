package llm

import (
	"context"
	"gogurt/types"
)

// Example implementation that satisfies agent.LLM
type MyLLM struct{}

func (m *MyLLM) Generate(ctx context.Context, messages []types.ChatMessage) (*types.ChatMessage, error) {
	// real implementation here
	return &types.ChatMessage{Role: types.RoleAssistant, Content: "response"}, nil
}