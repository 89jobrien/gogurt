package pipes

import (
	"context"
	"encoding/json"
	"fmt"
	"gogurt/internal/agent"
	"gogurt/internal/config"
	"gogurt/internal/factories"
	"gogurt/internal/tools"
	"gogurt/internal/tools/file_tools"
	"strings"
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
		file_tools.ReadFileTool,
		file_tools.WriteFileTool,
		file_tools.ListFilesTool,
	})
	for _, err := range errs {
		if err != nil {
			// In a real application, you might want to handle this more gracefully
			fmt.Printf("Warning: could not register tool: %v\n", err)
		}
	}

	planner := agent.NewPlannerAgent(llm)
	worker := agent.NewWorkerAgent(registry)

	return &WorkflowPipe{
		planner: planner,
		worker:  worker,
	}, nil
}

// Run executes the full plan-and-execute workflow asynchronously.
func (p *WorkflowPipe) Run(ctx context.Context, prompt string) (<-chan string, <-chan error) {
	resultCh := make(chan string, 1)
	errorCh := make(chan error, 1)

	go func() {
		defer close(resultCh)
		defer close(errorCh)

		// 1. Create the prompt for the planner
		registry := tools.NewRegistry()
		registry.RegisterBatch([]*tools.Tool{
			tools.UppercaseTool,
			tools.ConcatenateTool,
			tools.ReverseTool,
			tools.PalindromeTool,
			tools.AddTool,
			tools.SubtractTool,
			tools.MultiplyTool,
			tools.DivideTool,
			file_tools.ReadFileTool,
			file_tools.WriteFileTool,
			file_tools.ListFilesTool,
		})

		var toolDescriptions []string
		for _, tool := range registry.ListTools() {
			toolDescriptions = append(toolDescriptions, tool.Describe())
		}

		plannerPrompt := fmt.Sprintf(
			"Based on the user's goal, create a plan consisting of a sequence of tool calls. "+
				"Here are the available tools:\n\n%s\n\n"+
				"Goal: %s\n\n"+
				"Return ONLY a valid, flat JSON array of objects, where each object has a 'tool' and 'args' key. "+
				"Do not include any comments or nested arrays. For example: "+
				"[{\"tool\": \"duckduckgo_search\", \"args\": {\"query\": \"What is the capital of New Jersey?\", \"num_results\": 3}}]",
			strings.Join(toolDescriptions, "\n"),
			prompt,
		)

		// 2. Plan the steps asynchronously
		planResultCh, planErrCh := p.planner.Invoke(ctx, plannerPrompt)
		var plan []agent.PlannedStep

		select {
		case planResult := <-planResultCh:
			var ok bool
			plan, ok = planResult.([]agent.PlannedStep)
			if !ok {
				errorCh <- fmt.Errorf("planner returned invalid type: expected []agent.PlannedStep, got %T", planResult)
				return
			}
		case err := <-planErrCh:
			errorCh <- fmt.Errorf("planning phase failed: %w", err)
			return
		case <-ctx.Done():
			errorCh <- ctx.Err()
			return
		}

		if len(plan) == 0 {
			resultCh <- "No plan was generated to achieve the goal."
			return
		}

		// 3. Execute the steps sequentially
		var lastResult any
		for i, step := range plan {
			argsJSON, err := json.Marshal(step.Args)
			if err != nil {
				errorCh <- fmt.Errorf("failed to marshal args for step %d: %w", i+1, err)
				return
			}

			task := fmt.Sprintf("%s:%s", step.Tool, string(argsJSON))
			workerResultCh, workerErrCh := p.worker.Invoke(ctx, task)

			select {
			case result := <-workerResultCh:
				lastResult = result
			case err := <-workerErrCh:
				errorCh <- fmt.Errorf("execution of step %d ('%s') failed: %w", i+1, step.Tool, err)
				return
			case <-ctx.Done():
				errorCh <- ctx.Err()
				return
			}
		}

		resultCh <- fmt.Sprintf("%v", lastResult)
	}()

	return resultCh, errorCh
}