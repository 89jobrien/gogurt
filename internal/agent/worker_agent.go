package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"gogurt/internal/logger"
	"gogurt/internal/state"
	"gogurt/internal/tools"
	"gogurt/internal/types"
	"reflect"
	"strings"
)

// WorkerAgent executes a single tool call.
type WorkerAgent struct {
	state state.AgentState
	tools *tools.Registry
}

// NewWorkerAgent creates a new WorkerAgent.
func NewWorkerAgent(registry *tools.Registry) Agent {
	logger.Info("Creating WorkerAgent")
	return &WorkerAgent{
		state: state.NewMemoryState(),
		tools: registry,
	}
}

// Init initializes the agent with a given configuration.
func (a *WorkerAgent) Init(ctx context.Context, config types.AgentConfig) error {
	return nil
}

// Invoke takes a tool call string and executes it.
func (a *WorkerAgent) Invoke(ctx context.Context, input any) (any, error) {
	task, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("invalid input type for WorkerAgent: expected string, got %T", input)
	}
	logger.InfoCtx(ctx, "WorkerAgent invoked with task: %s", task)

	parts := strings.SplitN(task, ":", 2)
	toolName := parts[0]
	var args string
	if len(parts) > 1 {
		args = parts[1]
	}

	tool := a.tools.Get(toolName)
	if tool == nil {
		return nil, fmt.Errorf("tool '%s' not found", toolName)
	}

	// Dynamically handle stateful vs stateless tools
	var result any
	var err error
	if tool.Func.Type().NumIn() == 2 {
		// This is a stateful tool
		result, err = a.callStatefulTool(tool, args)
	} else {
		// This is a stateless tool
		result, err = tool.Call(args)
	}

	if err != nil {
		logger.ErrorCtx(ctx, "Tool call '%s' failed: %v", toolName, err)
		return nil, fmt.Errorf("tool call failed for '%s': %w", toolName, err)
	}

	logger.InfoCtx(ctx, "Tool '%s' executed successfully.", toolName)
	return result, nil
}

func (a *WorkerAgent) callStatefulTool(tool *tools.Tool, jsonArgs string) (any, error) {
	inputType := tool.Func.Type().In(0)
	inputValue := reflect.New(inputType).Interface()
	if err := json.Unmarshal([]byte(jsonArgs), &inputValue); err != nil {
		return nil, fmt.Errorf("error unmarshaling arguments: %w", err)
	}

	results := tool.Func.Call([]reflect.Value{
		reflect.ValueOf(inputValue).Elem(),
		reflect.ValueOf(a.state),
	})

	if !results[1].IsNil() {
		return nil, results[1].Interface().(error)
	}
	return results[0].Interface(), nil
}

// InvokeAsync is the asynchronous version of Invoke.
func (a *WorkerAgent) InvokeAsync(ctx context.Context, input any) (<-chan any, <-chan error) {
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
func (a *WorkerAgent) OnMessage(ctx context.Context, msg *types.StateMessage) (*types.StateMessage, error) {
	result, err := a.Invoke(ctx, msg.Message)
	if err != nil {
		return nil, err
	}
	resultBytes, _ := json.Marshal(result)
	return NewStateMessage(types.RoleAssistant, string(resultBytes)), nil
}

// OnMessageAsync is the asynchronous version of OnMessage.
func (a *WorkerAgent) OnMessageAsync(ctx context.Context, msg *types.StateMessage) (<-chan *types.StateMessage, <-chan error) {
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
func (a *WorkerAgent) State() *state.AgentState {
	return &a.state
}

// Describe returns a description of the agent.
func (a *WorkerAgent) Describe() *types.AgentDescription {
	return &types.AgentDescription{
		Name:         "WorkerAgent",
		Capabilities: []string{"tool-execution"},
	}
}

func init() {
	RegisterAgent("WorkerAgent", func() Agent {
		return &WorkerAgent{}
	})
}