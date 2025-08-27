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

// Invoke takes a tool call string and executes it asynchronously.
func (a *WorkerAgent) Invoke(ctx context.Context, input any) (<-chan any, <-chan error) {
	resultCh := make(chan any, 1)
	errorCh := make(chan error, 1)

	go func() {
		defer close(resultCh)
		defer close(errorCh)

		task, ok := input.(string)
		if !ok {
			err := fmt.Errorf("invalid input type for WorkerAgent: expected string, got %T", input)
			errorCh <- err
			return
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
			errorCh <- fmt.Errorf("tool '%s' not found", toolName)
			return
		}

		// Handle stateful vs stateless tools asynchronously
		var toolResultCh <-chan any
		var toolErrCh <-chan error

		if tool.Func.Type().NumIn() == 2 {
			// This is a stateful tool
			toolResultCh, toolErrCh = a.asyncCallStatefulTool(ctx, tool, args)
		} else {
			// This is a stateless tool
			toolResultCh, toolErrCh = tool.AsyncCall(ctx, args)
		}

		select {
		case result := <-toolResultCh:
			logger.InfoCtx(ctx, "Tool '%s' executed successfully.", toolName)
			resultCh <- result
		case err := <-toolErrCh:
			logger.ErrorCtx(ctx, "Tool call '%s' failed: %v", toolName, err)
			errorCh <- fmt.Errorf("tool call failed for '%s': %w", toolName, err)
		case <-ctx.Done():
			errorCh <- ctx.Err()
		}
	}()

	return resultCh, errorCh
}

// asyncCallStatefulTool handles the execution of tools that require agent state.
func (a *WorkerAgent) asyncCallStatefulTool(ctx context.Context, tool *tools.Tool, jsonArgs string) (<-chan any, <-chan error) {
	resultCh := make(chan any, 1)
	errorCh := make(chan error, 1)

	go func() {
		defer close(resultCh)
		defer close(errorCh)

		inputType := tool.Func.Type().In(0)
		inputValue := reflect.New(inputType).Interface()
		if err := json.Unmarshal([]byte(jsonArgs), &inputValue); err != nil {
			logger.ErrorCtx(ctx, "Error unmarshaling arguments: %v", err)
			errorCh <- fmt.Errorf("error unmarshaling arguments: %w", err)
			return
		}

		// The actual function call is blocking, but it's wrapped in a goroutine.
		results := tool.Func.Call([]reflect.Value{
			reflect.ValueOf(inputValue).Elem(),
			reflect.ValueOf(a.state),
		})

		if !results[1].IsNil() {
			errorCh <- results[1].Interface().(error)
			return
		}
		resultCh <- results[0].Interface()
	}()

	return resultCh, errorCh
}

// OnMessage handles agent-to-agent communication asynchronously.
func (a *WorkerAgent) OnMessage(ctx context.Context, msg *types.StateMessage) (<-chan *types.StateMessage, <-chan error) {
	resultCh := make(chan *types.StateMessage, 1)
	errorCh := make(chan error, 1)

	go func() {
		defer close(resultCh)
		defer close(errorCh)

		invokeResultCh, invokeErrCh := a.Invoke(ctx, msg.Message)

		select {
		case result := <-invokeResultCh:
			resultBytes, err := json.Marshal(result)
			if err != nil {
				errorCh <- err
				return
			}
			resultCh <- NewStateMessage(types.RoleAssistant, string(resultBytes))
		case err := <-invokeErrCh:
			errorCh <- err
		case <-ctx.Done():
			errorCh <- ctx.Err()
		}
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