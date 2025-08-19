package openai

import (
	"context"
	"os"

	"gogurt/types"

	"github.com/sashabaranov/go-openai"
)

// OpenAI is a client for the OpenAI API.
type OpenAI struct {
	client *openai.Client
}

// New creates a new OpenAI client.
// The return type is the interface, which is correct.
func New() types.LLM {
    return &OpenAI{
        client: openai.NewClient(os.Getenv("OPENAI_API_KEY")),
    }
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

	resp, err := o.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    openai.GPT4o, // Using a more recent model
			Messages: apiMessages,
		},
	)

	if err != nil {
		return nil, err
	}

	responseMessage := resp.Choices[0].Message
	return &types.ChatMessage{
		Role:    types.Role(responseMessage.Role),
		Content: responseMessage.Content,
	}, nil
}