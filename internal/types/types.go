package types

import (
	"context"
	"gogurt/internal/state"
	"time"
)

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

type StateMessageMeta struct {
	Current *StateMessage
	Next *StateMessage
	Previous *StateMessage
	CurrentState *state.AgentState
}

type StateMessage struct {
	Id      string
	Sender  Role
	Message string
	Timestamp time.Time
	Meta *StateMessageMeta
}

type ToolCall struct {
	Name string
	Args map[string]any
}



type PipelineStep func(context.Context, any) (any, error)

type NextStep func(context.Context, any) (*AgentCallResult, error)

type EndStep func(context.Context, any) (*AgentCallResult, error)


// Logging
type LogLevel int

const (
	DEBUG LogLevel = -1
	INFO LogLevel = iota
	WARNING
	ERROR
	FATAL
)

func (l LogLevel) String() string {
	switch l {
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

type LogFormat string

const (
	FormatText LogFormat = "text"
	FormatJSON LogFormat = "json"
)

type TimeStamp struct {
	time.Time
}