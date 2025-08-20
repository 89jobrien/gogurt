package agent

import (
	"context"
	"gogurt/types"
	"testing"
)

type mockLLM struct {
	responses []string
	callCount int
}

func (m *mockLLM) Generate(ctx context.Context, messages []types.ChatMessage) (*types.ChatMessage, error) {
	response := m.responses[m.callCount]
	m.callCount++
	return &types.ChatMessage{Role: types.RoleAssistant, Content: response}, nil
}

func TestAgent_Run(t *testing.T) {

	llm := &mockLLM{
		responses: []string{
			"The weather is sunny.",
		},
	}

	agent := New(llm, 5)
	result, err := agent.Run(context.Background(), "What is the weather in New York?")
	if err != nil {
		t.Fatalf("Agent.Run() error = %v", err)
	}

	expectedResult := "The weather is sunny."
	if result != expectedResult {
		t.Errorf("Agent.Run() = %v, want %v", result, expectedResult)
	}
}