package openai

import (
	"context"
	"fmt"
	"os"

	"gogurt/config"
	"gogurt/types"

	"github.com/sashabaranov/go-openai"
)

// OpenAI is a client for the OpenAI API.
type OpenAI struct {
	client *openai.Client
}

// New creates a new OpenAI client.
// The return type is the interface, which is correct.
func New(cfg *config.Config) (types.LLM, error) {
	apiKey := ""
	if cfg != nil {
		apiKey = cfg.OpenAIAPIKey
	}
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("openai api key not provided (config.OpenAIAPIKey or OPENAI_API_KEY)")
	}

	client := openai.NewClient(apiKey)
	return &OpenAI{client: client}, nil
}

// Generate generates a response from the OpenAI API.
func (o *OpenAI) Generate(ctx context.Context, messages []types.ChatMessage) (*types.ChatMessage, error) {
	apiMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		apiMessages[i] = openai.ChatCompletionMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		}
	}

	resp, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT4o,
		Messages: apiMessages,
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned from OpenAI")
	}

	responseMessage := resp.Choices[0].Message
	return &types.ChatMessage{
		Role:    types.Role(responseMessage.Role),
		Content: responseMessage.Content,
	}, nil
}