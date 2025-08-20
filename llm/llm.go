package llm

import (
	"context"

	"gogurt/types"
)

// LLM is the interface that all model adapters should implement.
//
// Generate produces a single final assistant message synchronously.
// Stream produces output incrementally by invoking the provided onToken
// callback for each streamed chunk (typically tokens). The Stream method
// returns the final assembled assistant message (or an error).
type LLM interface {
	Generate(ctx context.Context, messages []types.ChatMessage) (*types.ChatMessage, error)

	// Stream streams the model output. onToken is called for each chunk
	// of streamed content (e.g., tokens or partial text). If onToken
	// returns an error, streaming should stop and that error should be returned.
	// The method returns the final assembled assistant message and any error.
	Stream(ctx context.Context, messages []types.ChatMessage, onToken func(token string) error) (*types.ChatMessage, error)
}

// Example implementation of the LLM interface.
type MyLLM struct{}

// compile-time check that MyLLM implements LLM.
var _ LLM = (*MyLLM)(nil)

// NewMyLLM returns a new example LLM implementation.
func NewMyLLM() LLM {
	return &MyLLM{}
}

func (m *MyLLM) Generate(ctx context.Context, messages []types.ChatMessage) (*types.ChatMessage, error) {
	return &types.ChatMessage{Role: types.RoleAssistant, Content: "response"}, nil
}

// Stream example: streams the word "response" one rune at a time.
// Calls onToken for each chunk and returns the final assembled message.
func (m *MyLLM) Stream(ctx context.Context, messages []types.ChatMessage, onToken func(token string) error) (*types.ChatMessage, error) {
	content := "response"
	for _, r := range content {
		// If context cancelled, stop early.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if err := onToken(string(r)); err != nil {
			return nil, err
		}
	}
	return &types.ChatMessage{Role: types.RoleAssistant, Content: content}, nil
}