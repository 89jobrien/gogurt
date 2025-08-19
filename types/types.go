package types

import "context"

type Role string

const (
	RoleUser      Role = "user"
	RoleSystem    Role = "system"
	RoleAssistant Role = "assistant"
)

type ChatMessage struct {
	Role    Role
	Content string
}

type LLM interface {
	Generate(ctx context.Context, messages []ChatMessage) (*ChatMessage, error)
}