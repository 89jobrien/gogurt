package azure

import (
	"context"
	"os"

	"gogurt/types"

	"github.com/sashabaranov/go-openai"
)

type Azure struct {
	client         *openai.Client
	deploymentName string
}

func New() types.LLM {
	endpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
	apiKey := os.Getenv("AZURE_OPENAI_API_KEY")
	deploymentName := os.Getenv("AZURE_OPENAI_DEPLOYMENT_NAME")

	config := openai.DefaultAzureConfig(apiKey, endpoint)

	client := openai.NewClientWithConfig(config)

	return &Azure{
		client:         client,
		deploymentName: deploymentName,
	}
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