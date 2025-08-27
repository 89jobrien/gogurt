package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"gogurt/internal/llm"
	"gogurt/internal/logger"
	"gogurt/internal/state"
	"gogurt/internal/types"
	"gogurt/internal/utils"
)

// PlannedStep defines the structure for a single step in the generated plan.
type PlannedStep struct {
	Tool string         `json:"tool"`
	Args map[string]any `json:"args"`
}

// PlannerAgent is responsible for breaking down a high-level goal into a sequence of tool calls.
type PlannerAgent struct {
	llm   llm.LLM
	state state.AgentState
}

// NewPlannerAgent creates a new PlannerAgent.
func NewPlannerAgent(llm llm.LLM) Agent {
	logger.Info("Creating PlannerAgent with LLM: %v", llm)
	return &PlannerAgent{
		llm:   llm,
		state: state.NewMemoryState(),
	}
}

// Init initializes the agent with a given configuration.
func (a *PlannerAgent) Init(ctx context.Context, config types.AgentConfig) error {
	return nil
}

// Invoke takes a prompt as a string and returns a plan as a slice of PlannedSteps.
func (a *PlannerAgent) Invoke(ctx context.Context, input any) (any, error) {
	prompt, ok := input.(string)
	if !ok {
		logger.ErrorCtx(ctx, "Invalid input type for PlannerAgent: expected string, got %T", input)
		return nil, fmt.Errorf("invalid input type for PlannerAgent: expected string, got %T", input)
	}

	logger.InfoCtx(ctx, "PlannerAgent invoked with prompt.")

	messages := []types.ChatMessage{
		{Role: types.RoleSystem, Content: "You are a planning agent that creates a sequence of tool calls to achieve a goal."},
		{Role: types.RoleUser, Content: prompt},
	}

	resp, err := a.llm.Generate(ctx, messages)
	if err != nil {
		logger.ErrorCtx(ctx, "LLM plan generation failed: %v", err)
		return nil, fmt.Errorf("failed to generate plan: %w", err)
	}

	if resp == nil {
		logger.ErrorCtx(ctx, "LLM returned a nil response")
		return nil, fmt.Errorf("no response from LLM")
	}
	logger.InfoCtx(ctx, "LLM response received: %s", resp.Content)

	jsonContent := utils.ExtractJSONArray(resp.Content)
	if jsonContent == "" {
		logger.WarnCtx(ctx, "No JSON array found in LLM response: %s", resp.Content)
		return nil, fmt.Errorf("no JSON array found in LLM response: %s", resp.Content)
	}

	var plan []PlannedStep
	if err := json.Unmarshal([]byte(jsonContent), &plan); err != nil {
		logger.ErrorCtx(ctx, "Failed to unmarshal plan from LLM response: %v. Content: %s", err, jsonContent)
		return nil, fmt.Errorf("failed to unmarshal plan from LLM response: %w. Response content: %s", err, jsonContent)
	}

	a.state.Set("plan", plan)
	logger.InfoCtx(ctx, "Plan generated successfully: %v", plan)
	return plan, nil
}

// InvokeAsync is the asynchronous version of Invoke.
func (a *PlannerAgent) InvokeAsync(ctx context.Context, input any) (<-chan any, <-chan error) {
	resultCh := make(chan any, 1)
	errorCh := make(chan error, 1)
	go func() {
		defer close(resultCh)
		defer close(errorCh)
		res, err := a.Invoke(ctx, input)
		if err != nil {
			errorCh <- err
			return
		}
		resultCh <- res
	}()
	return resultCh, errorCh
}

// OnMessage handles agent-to-agent communication.
func (a *PlannerAgent) OnMessage(ctx context.Context, msg *types.StateMessage) (*types.StateMessage, error) {
	plan, err := a.Invoke(ctx, msg.Message)
	if err != nil {
		return nil, err
	}
	planBytes, _ := json.Marshal(plan)
	return NewStateMessage(types.RoleAssistant, string(planBytes)), nil
}

// OnMessageAsync is the asynchronous version of OnMessage.
func (a *PlannerAgent) OnMessageAsync(ctx context.Context, msg *types.StateMessage) (<-chan *types.StateMessage, <-chan error) {
	resultCh := make(chan *types.StateMessage, 1)
	errorCh := make(chan error, 1)
	go func() {
		defer close(resultCh)
		defer close(errorCh)
		res, err := a.OnMessage(ctx, msg)
		if err != nil {
			errorCh <- err
			return
		}
		resultCh <- res
	}()
	return resultCh, errorCh
}

// State returns the agent's current state.
func (a *PlannerAgent) State() *state.AgentState {
	return &a.state
}

// Describe returns a description of the agent.
func (a *PlannerAgent) Describe() *types.AgentDescription {
	return &types.AgentDescription{
		Name:         "PlannerAgent",
		Capabilities: []string{"task-planning", "decomposition"},
	}
}

func init() {
	RegisterAgent("PlannerAgent", func() Agent {
		return &PlannerAgent{}
	})
}