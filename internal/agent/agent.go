package agent

import (
	"context"
	"fmt"
	"gogurt/internal/types"
)

// Agent interface with inspection and description
type Agent interface {
	Init(ctx context.Context, config types.AgentConfig) error
	Invoke(ctx context.Context, input any) (any, error)
	InvokeAsync(ctx context.Context, input any) (<-chan any, <-chan error)
	Delegate(ctx context.Context, task any) (any, error)
	Planner() Planner
	State() any
	Capabilities() []string
	Describe() *types.AgentDescription
}

// Describes an agent programmatically
type AgentRegistry map[string]func() Agent

var RegisteredAgents = make(AgentRegistry)

func RegisterAgent(name string, factory func() Agent) {
	RegisteredAgents[name] = factory
}

func NewAgent(config types.AgentConfig) (Agent, error) {
	factory, ok := RegisteredAgents[config.Name]
	if !ok {
		return nil, fmt.Errorf("agent not registered: %s", config.Name)
	}
	return factory(), nil
}
