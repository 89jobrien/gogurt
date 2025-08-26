package agent

import (
	"context"
	"gogurt/internal/state"
	"gogurt/internal/types"
)

type ResearchAgent struct {
	state state.AgentState
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
func (a *ResearchAgent) State() *state.AgentState {
	return &a.state
}
func (a *ResearchAgent) Capabilities() []string { return []string{"websearch"} }
func (a *ResearchAgent) Describe() *types.AgentDescription {
	return &types.AgentDescription{
		Name:         "ResearchAgent",
		Capabilities: a.Capabilities(),
		Tools:        []string{"WebSearchTool"},
	}
}

func (a *ResearchAgent) OnMessage(ctx context.Context, msg *types.StateMessage) (*types.StateMessage, error) {
	return &types.StateMessage{
		Id:        "msg-websearch",
		Sender:    types.Role("websearch"),
		Message:   "web search result",
		Timestamp: types.TimeStamp.Local(types.TimeStamp{}),
		Meta:      nil,
	}, nil
}

func (a *ResearchAgent) OnMessageAsync(ctx context.Context, msg *types.StateMessage) (<-chan *types.StateMessage, <-chan error) {
	msgCh := make(chan *types.StateMessage, 1)
	errCh := make(chan error, 1)
	go func() {
		defer close(msgCh)
		defer close(errCh)
		m, err := a.OnMessage(ctx, msg)
		if err != nil {
			errCh <- err
		} else {
			msgCh <- m
		}
	}()
	return msgCh, errCh
}

func init() {
	RegisterAgent("ResearchAgent", func() Agent { 
		return &ResearchAgent{state: state.NewMemoryState()} 
	})
}