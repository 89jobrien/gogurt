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
	"gogurt/internal/tools/web"
	"strings"
)

// DDGSPipe orchestrates a multi-step task by first planning and then executing.
type DDGSPipe struct {
	planner agent.Agent
	worker  agent.Agent
}

// NewDDGSPipe creates a new DDGSPipe.
func NewDDGSPipe(ctx context.Context, cfg *config.Config) (*DDGSPipe, error) {
	llm := factories.GetLLM(cfg)
	registry := tools.NewRegistry()
	errs := registry.RegisterBatch([]*tools.Tool{
		file_tools.ReadFileTool,
		file_tools.WriteFileTool,
		file_tools.ListFilesTool,
		web.DuckDuckGoSearchTool,
	})
	for _, err := range errs {
		if err != nil {
			// In a real application, you'd use a structured logger.
			fmt.Printf("Warning: could not register tool: %v\n", err)
		}
	}

	planner := agent.NewPlannerAgent(llm)
	worker := agent.NewWorkerAgent(registry)

	return &DDGSPipe{
		planner: planner,
		worker:  worker,
	}, nil
}

// Run executes the full plan-and-execute workflow synchronously.
// It is a blocking wrapper around the asynchronous ARun method.
func (p *DDGSPipe) Run(ctx context.Context, prompt string) (string, error) {
	resultCh, errCh := p.ARun(ctx, prompt)
	select {
	case result := <-resultCh:
		return result, nil
	case err := <-errCh:
		return "", err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// ARun is the asynchronous version of Run. It executes the entire workflow in a non-blocking manner.
func (p *DDGSPipe) ARun(ctx context.Context, prompt string) (<-chan string, <-chan error) {
	resultCh := make(chan string, 1)
	errorCh := make(chan error, 1)

	go func() {
		defer close(resultCh)
		defer close(errorCh)

		// 1. Create the prompt for the planner
		registry := tools.NewRegistry()
		// Note: Re-registering tools here isn't ideal. In a larger app, the registry
		// would be shared or passed in. For this pipe, we'll keep it self-contained.
		registry.RegisterBatch([]*tools.Tool{
			file_tools.ReadFileTool,
			file_tools.WriteFileTool,
			file_tools.ListFilesTool,
			web.DuckDuckGoSearchTool,
		})

		var toolDescriptions []string
		for _, tool := range registry.ListTools() {
			toolDescriptions = append(toolDescriptions, tool.Describe())
		}

		plannerPrompt := fmt.Sprintf(
			"Based on the user's goal, create a plan consisting of a sequence of tool calls. "+
				"Here are the available tools:\n\n%s\n\n"+
				"Goal: %s\n\n"+
				"Important: The 'duckduckgo_search' tool directly returns the search results. No further steps are needed to read or process this output. "+
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

		// 3. Execute the steps sequentially, but each step is an async call
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