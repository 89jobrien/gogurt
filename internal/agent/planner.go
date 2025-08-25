package agent

import (
	"context"
	"gogurt/internal/state"
)

type Planner interface {
	Plan(ctx context.Context, goal string, state state.AgentState) ([]string, error)
}
