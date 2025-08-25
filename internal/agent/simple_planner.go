package agent

import (
	"context"
	"gogurt/internal/state"
)

// SimplePlanner breaks the goal into words as "steps"
type SimplePlanner struct{}

func (p *SimplePlanner) Plan(ctx context.Context, goal string, state state.AgentState) ([]string, error) {
	// Example: split input into fake sequential steps
	// In practice, this could be a sophisticated decomposition or reasoning engine
	return []string{"Step 1: " + goal, "Step 2: Complete"}, nil
}
