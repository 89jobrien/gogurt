package agent

import (
	"context"
	"sync"
)

// Orchestrator coordinates multiple agents
type Orchestrator struct {
	Agents []Agent
}

// RunParallel invokes all agents with the same input and collects results
func (o *Orchestrator) RunParallel(ctx context.Context, input string) ([]*AgentCallResult, error) {
	var wg sync.WaitGroup
	results := make([]*AgentCallResult, len(o.Agents))
	errCh := make(chan error, len(o.Agents))

	for i, agent := range o.Agents {
		wg.Add(1)
		go func(idx int, ag Agent) {
			defer wg.Done()
			res, err := ag.Invoke(ctx, input)
			if err != nil {
				results[idx] = &AgentCallResult{Error: err}
				errCh <- err
			} else {
				results[idx] = res
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

// RunChained runs agents sequentially, using each output as the next input.
func (o *Orchestrator) RunChained(ctx context.Context, input string) (*AgentCallResult, error) {
	var chainedResult *AgentCallResult
	var err error
	currentInput := input
	for _, agent := range o.Agents {
		chainedResult, err = agent.Invoke(ctx, currentInput)
		if err != nil {
			return chainedResult, err
		}
		currentInput = chainedResult.Output
	}
	return chainedResult, nil
}
