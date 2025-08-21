package agent

import (
	"context"
	"gogurt/state"
)

type Planner interface {
    Plan(ctx context.Context, goal string, state state.AgentState) ([]string, error)
}