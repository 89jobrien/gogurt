// In llm/azure/azure.go
package azure

import (
	"context"
	"fmt"

	"gogurt/config"
	"gogurt/types"

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