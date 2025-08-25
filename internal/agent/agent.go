package agent

import (
	"context"
	"gogurt/internal/state"
)

// Agent interface
type Agent interface {
    Init(ctx context.Context, params map[string]interface{}) error // Flexible initialization
    Invoke(ctx context.Context, input string) (*AgentCallResult, error) // Synchronous call with chaining/metadata
    InvokeAsync(ctx context.Context, input string) (<-chan AgentCallResult, error) // Async/streaming call
    State() state.AgentState
    Planner() Planner
    Delegate(ctx context.Context, input string, agents []Agent) ([]*AgentCallResult, error) // Multi-agent orchestration
}