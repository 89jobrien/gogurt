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
			fmt.Printf("Warning: could not register tool: %v", err)
		}
	}

	planner := agent.NewPlannerAgent(llm)
	worker := agent.NewWorkerAgent(registry)

	return &DDGSPipe{
		planner: planner,
		worker:  worker,
	}, nil
}

// Run executes the full plan-and-execute workflow.
func (p *DDGSPipe) Run(ctx context.Context, prompt string) (string, error) {
	// 1. Create the prompt for the planner
	registry := tools.NewRegistry()
	errs := registry.RegisterBatch([]*tools.Tool{
		file_tools.ReadFileTool,
		file_tools.WriteFileTool,
		file_tools.ListFilesTool,
		web.DuckDuckGoSearchTool,
	})
	for _, err := range errs {
		if err != nil {
			fmt.Printf("Warning: could not register tool: %v", err)
		}
	}

	toolDescriptions := []string{}
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

	// 2. Plan the steps
	planResult, err := p.planner.Invoke(ctx, plannerPrompt)
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

	// 3. Execute the steps sequentially
	var lastResult any

	for i, step := range plan {
		argsJSON, err := json.Marshal(step.Args)
		if err != nil {
			return "", fmt.Errorf("failed to marshal args for step %d: %w", i+1, err)
		}

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
func (p *DDGSPipe) ARun(ctx context.Context, prompt string) (<-chan string, <-chan error) {
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