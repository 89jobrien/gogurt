package agent

import (
	"context"
	"fmt"
	"gogurt/internal/state"
	"gogurt/internal/types"
	"time"
)

// Utility: generate a new message with all fields filled
func NewStateMessage(sender types.Role, message string) *types.StateMessage {
	return &types.StateMessage{
		Id:        fmt.Sprintf("msg-%d", time.Now().UnixNano()),
		Sender:    sender,
		Message:   message,
		Timestamp: time.Now(),
		Meta:      &types.StateMessageMeta{},
	}
}

// Agent interface supporting both human/system and agent-to-agent communication.
type Agent interface {
	// Human/system interaction (CLI, API, etc.)
	Invoke(ctx context.Context, input any) (<-chan any, <-chan error)

	// Agent-to-agent communication (robust workflow)
	OnMessage(ctx context.Context, msg *types.StateMessage) (<-chan *types.StateMessage, <-chan error)

	// Returns the current agent state (may be nil/unimplemented)
	State() *state.AgentState

	// Describe agent programmatically (metadata, summary, etc.)
	Describe() *types.AgentDescription

	// Optionally: initialization/config
	Init(ctx context.Context, config types.AgentConfig) error
}

// Registry for agent factories
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
	if factory == nil {
		return nil, fmt.Errorf("agent factory for %q is nil", config.Name)
	}
	return factory(), nil
}

// MultiAgentCoordinator orchestrates communication and workflows between multiple agents using StateMessage.
type MultiAgentCoordinator struct {
	Agents   map[string]Agent
	Workflow []string
}

// SendMessageTo sends a message to an agent and links meta information.
func (mac *MultiAgentCoordinator) SendMessageTo(ctx context.Context, agentID string, input *types.StateMessage) (<-chan *types.StateMessage, <-chan error) {
	outCh := make(chan *types.StateMessage, 1)
	errCh := make(chan error, 1)

	go func() {
		defer close(outCh)
		defer close(errCh)

		agent, ok := mac.Agents[agentID]
		if !ok {
			errCh <- &AgentError{"agent not found: " + agentID}
			return
		}

		outputCh, errChFromAgent := agent.OnMessage(ctx, input)
		select {
		case output := <-outputCh:
			// Thread history in meta
			if output.Meta == nil {
				output.Meta = &types.StateMessageMeta{}
			}
			output.Meta.Previous = input
			output.Meta.Current = output
			output.Meta.CurrentState = agent.State()
			outCh <- output
		case err := <-errChFromAgent:
			errCh <- err
		case <-ctx.Done():
			errCh <- ctx.Err()
		}
	}()

	return outCh, errCh
}

// BroadcastMessage sends a message through each agent in Workflow sequentially, threading StateMessageMeta.
func (mac *MultiAgentCoordinator) BroadcastMessage(ctx context.Context, initial *types.StateMessage) (<-chan *types.StateMessage, <-chan error) {
	outCh := make(chan *types.StateMessage)
	errCh := make(chan error, 1)

	go func() {
		defer close(outCh)
		defer close(errCh)

		var finalMessages []*types.StateMessage
		curr := initial

		for _, agentID := range mac.Workflow {
			respCh, broadcastErrCh := mac.SendMessageTo(ctx, agentID, curr)

			select {
			case resp := <-respCh:
				if curr.Meta != nil {
					curr.Meta.Next = resp // Link next in the chain
				}
				finalMessages = append(finalMessages, resp)
				curr = resp
			case err := <-broadcastErrCh:
				errCh <- err
				return
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			}
		}

		for _, msg := range finalMessages {
			outCh <- msg
		}
	}()

	return outCh, errCh
}

// AgentError for errors originating in agent communication.
type AgentError struct{ msg string }

func (e *AgentError) Error() string { return e.msg }