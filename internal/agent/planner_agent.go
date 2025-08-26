package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"gogurt/internal/llm"
	"gogurt/internal/logger"
	"gogurt/internal/state"
	"gogurt/internal/tools"
	"gogurt/internal/types"
	"strings"
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
	tools *tools.Registry
}

// NewPlannerAgent creates a new PlannerAgent.
func NewPlannerAgent(llm llm.LLM, registry *tools.Registry) Agent {
	logger.Info("Creating PlannerAgent:\n    LLM: %v\n   Registry: %v", llm, registry)
	return &PlannerAgent{
		llm:   llm,
		state: state.NewMemoryState(),
		tools: registry,
	}
}

// Init initializes the agent with a given configuration.
func (a *PlannerAgent) Init(ctx context.Context, config types.AgentConfig) error {
	// Initialization logic for PlannerAgent, if any, would go here.
	// For now, we can leave it empty.
	return nil
}

// Invoke takes a high-level goal as a string and returns a plan as a slice of PlannedSteps.
func (a *PlannerAgent) Invoke(ctx context.Context, input any) (any, error) {
	goal, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("invalid input type for PlannerAgent: expected string, got %T", input)
	}
	logger.Info("PlannerAgent invoked with goal: %v", goal)

	toolDescriptions := []string{}
	for _, tool := range a.tools.ListTools() {
		toolDescriptions = append(toolDescriptions, tool.Describe())
	}

	prompt := fmt.Sprintf(
		"Based on the user's goal, create a plan consisting of a sequence of tool calls. "+
			"Here are the available tools:\n\n%s\n\n"+
			"Goal: %s\n\n"+
			"Return ONLY the plan as a JSON array of objects, where each object has a 'tool' and 'args' key. "+
			"For example: [{\"tool\": \"read_file\", \"args\": {\"Filename\": \"example.txt\"}}]",
		strings.Join(toolDescriptions, "\n"),
		goal,
	)

	messages := []types.ChatMessage{
		{Role: types.RoleSystem, Content: "You are a planning agent that creates a sequence of tool calls to achieve a goal."},
		{Role: types.RoleUser, Content: prompt},
	}
	logger.Info("Messages: %v", messages)

	resp, err := a.llm.Generate(ctx, messages)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to generate plan: %v", err)
		return nil, fmt.Errorf("failed to generate plan: %w", err)
	}

	if resp == nil {
		logger.ErrorCtx(ctx, "No response from LLM")
		return nil, fmt.Errorf("no response from LLM")
	}
	
	logger.InfoCtx(ctx, "LLM response: %v", resp.Content)

	// Extract the JSON array from the LLM's response.
	jsonContent := extractJSONArray(resp.Content)
	if jsonContent == "" {
		logger.ErrorCtx(ctx, "No JSON array found in LLM response: %v", resp.Content)
		return nil, fmt.Errorf("no JSON array found in LLM response: %v", resp.Content)
	}

	var plan []PlannedStep
	if err := json.Unmarshal([]byte(jsonContent), &plan); err != nil {
		logger.ErrorCtx(ctx, "Failed to unmarshal plan from LLM response: %v. Response content: %v", err, jsonContent)
		return nil, fmt.Errorf("failed to unmarshal plan from LLM response: %v. Response content: %v", err, jsonContent)
	}

	// Add the plan to the state.
	err = a.state.Set("plan", plan)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to add plan to state: %v", err)
		return nil, fmt.Errorf("failed to add plan to state: %w", err)
	}

	logger.InfoCtx(ctx, "Plan: %v", plan)

	return plan, nil
}

// extractJSONArray finds and returns the first JSON array from a string.
func extractJSONArray(s string) string {
	start := strings.Index(s, "[")
	end := strings.LastIndex(s, "]")

	if start != -1 && end != -1 && start < end {
		return s[start : end+1]
	}

	return ""
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
	logger.Info("resultCh: %v", resultCh)
	return resultCh, errorCh
}

// OnMessage handles agent-to-agent communication.
func (a *PlannerAgent) OnMessage(ctx context.Context, msg *types.StateMessage) (*types.StateMessage, error) {
	plan, err := a.Invoke(ctx, msg.Message)
	if err != nil {
		return nil, err
	}
	planBytes, _ := json.Marshal(plan)
	logger.Info("Plan: %v", plan)
	return NewStateMessage(types.RoleAssistant, string(planBytes)), nil
}

// OnMessageAsync is the asynchronous version of OnMessage.
func (a *PlannerAgent) OnMessageAsync(ctx context.Context, msg *types.StateMessage) (<-chan *types.StateMessage, <-chan error) {
	logger.Info("Received message from %s: %s", msg.Sender, msg.Message)
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
	logger.Info("State: %v", a.state)
	return &a.state
}

// Describe returns a description of the agent.
func (a *PlannerAgent) Describe() *types.AgentDescription {
	logger.Info("PlannerAgent Describe: %v", a.state)
	return &types.AgentDescription{
		Name:         "PlannerAgent",
		Capabilities: []string{"task-planning", "decomposition"},
	}
}

func init() {
	RegisterAgent("PlannerAgent", func() Agent {
		logger.Info("Initializing PlannerAgent")
		return &PlannerAgent{}
	})
}
