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

// SerpApiPipe orchestrates a multi-step task by first planning, then executing, and finally synthesizing a result.
type SerpApiPipe struct {
	planner agent.Agent
	worker  agent.Agent
	llm     llm.LLM
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
		llm:     llm,
	}, nil
}

// Run executes the full plan-and-execute workflow asynchronously.
func (p *SerpApiPipe) Run(ctx context.Context, prompt string) (<-chan string, <-chan error) {
	resultCh := make(chan string, 1)
	errorCh := make(chan error, 1)

	go func() {
		defer close(resultCh)
		defer close(errorCh)

		logger.Info("Running SerpApiPipe")
		// 1. Set up prompt for the planner
		registry := tools.NewRegistry()
		registry.RegisterBatch([]*tools.Tool{
			stateful.ReadScratchpadTool,
			stateful.SaveToScratchpadTool,
			web.SerpAPISearchTool,
		})

		var toolDescriptions []string
		for _, tool := range registry.ListTools() {
			toolDescriptions = append(toolDescriptions, tool.Describe())
		}

		tmpl, err := prompts.NewPromptTemplate(planner.SerpApiPlannerPrompt)
		if err != nil {
			errorCh <- fmt.Errorf("failed to create planner prompt template: %w", err)
			return
		}

		plannerPrompt, err := tmpl.Format(map[string]string{
			"tool_descriptions": strings.Join(toolDescriptions, "\n"),
			"goal":              prompt,
		})
		if err != nil {
			errorCh <- fmt.Errorf("failed to format planner prompt: %w", err)
			return
		}

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
			logger.Info("No plan was generated to achieve the goal.")
			resultCh <- "No plan was generated to achieve the goal."
			return
		}
		logger.Info("Plan: %v", plan)

		// 3. Execute the steps sequentially
		var lastResult any
		for i, step := range plan {
			argsJSON, err := json.Marshal(step.Args)
			if err != nil {
				errorCh <- fmt.Errorf("failed to marshal args for step %d: %w", i+1, err)
				return
			}

			logger.Info("Executing step %d ('%s') with args: %s", i+1, step.Tool, string(argsJSON))
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

		// 4. Synthesize the final answer asynchronously
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

		finalAnswerCh, synthErrCh := p.llm.AGenerate(ctx, synthesisMessages)
		select {
		case finalAnswer := <-finalAnswerCh:
			if finalAnswer == nil {
				errorCh <- fmt.Errorf("no response from LLM for synthesis")
				return
			}
			logger.Info("Result: %v", finalAnswer.Content)
			resultCh <- finalAnswer.Content
		case err := <-synthErrCh:
			errorCh <- fmt.Errorf("final answer synthesis failed: %w", err)
			return
		case <-ctx.Done():
			errorCh <- ctx.Err()
			return
		}
	}()

	return resultCh, errorCh
}