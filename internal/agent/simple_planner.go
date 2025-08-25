package agent

import (
	"context"
	"gogurt/internal/state"
)

type SimplePlanner struct{}

func (p *SimplePlanner) Plan(ctx context.Context, goal string, state state.AgentState) ([]string, error) {
	// Example: split input into fake sequential steps
	return []string{"Step 1: " + goal, "Step 2: Complete"}, nil
}
