package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gogurt/tools"
	"gogurt/types"
)

type Agent struct {
    llm    types.LLM
    tools  []*tools.Tool
    memory []types.ChatMessage
}

func New(llm types.LLM, tools ...*tools.Tool) *Agent {
	return &Agent{
		llm:   llm,
		tools: tools,
	}
}

// Run executes the agent with a given prompt.
func (a *Agent) Run(ctx context.Context, prompt string) (string, error) {
	a.memory = append(a.memory, types.ChatMessage{Role: types.RoleUser, Content: prompt})

	for {
		response, err := a.llm.Generate(ctx, a.memory)
		if err != nil {
			return "", err
		}

		a.memory = append(a.memory, *response)

		// Check if the response is a tool call
		if after, ok :=strings.CutPrefix(response.Content, "TOOL_CALL:"); ok  {
			toolCall := after
			var toolData map[string]string
			if err := json.Unmarshal([]byte(toolCall), &toolData); err != nil {
				return "", fmt.Errorf("invalid tool call format: %w", err)
			}

			toolName, args := toolData["name"], toolData["arguments"]
			tool, err := a.findTool(toolName)
			if err != nil {
				return "", err
			}

			result, err := tool.Call(args)
			if err != nil {
				return "", fmt.Errorf("error calling tool %s: %w", toolName, err)
			}

			resultJSON, _ := json.Marshal(result)
			a.memory = append(a.memory, types.ChatMessage{Role: types.RoleSystem, Content: string(resultJSON)})
		} else {
			return response.Content, nil
		}
	}
}

func (a *Agent) findTool(name string) (*tools.Tool, error) {
	for _, tool := range a.tools {
		if tool.Name == name {
			return tool, nil
		}
	}
	return nil, fmt.Errorf("tool %s not found", name)
}