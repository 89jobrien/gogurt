package pipes

import (
	"context"
	"encoding/json"
	"fmt"
	"gogurt/internal/agent"
	"gogurt/internal/config"
	"gogurt/internal/factories"
	"gogurt/internal/tools"
)

// WorkflowPipe orchestrates a multi-step task by first planning and then executing.
type WorkflowPipe struct {
	planner agent.Agent
	worker  agent.Agent
}

// NewWorkflowPipe creates a new WorkflowPipe.
func NewWorkflowPipe(ctx context.Context, cfg *config.Config) (*WorkflowPipe, error) {
	llm := factories.GetLLM(cfg)
	registry := tools.NewRegistry()
	// Register all simple tools for the workflow
	errs := registry.RegisterBatch([]*tools.Tool{
		tools.UppercaseTool,
		tools.ConcatenateTool,
		tools.ReverseTool,
		tools.PalindromeTool,
		tools.AddTool,
		tools.SubtractTool,
		tools.MultiplyTool,
		tools.DivideTool,
	})
	for _, err := range errs {
		if err != nil {
			// In a real application, you might want to handle this more gracefully
			fmt.Printf("Warning: could not register tool: %v\n", err)
		}
	}

	planner := agent.NewPlannerAgent(llm, registry)
	worker := agent.NewWorkerAgent(registry)

	return &WorkflowPipe{
		planner: planner,
		worker:  worker,
	}, nil
}

// Run executes the full plan-and-execute workflow.
func (p *WorkflowPipe) Run(ctx context.Context, prompt string) (string, error) {
	// 1. Plan the steps
	planResult, err := p.planner.Invoke(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("planning phase failed: %w", err)
	}

	plan, ok := planResult.([]agent.PlannedStep)
	if !ok {
		return "", fmt.Errorf("planner returned invalid type: expected []agent.PlannedStep, got %T", planResult)
	}

	if len(plan) == 0 {
		return "No plan was generated to achieve the goal.", nil
	}

	// 2. Execute the steps sequentially
	var lastResult any

	for i, step := range plan {
		// For steps after the first, explicitly inject the previous result
		// as the primary argument. This is a convention-based approach,
		// assuming the main input field for text tools is "Text".
		if i > 0 {
			if step.Args == nil {
				step.Args = make(map[string]any)
			}
			// Overwrite the "Text" argument with the output of the previous step.
			step.Args["Text"] = lastResult
		}

		argsJSON, err := json.Marshal(step.Args)
		if err != nil {
			return "", fmt.Errorf("failed to marshal args for step %d: %w", i+1, err)
		}

		// The worker agent expects a single string in the format "tool_name:json_args"
		task := fmt.Sprintf("%s:%s", step.Tool, string(argsJSON))

		result, err := p.worker.Invoke(ctx, task)
		if err != nil {
			return "", fmt.Errorf("execution of step %d ('%s') failed: %w", i+1, step.Tool, err)
		}
		lastResult = result
	}

	return fmt.Sprintf("%v", lastResult), nil
}

// ARun is the asynchronous version of Run.
func (p *WorkflowPipe) ARun(ctx context.Context, prompt string) (<-chan string, <-chan error) {
	resultCh := make(chan string, 1)
	errorCh := make(chan error, 1)
	go func() {
		defer close(resultCh)
		defer close(errorCh)
		res, err := p.Run(ctx, prompt)
		if err != nil {
			errorCh <- err
			return
		}
		resultCh <- res
	}()
	return resultCh, errorCh
}
