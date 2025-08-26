package llm

import (
	"context"
	"gogurt/internal/types"
)

// type CompletionRequest struct {
// 	Prompt string
// 	// add additional fields as needed
// }

// type CompletionResponse struct {
// 	Text string
// 	// add additional fields as needed
// }

// // LLM is the interface for a large language model
// type LLM interface {
// 	Generate(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)
// 	AGenerate(ctx context.Context, req *CompletionRequest) (<-chan *CompletionResponse, <-chan error)
// 	Stream(ctx context.Context, req *CompletionRequest) (<-chan string, <-chan error)
// }

type LLM interface {
	Generate(ctx context.Context, messages []types.ChatMessage) (*types.ChatMessage, error)
	AGenerate(ctx context.Context, messages []types.ChatMessage) (<-chan *types.ChatMessage, <-chan error)
	Stream(ctx context.Context, messages []types.ChatMessage, onToken func(token string) error) (*types.ChatMessage, error)
	AStream(ctx context.Context, messages []types.ChatMessage) (<-chan string, <-chan error)
	HealthCheck(ctx context.Context) error
	Metadata() map[string]any
}
