package pipes

import (
	"context"

	"gogurt/internal/agent"
)

// OrchestratorPipe uses an Orchestrator to run a workflow.
type OrchestratorPipe struct {
	Orchestrator *agent.Orchestrator
}

// NewOrchestratorPipe creates a new pipe from an agent orchestrator.
func NewOrchestratorPipe(orch *agent.Orchestrator) Pipe {
	return &OrchestratorPipe{Orchestrator: orch}
}

// Run provides a non-blocking, asynchronous call to the orchestration workflow.
func (p *OrchestratorPipe) Run(ctx context.Context, prompt string) (<-chan string, <-chan error) {
	outCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		defer close(outCh)
		defer close(errCh)

		// Call the asynchronous RunPiped method from the orchestrator
		resultCh, orchErrCh := p.Orchestrator.RunPiped(ctx, prompt)

		select {
		case result := <-resultCh:
			// Handle the result from the orchestrator
			if result == nil {
				outCh <- ""
				return
			}
			if result.Error != nil {
				errCh <- result.Error
				return
			}
			outCh <- result.Output
		case err := <-orchErrCh:
			// Handle any errors from the orchestrator's execution
			errCh <- err
		case <-ctx.Done():
			// Handle context cancellation
			errCh <- ctx.Err()
		}
	}()

	return outCh, errCh
}