package azure

import (
	"context"
	"fmt"
	"io"
	"strings"

	"gogurt/internal/config"
	"gogurt/internal/types"

	"github.com/sashabaranov/go-openai"
)

type Azure struct {
	client         *openai.Client
	deploymentName string
}

func New(cfg *config.Config) (types.LLM, error) {
	// Validate that the required configuration is present.
	if cfg.AzureOpenAIAPIKey == "" {
		return nil, fmt.Errorf("azure api key not provided")
	}
	if cfg.AzureOpenAIEndpoint == "" {
		return nil, fmt.Errorf("azure endpoint not provided")
	}
	if cfg.AzureDeployment == "" {
		return nil, fmt.Errorf("azure deployment name not provided")
	}

	// Use the specific Azure configuration from the go-openai library.
	config := openai.DefaultAzureConfig(cfg.AzureOpenAIAPIKey, cfg.AzureOpenAIEndpoint)
	client := openai.NewClientWithConfig(config)

	return &Azure{
		client:         client,
		deploymentName: cfg.AzureDeployment,
	}, nil
}

func (a *Azure) Generate(ctx context.Context, messages []types.ChatMessage) (*types.ChatMessage, error) {
	apiMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		apiMessages[i] = openai.ChatCompletionMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		}
	}

	resp, err := a.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    a.deploymentName,
			Messages: apiMessages,
		},
	)

	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned from Azure OpenAI")
	}

	responseMessage := resp.Choices[0].Message
	return &types.ChatMessage{
		Role:    types.Role(responseMessage.Role),
		Content: responseMessage.Content,
	}, nil
}

// Stream streams the model output. onToken is called for each streamed chunk.
func (a *Azure) Stream(ctx context.Context, messages []types.ChatMessage, onToken func(token string) error) (*types.ChatMessage, error) {
	apiMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		apiMessages[i] = openai.ChatCompletionMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		}
	}

	req := openai.ChatCompletionRequest{
		Model:    a.deploymentName,
		Messages: apiMessages,
		Stream:   true,
	}

	stream, err := a.client.CreateChatCompletionStream(ctx, req)
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
				// Forward token to callback
				if err := onToken(token); err != nil {
					// stop streaming if callback requests it
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
