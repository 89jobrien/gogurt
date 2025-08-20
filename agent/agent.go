package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"gogurt/tools"
	"gogurt/types"
	"strings"
)

type Agent struct {
    llm           types.LLM
    tools         []*tools.Tool
    memory        []types.ChatMessage
    MaxIterations int
}

func New(llm types.LLM, maxIterations int, tools ...*tools.Tool) *Agent {
    return &Agent{
        llm:           llm,
        tools:         tools,
        MaxIterations: maxIterations,
    }
}

func (a *Agent) Run(ctx context.Context, prompt string) (string, error) {
    a.memory = append(a.memory, types.ChatMessage{Role: types.RoleUser, Content: prompt})
    // max iterations is used to prevent infinite loops
    for i := 0; i < a.MaxIterations; i++ {
        response, err := a.llm.Generate(ctx, a.memory)
        if err != nil {
            return "", fmt.Errorf("failed to generate response from LLM: %w", err)
        }
        a.memory = append(a.memory, *response)
        if after, ok := strings.CutPrefix(response.Content, "TOOL_CALL:"); ok {
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
            resultJSON, err := json.Marshal(result)
            if err != nil {
                return "", fmt.Errorf("error marshalling tool result: %w", err)
            }
            a.memory = append(a.memory, types.ChatMessage{Role: types.RoleSystem, Content: string(resultJSON)})
        } else {
            return response.Content, nil
        }
    }
    return "", fmt.Errorf("agent reached iteration limit")
}

func (a *Agent) findTool(name string) (*tools.Tool, error) {
    for _, tool := range a.tools {
        if tool.Name == name {
            return tool, nil
        }
    }
    return nil, fmt.Errorf("tool %s not found", name)
}
