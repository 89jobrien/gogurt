package agent

import (
	"context"
	"fmt"

	"gogurt/internal/types"
)

type PipelineStep func(ctx context.Context, input any) (any, error)

type PipelineAgent struct {
	Steps []PipelineStep
}

func (p *PipelineAgent) Invoke(ctx context.Context, input string) (*types.AgentCallResult, error) {
	var data any = input
	var err error
	for _, step := range p.Steps {
		data, err = step(ctx, data)
		if err != nil {
			return nil, err
		}
	}
	return &types.AgentCallResult{Output: fmt.Sprintf("%v", data), Metadata: map[string]any{"pipeline": true}}, nil
}
