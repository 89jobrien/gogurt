package stateful

import (
	"fmt"
	"gogurt/internal/state"
	"gogurt/internal/tools"
	"reflect"
)

// --- SaveToScratchpad Tool ---

type SaveToScratchpadArgs struct {
	Key     string `json:"key"`
	Content string `json:"content"`
}

func SaveToScratchpad(args SaveToScratchpadArgs, agentState state.AgentState) (string, error) {
	if err := agentState.Set(args.Key, args.Content); err != nil {
		return "", fmt.Errorf("failed to save to scratchpad: %w", err)
	}
	return fmt.Sprintf("Successfully saved to scratchpad with key: %s", args.Key), nil
}

var SaveToScratchpadTool = &tools.Tool{
	Name:        "save_to_scratchpad",
	Description: "Saves content to the agent's in-memory scratchpad.",
	Func:        reflect.ValueOf(SaveToScratchpad),
	InputSchema: tools.GenInputSchema(reflect.TypeOf(SaveToScratchpadArgs{})),
	Example:     `{"key":"search_results", "content":"..."}`,
	Metadata:    map[string]any{"category": "state"},
}

// --- ReadScratchpad Tool ---

type ReadScratchpadArgs struct {
	Key string `json:"key"`
}

func ReadScratchpad(args ReadScratchpadArgs, agentState state.AgentState) (string, error) {
	val, err := agentState.Get(args.Key)
	if err != nil {
		return "", fmt.Errorf("failed to read from scratchpad: %w", err)
	}
	if val == nil {
		return "", fmt.Errorf("no value found in scratchpad for key: %s", args.Key)
	}
	return fmt.Sprintf("%v", val), nil
}

var ReadScratchpadTool = &tools.Tool{
	Name:        "read_scratchpad",
	Description: "Reads content from the agent's in-memory scratchpad.",
	Func:        reflect.ValueOf(ReadScratchpad),
	InputSchema: tools.GenInputSchema(reflect.TypeOf(ReadScratchpadArgs{})),
	Example:     `{"key":"search_results"}`,
	Metadata:    map[string]any{"category": "state"},
}