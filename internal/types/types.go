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

// represents a chunk of text from a source.
type Document struct {
	PageContent string
	Metadata    map[string]any
}

// interface that document loaders must satisfy
type DocumentLoader interface {
	LoadDocuments() ([]Document, error)
}

type AgentConfig struct {
	Name     string
	AiClient map[string]any
	Tools    []string
	Params   map[string]any
	Children []AgentConfig
}

// Used for agent introspection
type AgentDescription struct {
	Name         string
	AiClient     map[string]any
	Capabilities []string
	Tools        []string
	Children     []*AgentDescription
}

type AgentCallResult struct {
	Output   string
	Error    error
	Metadata map[string]interface{}
	Next     *AgentCallResult
}

type PipelineStep func(context.Context, any) (any, error)

type NextStep func(context.Context, any) (*AgentCallResult, error)

type EndStep func(context.Context, any) (*AgentCallResult, error)

