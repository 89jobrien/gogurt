package agent

import (
	"context"
	"gogurt/internal/state"
)

type Planner interface {
	Plan(ctx context.Context, goal string, state state.AgentState) ([]string, error)
	APlan(ctx context.Context, goal string, state state.AgentState) (<-chan []string, <-chan error)
}

// synchronous-to-async adapter for any Planner implementation:
func APlanAdapter(p Planner, ctx context.Context, goal string, state state.AgentState) (<-chan []string, <-chan error) {
	planCh := make(chan []string, 1)
	errCh := make(chan error, 1)
	go func() {
		defer close(planCh)
		defer close(errCh)
		plan, err := p.Plan(ctx, goal, state)
		if err != nil {
			errCh <- err
			return
		}
		planCh <- plan
	}()
	return planCh, errCh
}
