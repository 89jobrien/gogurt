package agent

import (
	"context"
	"fmt"
	"sync"

	"gogurt/internal/state"
	"gogurt/internal/types"
)

// Orchestrator coordinates multiple agents, threading state updates.
type Orchestrator struct {
	Agents []Agent
	State  *state.AgentState
}

// RunParallel invokes all agents with the same input and collects results (and states) concurrently.
func (o *Orchestrator) RunParallel(ctx context.Context, input string) ([]*types.AgentCallResult, error) {
	var wg sync.WaitGroup
	results := make([]*types.AgentCallResult, len(o.Agents))
	states := make([]state.AgentState, len(o.Agents))
	errCh := make(chan error, len(o.Agents))
	for i, agent := range o.Agents {
		wg.Add(1)
		go func(idx int, ag Agent) {
			defer wg.Done()
			res, err := ag.Invoke(ctx, input)
			states[idx] = *ag.State()
			if err != nil {
				results[idx] = &types.AgentCallResult{Error: err, Metadata: map[string]any{"state": ag.State()}}
				errCh <- err
			} else if acr, ok := res.(*types.AgentCallResult); ok {
				if acr.Metadata == nil {
					acr.Metadata = map[string]any{}
				}
				acr.Metadata["state"] = ag.State()
				results[idx] = acr
			} else {
				results[idx] = &types.AgentCallResult{
					Error: fmt.Errorf("invalid AgentCallResult (got %T)", res),
					Metadata: map[string]any{"state": ag.State()},
				}
				errCh <- results[idx].Error
			}
		}(i, agent)
	}
	wg.Wait()
	close(errCh)
	for _, st := range states {
		if st != nil {
			*o.State = st
		}
	}
	for err := range errCh {
		return results, err
	}
	return results, nil
}

// RunPiped runs agents sequentially, using each output and updating state.
func (o *Orchestrator) RunPiped(ctx context.Context, input string) (*types.AgentCallResult, error) {
	var pipedRes *types.AgentCallResult
	currentInput := input
	for _, agent := range o.Agents {
		res, err := agent.Invoke(ctx, currentInput)
		if err != nil {
			return &types.AgentCallResult{Error: err, Metadata: map[string]any{"state": agent.State()}}, err
		}
		var acr *types.AgentCallResult
		if r, ok := res.(*types.AgentCallResult); ok {
			acr = r
		} else {
			acr = &types.AgentCallResult{
				Error: fmt.Errorf("invalid AgentCallResult (got %T)", res),
			}
		}
		if acr.Metadata == nil {
			acr.Metadata = map[string]any{}
		}
		acr.Metadata["state"] = agent.State()
		pipedRes = acr
		o.State = agent.State()
		if pipedRes.Error != nil {
			return pipedRes, pipedRes.Error
		}
		currentInput = fmt.Sprintf("%v", pipedRes.Output)
	}
	return pipedRes, nil
}