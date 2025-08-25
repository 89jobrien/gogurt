package pipes

import (
	"context"
	"fmt"
	"gogurt/internal/agent"
	"gogurt/internal/tools"
)

// ToolPipe invokes agent, then uses tool registry
type ToolPipe struct {
	Agent    agent.Agent
	Registry *tools.Registry
}

func (tp *ToolPipe) Run(ctx context.Context, prompt string) (string, error) {
	result, err := tp.Agent.Invoke(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("agent error: %w", err)
	}

	var toolName, jsonArgs string
	n, err := fmt.Sscanf(result.Output, "%s|%s", &toolName, &jsonArgs)
	if err != nil || n != 2 {
		return "", fmt.Errorf("agent output format invalid: %v", result.Output)
	}

	fmt.Printf("Using tool: %v\v", toolName)
	toolResult, err := tp.Registry.Call(toolName, jsonArgs)
	if err != nil {
		return "", fmt.Errorf("tool error: %w\v", err)
	}
	fmt.Printf("Tool result: %v\n", toolResult)

	return fmt.Sprintf("%v", toolResult), nil
}

// // in main code
// // Setup registry
// reg := tools.NewRegistry()
// reg.Register(addTool)
// reg.Register(subtractTool)

// // Setup agent (needs to emit e.g. "add|{\"a\":1,\"b\":2}")
// myAgent := &YourAgentImplementation{...}

// // Create pipe
// pipe := pipes.ToolPipe{
//     Agent: myAgent,
//     Registry: reg,
// }

// // Run with a prompt
// out, err := pipe.Run(context.Background(), "Add 1 and 2")
// fmt.Println("Final output:", out)
