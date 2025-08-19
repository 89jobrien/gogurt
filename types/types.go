package types

import "context"

// Role defines the role of the message sender.
type Role string

const (
	RoleUser      Role = "user"
	RoleSystem    Role = "system"
	RoleAssistant Role = "assistant"
)

// ChatMessage represents a single message in a conversation.
type ChatMessage struct {
	Role    Role
	Content string
}

// LLM is the interface that language model implementations must satisfy.
type LLM interface {
	Generate(ctx context.Context, messages []ChatMessage) (*ChatMessage, error)
}