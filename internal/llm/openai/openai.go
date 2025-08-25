package openai

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"gogurt/internal/config"
	"gogurt/internal/llm"
	"gogurt/internal/types"

	"github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	client *openai.Client
}

func New(cfg *config.Config) (llm.LLM, error) {
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

// generates a response from the OpenAI API.
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

// Stream streams the model output. onToken is called for each streamed chunk.
func (o *OpenAI) Stream(ctx context.Context, messages []types.ChatMessage, onToken func(token string) error) (*types.ChatMessage, error) {
	apiMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		apiMessages[i] = openai.ChatCompletionMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		}
	}

	req := openai.ChatCompletionRequest{
		Model:    openai.GPT4o,
		Messages: apiMessages,
	}

	stream, err := o.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	var responseContent strings.Builder
	var responseRole types.Role

	for {
		part, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for _, choice := range part.Choices {
			if responseRole == "" && choice.Delta.Role != "" {
				responseRole = types.Role(choice.Delta.Role)
			}
			if choice.Delta.Content != "" {
				token := choice.Delta.Content
				if err := onToken(token); err != nil {
					return nil, err
				}
				responseContent.WriteString(token)
			}
		}
	}

	if responseRole == "" {
		responseRole = types.RoleAssistant
	}

	return &types.ChatMessage{
		Role:    responseRole,
		Content: responseContent.String(),
	}, nil
}

func (o *OpenAI) HealthCheck(ctx context.Context) error {
	// Simple completion with small dummy prompt, expect nil error on healthy.
	_, err := o.Generate(ctx, []types.ChatMessage{{Role: "system", Content: "ping"}})
	return err
}
func (o *OpenAI) Metadata() map[string]interface{} {
	return map[string]interface{}{
		"provider": "OpenAI",
		"version":  "latest",
	}
}
