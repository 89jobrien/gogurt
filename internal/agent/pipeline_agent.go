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
	return &types.AgentCallResult{
		Output:   fmt.Sprintf("%v", data),
		Metadata: map[string]any{"pipeline": true},
	}, nil
}

// Async Invoke
func (p *PipelineAgent) AInvoke(ctx context.Context, input string) (<-chan *types.AgentCallResult, <-chan error) {
	resultCh := make(chan *types.AgentCallResult, 1)
	errorCh := make(chan error, 1)
	go func() {
		defer close(resultCh)
		defer close(errorCh)
		res, err := p.Invoke(ctx, input)
		if err != nil {
			errorCh <- err
			return
		}
		resultCh <- res
	}()
	return resultCh, errorCh
}
