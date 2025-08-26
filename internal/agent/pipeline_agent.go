package agent

import (
	"context"
	"fmt"
	"gogurt/internal/state"
	"gogurt/internal/types"
)

// PipelineStep is a single step, which should update state as needed and return output and state.
type PipelineStep func(ctx context.Context, input any, st state.AgentState) (any, state.AgentState, error)

type PipelineAgent struct {
	Steps []PipelineStep
	State state.AgentState // current state of the agent, updated after each step
}

func (p *PipelineAgent) Invoke(ctx context.Context, input string) (*types.AgentCallResult, error) {
	var data any = input
	var err error
	currentState := p.State
	for _, step := range p.Steps {
		var newState state.AgentState
		data, newState, err = step(ctx, data, currentState)
		if err != nil {
			return nil, err
		}
		if newState != nil {
			currentState = newState
		}
	}
	p.State = currentState
	return &types.AgentCallResult{
		Output: fmt.Sprintf("%v", data),
		Metadata: map[string]any{
			"pipeline": true,
			"state":    currentState,
		},
	}, nil
}

// Async Invoke
func (p *PipelineAgent) AInvoke(ctx context.Context, input string) (<-chan *types.AgentCallResult, <-chan error) {
	resultCh := make(chan *types.AgentCallResult, 1)
	errorCh := make(chan error, 1)
	go func() {
		defer close(resultCh)
		defer close(errorCh)
		res, err := p.Invoke(ctx, input)
		if err != nil {
			errorCh <- err
			return
		}
		resultCh <- res
	}()
	return resultCh, errorCh
}