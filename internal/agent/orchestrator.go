package agent

import (
	"context"
	"fmt"
	"sync"

	"gogurt/internal/types"
)

// Orchestrator coordinates multiple agents
type Orchestrator struct {
	Agents []Agent
}

// RunParallel invokes all agents with the same input and collects results
func (o *Orchestrator) RunParallel(ctx context.Context, input string) ([]*types.AgentCallResult, error) {
	var wg sync.WaitGroup
	results := make([]*types.AgentCallResult, len(o.Agents))
	errCh := make(chan error, len(o.Agents))
	for i, agent := range o.Agents {
		wg.Add(1)
		go func(idx int, ag Agent) {
			defer wg.Done()
			res, err := ag.Invoke(ctx, input)
			if err != nil {
				results[idx] = &types.AgentCallResult{Error: err}
				errCh <- err
			} else if acr, ok := res.(*types.AgentCallResult); ok {
				results[idx] = acr
			} else {
				results[idx] = &types.AgentCallResult{
					Error: fmt.Errorf("invalid AgentCallResult (got %T)", res),
				}
				errCh <- results[idx].Error
			}
		}(i, agent)
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		return results, err
	}
	return results, nil
}

func (o *Orchestrator) RunPiped(ctx context.Context, input string) (*types.AgentCallResult, error) {
	var pipedRes *types.AgentCallResult
	currentInput := input
	for _, agent := range o.Agents {
		res, err := agent.Invoke(ctx, currentInput)
		var acr *types.AgentCallResult
		if r, ok := res.(*types.AgentCallResult); ok {
			acr = r
		} else if err != nil {
			acr = &types.AgentCallResult{Error: err}
		} else {
			acr = &types.AgentCallResult{
				Error: fmt.Errorf("invalid AgentCallResult (got %T)", res),
			}
		}
		pipedRes = acr
		if pipedRes.Error != nil {
			return pipedRes, pipedRes.Error
		}
		currentInput = pipedRes.Output
	}
	return pipedRes, nil
}