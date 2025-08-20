package types

import "context"

// defines the role of the message sender
type Role string

const (
	RoleUser      Role = "user"
	RoleSystem    Role = "system"
	RoleAssistant Role = "assistant"
)

// represents a single message in a conversation
type ChatMessage struct {
	Role    Role
	Content string
}

// interface that language model implementations must satisfy
type LLM interface {
	Generate(ctx context.Context, messages []ChatMessage) (*ChatMessage, error)
}

// represents a chunk of text from a source.
type Document struct {
	PageContent string
	Metadata    map[string]interface{}
}
