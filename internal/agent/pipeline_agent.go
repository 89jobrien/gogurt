package agent

import (
	"context"
	"fmt"
)

type PipelineStep func(ctx context.Context, input any) (any, error)

type PipelineAgent struct {
    Steps []PipelineStep
}

func (p *PipelineAgent) Invoke(ctx context.Context, input string) (*AgentCallResult, error) {
    var data any = input
    var err error
    for _, step := range p.Steps {
        data, err = step(ctx, data)
        if err != nil {
            return nil, err
        }
    }
    return &AgentCallResult{Output: fmt.Sprintf("%v", data), Metadata: map[string]any{"pipeline": true}}, nil
}
