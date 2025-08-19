package azure

import (
	"context"
	"fmt"
	"os"

	"gogurt/config"
	"gogurt/types"

	"github.com/sashabaranov/go-openai"
)

type Azure struct {
	client         *openai.Client
	deploymentName string
}

func New(cfg *config.Config) (types.LLM, error) {
    key := ""
    if cfg != nil {
        key = cfg.AzureOpenAIAPIKey
    }
    if key == "" {
        key = os.Getenv("AZURE_KEY")
    }
    if key == "" {
        return nil, fmt.Errorf("azure key not provided")
    }
	
	client := openai.NewClient(cfg.AzureOpenAIAPIKey)
	deploymentName := cfg.AzureDeployment

	return &Azure{
		client:         client,
		deploymentName: deploymentName,
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

	responseMessage := resp.Choices[0].Message
	return &types.ChatMessage{
		Role:    types.Role(responseMessage.Role),
		Content: responseMessage.Content,
	}, nil
}