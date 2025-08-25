package agent

import (
	"context"
	"fmt"
	"gogurt/internal/state"
)

// ExampleAgent demonstrates a simple agent implementation
type ExampleAgent struct {
	state   state.AgentState // Use the shared interface
	planner Planner
}

func (a *ExampleAgent) Init(ctx context.Context, params map[string]interface{}) error {
	a.state = state.NewInMemoryState() // Returns AgentState
	a.planner = &SimplePlanner{}
	return nil
}

// Invoke synchronously produces a response.
func (a *ExampleAgent) Invoke(ctx context.Context, input string) (*AgentCallResult, error) {
	a.state.Set("last_input", input)
	response := fmt.Sprintf("Echo: %s", input)
	result := &AgentCallResult{
		Output:   response,
		Metadata: map[string]interface{}{"agent": "ExampleAgent"},
	}
	return result, nil
}

// InvokeAsync demonstrates async by sending one result on a channel.
func (a *ExampleAgent) InvokeAsync(ctx context.Context, input string) (<-chan AgentCallResult, error) {
	ch := make(chan AgentCallResult, 1)
	go func() {
		res, _ := a.Invoke(ctx, input)
		ch <- *res
		close(ch)
	}()
	return ch, nil
}

func (a *ExampleAgent) State() state.AgentState {
	return a.state
}

func (a *ExampleAgent) Planner() Planner {
	return a.planner
}

// Delegate passes the input to other agents and collects their results.
func (a *ExampleAgent) Delegate(ctx context.Context, input string, agents []Agent) ([]*AgentCallResult, error) {
	results := make([]*AgentCallResult, 0, len(agents))
	for _, ag := range agents {
		res, err := ag.Invoke(ctx, input)
		if err != nil {
			res = &AgentCallResult{Error: err}
		}
		results = append(results, res)
	}
	return results, nil
}
