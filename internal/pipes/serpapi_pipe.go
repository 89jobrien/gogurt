package pipes

import (
	"context"
	"encoding/json"
	"fmt"
	"gogurt/internal/agent"
	"gogurt/internal/config"
	"gogurt/internal/factories"
	"gogurt/internal/llm"
	"gogurt/internal/logger"
	"gogurt/internal/prompts"
	"gogurt/internal/prompts/planner"
	"gogurt/internal/tools"
	"gogurt/internal/tools/stateful"
	"gogurt/internal/tools/web"
	"gogurt/internal/types"
	"strings"
)

// SerpApiPipe orchestrates a multi-step task by first planning and then executing.
type SerpApiPipe struct {
	planner agent.Agent
	worker  agent.Agent
	llm     llm.LLM // Added the LLM for the synthesis step
}

// NewSerpApiPipe creates a new SerpApiPipe.
func NewSerpApiPipe(ctx context.Context, cfg *config.Config) (*SerpApiPipe, error) {
	llm := factories.GetLLM(cfg)
	registry := tools.NewRegistry()
	errs := registry.RegisterBatch([]*tools.Tool{
		stateful.ReadScratchpadTool,
		stateful.SaveToScratchpadTool,
		web.SerpAPISearchTool,
	})
	for _, err := range errs {
		if err != nil {
			fmt.Printf("Warning: could not register tool: %v\n", err)
		}
	}

	planner := agent.NewPlannerAgent(llm)
	worker := agent.NewWorkerAgent(registry)

	return &SerpApiPipe{
		planner: planner,
		worker:  worker,
		llm:     llm, // Populate the new llm field
	}, nil
}

// Run executes the full plan-and-execute workflow.
func (p *SerpApiPipe) Run(ctx context.Context, prompt string) (string, error) {
	logger.Info("Running SerpApiPipe")
	registry := tools.NewRegistry()
	errs := registry.RegisterBatch([]*tools.Tool{
		stateful.ReadScratchpadTool,
		stateful.SaveToScratchpadTool,
		web.SerpAPISearchTool,
	})
	for _, err := range errs {
		if err != nil {
			fmt.Printf("Warning: could not register tool: %v\n", err)
		}
	}

	toolDescriptions := []string{}
	for _, tool := range registry.ListTools() {
		toolDescriptions = append(toolDescriptions, tool.Describe())
	}

	// Use the correct, updated prompt template
	tmpl, err := prompts.NewPromptTemplate(planner.SerpApiPlannerPrompt)
	if err != nil {
		logger.Error("Failed to create planner prompt template: %v", err)
		return "", fmt.Errorf("failed to create planner prompt template: %w", err)
	}

	plannerPrompt, err := tmpl.Format(map[string]string{
		"tool_descriptions": strings.Join(toolDescriptions, "\n"),
		"goal":              prompt,
	})
	if err != nil {
		logger.Error("Failed to format planner prompt: %v", err)
		return "", fmt.Errorf("failed to format planner prompt: %w", err)
	}

	// 2. Plan the steps
	planResult, err := p.planner.Invoke(ctx, plannerPrompt)
	if err != nil {
		logger.Error("Planning phase failed: %v", err)
		return "", fmt.Errorf("planning phase failed: %w", err)
	}

	plan, ok := planResult.([]agent.PlannedStep)
	if !ok {
		logger.Error("Planner returned invalid type: expected []agent.PlannedStep, got %T", planResult)
		return "", fmt.Errorf("planner returned invalid type: expected []agent.PlannedStep, got %T", planResult)
	}

	if len(plan) == 0 {
		logger.Info("No plan was generated to achieve the goal.")
		return "No plan was generated to achieve the goal.", nil
	}

	logger.Info("Plan: %v", plan)
	// 3. Execute the steps sequentially
	var lastResult any

	for i, step := range plan {
		argsJSON, err := json.Marshal(step.Args)
		if err != nil {
			logger.Error("Failed to marshal args for step %d: %v", i+1, err)
			return "", fmt.Errorf("failed to marshal args for step %d: %w", i+1, err)
		}

		logger.Info("Executing step %d ('%s') with args: %s", i+1, step.Tool, string(argsJSON))
		task := fmt.Sprintf("%s:%s", step.Tool, string(argsJSON))

		result, err := p.worker.Invoke(ctx, task)
		if err != nil {
			logger.Error("Execution of step %d ('%s') failed: %v", i+1, step.Tool, err)
			return "", fmt.Errorf("execution of step %d ('%s') failed: %w", i+1, step.Tool, err)
		}
		lastResult = result
	}

	// 4. NEW: Synthesize the final answer
	logger.Info("Synthesizing final answer from tool results.")
	synthesisPrompt := fmt.Sprintf(
		"Based on the following information, please provide a direct answer to the user's original question.\n\n"+
			"Information:\n%v\n\n"+
			"Original Question: %s",
		lastResult,
		prompt,
	)

	synthesisMessages := []types.ChatMessage{
		{Role: types.RoleSystem, Content: "You are a helpful assistant that answers questions based on provided context."},
		{Role: types.RoleUser, Content: synthesisPrompt},
	}

	finalAnswer, err := p.llm.Generate(ctx, synthesisMessages)
	if err != nil {
		logger.Error("Final answer synthesis failed: %v", err)
		return "", fmt.Errorf("final answer synthesis failed: %w", err)
	}

	if finalAnswer == nil {
		logger.Error("LLM returned a nil response for synthesis")
		return "", fmt.Errorf("no response from LLM for synthesis")
	}

	logger.Info("Result: %v", finalAnswer.Content)
	return finalAnswer.Content, nil
}

// ARun is the asynchronous version of Run.
func (p *SerpApiPipe) ARun(ctx context.Context, prompt string) (<-chan string, <-chan error) {
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