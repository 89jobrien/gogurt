package agent

import (
	"context"
	"fmt"
	"gogurt/internal/state"
	"gogurt/internal/types"
	"sync"
)

// Orchestrator coordinates multiple agents, threading state updates.
type Orchestrator struct {
	Agents []Agent
	State  *state.AgentState
}

// RunParallel invokes all agents with the same input and collects results (and states) concurrently.
func (o *Orchestrator) RunParallel(ctx context.Context, input string) (<-chan *types.AgentCallResult, <-chan error) {
	resultsCh := make(chan *types.AgentCallResult, len(o.Agents))
	errCh := make(chan error, 1)

	go func() {
		defer close(resultsCh)
		defer close(errCh)

		var wg sync.WaitGroup
		for _, agent := range o.Agents {
			wg.Add(1)
			go func(ag Agent) {
				defer wg.Done()
				resCh, agentErrCh := ag.Invoke(ctx, input)
				select {
				case res := <-resCh:
					if acr, ok := res.(*types.AgentCallResult); ok {
						if acr.Metadata == nil {
							acr.Metadata = make(map[string]any)
						}
						acr.Metadata["state"] = ag.State()
						resultsCh <- acr
					} else {
						resultsCh <- &types.AgentCallResult{
							Error:    fmt.Errorf("invalid AgentCallResult (got %T)", res),
							Metadata: map[string]any{"state": ag.State()},
						}
					}
				case err := <-agentErrCh:
					resultsCh <- &types.AgentCallResult{Error: err, Metadata: map[string]any{"state": ag.State()}}
				case <-ctx.Done():
					resultsCh <- &types.AgentCallResult{Error: ctx.Err(), Metadata: map[string]any{"state": ag.State()}}
				}
			}(agent)
		}
		wg.Wait()
	}()

	return resultsCh, errCh
}

// RunPiped runs agents sequentially, using each output and updating state.
func (o *Orchestrator) RunPiped(ctx context.Context, input string) (<-chan *types.AgentCallResult, <-chan error) {
	resultCh := make(chan *types.AgentCallResult, 1)
	errCh := make(chan error, 1)

	go func() {
		defer close(resultCh)
		defer close(errCh)

		var finalResult *types.AgentCallResult
		currentInput := input

		for _, agent := range o.Agents {
			resCh, agentErrCh := agent.Invoke(ctx, currentInput)
			select {
			case res := <-resCh:
				var acr *types.AgentCallResult
				if r, ok := res.(*types.AgentCallResult); ok {
					acr = r
				} else {
					acr = &types.AgentCallResult{
						Error: fmt.Errorf("invalid AgentCallResult (got %T)", res),
					}
				}

				if acr.Metadata == nil {
					acr.Metadata = make(map[string]any)
				}
				acr.Metadata["state"] = agent.State()

				finalResult = acr
				o.State = agent.State()

				if finalResult.Error != nil {
					resultCh <- finalResult
					return
				}
				currentInput = fmt.Sprintf("%v", finalResult.Output)
			case err := <-agentErrCh:
				errCh <- err
				return
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			}
		}
		resultCh <- finalResult
	}()

	return resultCh, errCh
}