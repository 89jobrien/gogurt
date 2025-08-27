package agent

import (
	"context"
	"gogurt/internal/logger"
	"gogurt/internal/state"
	"gogurt/internal/types"
	"time"
)

type ResearchAgent struct {
	state state.AgentState
}

func (a *ResearchAgent) Init(ctx context.Context, config types.AgentConfig) error {
	// The NewAgentGraph call seems to be for initializing child agents.
	// We'll keep it, assuming it's part of a larger design, but note that
	// the result isn't used in this agent's current form.
	_, err := NewAgentGraph(ctx, config, RegisteredAgents)
	if err != nil {
		return err
	}
	return nil
}

// Invoke is now fully asynchronous, returning channels for results and errors.
func (a *ResearchAgent) Invoke(ctx context.Context, input any) (<-chan any, <-chan error) {
	resultCh := make(chan any, 1)
	errorCh := make(chan error, 1)

	go func() {
		defer close(resultCh)
		defer close(errorCh)

		// Placeholder for actual web search logic
		// In a real implementation, you would call a search tool here.
		logger.InfoCtx(ctx, "ResearchAgent executing web search for: %v", input)
		time.Sleep(100 * time.Millisecond) // Simulate network latency

		select {
		case resultCh <- "web search result":
		case <-ctx.Done():
			errorCh <- ctx.Err()
		}
	}()

	return resultCh, errorCh
}

// State returns the agent's current state.
func (a *ResearchAgent) State() *state.AgentState {
	return &a.state
}

// Describe returns a description of the agent.
func (a *ResearchAgent) Describe() *types.AgentDescription {
	return &types.AgentDescription{
		Name:         "ResearchAgent",
		Capabilities: []string{"web-search", "data-retrieval"},
		Tools:        []string{"WebSearchTool"}, // Example tool
	}
}

// OnMessage handles agent-to-agent communication asynchronously.
func (a *ResearchAgent) OnMessage(ctx context.Context, msg *types.StateMessage) (<-chan *types.StateMessage, <-chan error) {
	resultCh := make(chan *types.StateMessage, 1)
	errorCh := make(chan error, 1)

	go func() {
		defer close(resultCh)
		defer close(errorCh)

		logger.InfoCtx(ctx, "ResearchAgent received message from %s: %s", msg.Sender, msg.Message)
		// This agent's response is a simulated search result based on the incoming message.
		responseMsg := NewStateMessage(types.Role("research_agent"), "web search result for: "+msg.Message)

		select {
		case resultCh <- responseMsg:
		case <-ctx.Done():
			errorCh <- ctx.Err()
		}
	}()

	return resultCh, errorCh
}

func init() {
	RegisterAgent("ResearchAgent", func() Agent {
		logger.Info("Initializing ResearchAgent")
		return &ResearchAgent{state: state.NewMemoryState()}
	})
}