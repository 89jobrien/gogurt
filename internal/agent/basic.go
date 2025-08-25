package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"gogurt/internal/state"
	"gogurt/internal/tools"
	"gogurt/internal/types"
	"strings"
)

// basic is the concrete implementation of the Agent interface.
type basic struct {
	llm           types.LLM
	tools         []*tools.Tool
	memory        []types.ChatMessage
	MaxIterations int
}

var _ Agent = (*basic)(nil)

func New(llm types.LLM, maxIterations int, tools ...*tools.Tool) Agent {
	return &basic{
		llm:           llm,
		tools:         tools,
		MaxIterations: maxIterations,
	}
}

// Run drives the agent loop: ask the LLM, interpret TOOL_CALL responses,
// call tools as needed, and return the final assistant answer.
func (a *basic) Invoke(ctx context.Context, prompt string) (*AgentCallResult, error) {
	a.memory = append(a.memory, types.ChatMessage{Role: types.RoleUser, Content: prompt})

	for i := 0; i < a.MaxIterations; i++ {
		response, err := a.llm.Generate(ctx, a.memory)
		if err != nil {
			return nil, fmt.Errorf("failed to generate response from LLM: %w", err)
		}

		a.memory = append(a.memory, *response)

		if after, ok := strings.CutPrefix(response.Content, "TOOL_CALL:"); ok {
			toolCall := after
			var toolData map[string]string
			if err := json.Unmarshal([]byte(toolCall), &toolData); err != nil {
				return nil, fmt.Errorf("invalid tool call format: %w", err)
			}
			toolName, args := toolData["name"], toolData["arguments"]
			tool, err := a.findTool(toolName)
			if err != nil {
				return nil, err
			}
			result, err := tool.Call(args)
			if err != nil {
				return nil, fmt.Errorf("error calling tool %s: %w", toolName, err)
			}
			resultJSON, err := json.Marshal(result)
			if err != nil {
				return nil, fmt.Errorf("error marshalling tool result: %w", err)
			}
			a.memory = append(a.memory, types.ChatMessage{Role: types.RoleSystem, Content: string(resultJSON)})
		} else {
			return &AgentCallResult{
				Output:   response.Content,
				Metadata: map[string]interface{}{"agent": "basic"},
			}, nil
		}
	}

	return nil, fmt.Errorf("agent reached iteration limit")
}

func (a *basic) findTool(name string) (*tools.Tool, error) {
	for _, tool := range a.tools {
		if tool.Name == name {
			return tool, nil
		}
	}
	return nil, fmt.Errorf("tool %s not found", name)
}

func (b *basic) Delegate(ctx context.Context, input string, agents []Agent) ([]*AgentCallResult, error) {
    // Minimal empty implementation for interface satisfaction
    return nil, nil
}

func (b *basic) InvokeAsync(ctx context.Context, input string) (<-chan AgentCallResult, error) {
	// Provide suitable logic, or a minimal stub if not used:
	return nil, nil
}

func (b *basic) State() state.AgentState {
	// Provide suitable logic, or a minimal stub if not used:
	return nil
}

func (b *basic) Planner() Planner {
	// Provide suitable logic, or a minimal stub if not used:
	return nil
}

func (b *basic) Init(ctx context.Context, params map[string]interface{}) error {
	// Provide suitable logic, or a minimal stub if not used:
	return nil
}