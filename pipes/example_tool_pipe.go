package pipes

import (
	"context"
	"fmt"
	"gogurt/agent"
	"gogurt/tools"
)

// ToolPipe invokes agent, then uses tool registry
type ToolPipe struct {
	Agent    agent.Agent
	Registry *tools.Registry
}

func (tp *ToolPipe) Run(ctx context.Context, prompt string) (string, error) {
	// 1. Invoke agent on prompt. Suppose agent returns tool name and json args.
	result, err := tp.Agent.Invoke(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("agent error: %w", err)
	}

	// For this example, let's suppose AgentCallResult.Output is like: "add|{\"a\":1,\"b\":2}"
	var toolName, jsonArgs string
	n, err := fmt.Sscanf(result.Output, "%s|%s", &toolName, &jsonArgs)
	if err != nil || n != 2 {
		return "", fmt.Errorf("agent output format invalid: %v", result.Output)
	}

	// 2. Call tool via registry
	fmt.Printf("Using tool: %v\v", toolName)
	toolResult, err := tp.Registry.Call(toolName, jsonArgs)
	if err != nil {
		return "", fmt.Errorf("tool error: %w\v", err)
	}
	fmt.Printf("Tool result: %v\n", toolResult)

	return fmt.Sprintf("%v", toolResult), nil
}