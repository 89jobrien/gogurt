package agent

import (
	"context"
	"gogurt/internal/types"
)

type ResearchAgent struct {
	// internal state and tools
}

func (a *ResearchAgent) Init(ctx context.Context, config types.AgentConfig) error {
	_, _ = NewAgentGraph(ctx, config, nil)
	return nil
}

func (a *ResearchAgent) Invoke(ctx context.Context, input any) (any, error) {
	// Run web search, etc
	return "web search result", nil
}
func (a *ResearchAgent) InvokeAsync(ctx context.Context, input any) (<-chan any, <-chan error) {
	// Optional
	return nil, nil
}
func (a *ResearchAgent) Delegate(ctx context.Context, task any) (any, error) {
	// Could delegate to subordinate agents
	return nil, nil
}
func (a *ResearchAgent) Planner() Planner       { return nil }
func (a *ResearchAgent) State() any             { return nil }
func (a *ResearchAgent) Capabilities() []string { return []string{"websearch"} }
func (a *ResearchAgent) Describe() *types.AgentDescription {
	return &types.AgentDescription{
		Name:         "ResearchAgent",
		Capabilities: a.Capabilities(),
		Tools:        []string{"WebSearchTool"},
	}
}

// Register the agent automatically via init()
func init() {
	RegisterAgent("ResearchAgent", func() Agent { return &ResearchAgent{} })
}
