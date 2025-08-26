package pipes

import (
	"context"

	"gogurt/internal/agent"
)

type OrchestratorPipe struct {
	Orchestrator *agent.Orchestrator
}

func NewOrchestratorPipe(orch *agent.Orchestrator) Pipe {
	return &OrchestratorPipe{Orchestrator: orch}
}

func (p *OrchestratorPipe) Run(ctx context.Context, prompt string) (string, error) {
	result, err := p.Orchestrator.RunPiped(ctx, prompt)
	if err != nil {
		return "", err
	}
	if result == nil || result.Output == "" {
		return "", nil
	}
	return result.Output, nil
}

// ARun implements Pipe.
func (p *OrchestratorPipe) ARun(ctx context.Context, prompt string) (<-chan string, <-chan error) {
	outCh := make(chan string, 1)
	errCh := make(chan error, 1)
	go func() {
		defer close(outCh)
		defer close(errCh)
		result, err := p.Orchestrator.RunPiped(ctx, prompt)
		if err != nil {
			errCh <- err
			return
		}
		if result == nil || result.Output == "" {
			outCh <- ""
			return
		}
		outCh <- result.Output
	}()
	return outCh, errCh
}