package azure

import (
	"context"
	"gogurt/internal/config"
	"gogurt/internal/llm"
	"gogurt/internal/types"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type AzureLLM struct {
	client         *openai.Client
	deploymentName string
}

// HealthCheck implements types.LLM.
func (a *AzureLLM) HealthCheck(ctx context.Context) error {
	_, err := a.Generate(ctx, []types.ChatMessage{{Role: "system", Content: "ping"}})
	return err
}

// Metadata implements types.LLM.
func (a *AzureLLM) Metadata() map[string]any {
	md := make(map[string]any)
	md["model"] = a.deploymentName
	return md
}

func New(cfg *config.Config) (llm.LLM, error) {
	client := openai.NewClient(cfg.AzureOpenAIAPIKey)
	return &AzureLLM{
		client:         client,
		deploymentName: cfg.AzureDeployment,
	}, nil
}

// Convert types.ChatMessage to openai.ChatCompletionMessage
func toOpenAIMessages(messages []types.ChatMessage) []openai.ChatCompletionMessage {
	apiMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		apiMessages[i] = openai.ChatCompletionMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		}
	}
	return apiMessages
}

// Generate generates a response from the Azure API.
func (a *AzureLLM) Generate(ctx context.Context, messages []types.ChatMessage) (*types.ChatMessage, error) {
	apiMessages := toOpenAIMessages(messages)
	res, err := a.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    a.deploymentName,
		Messages: apiMessages,
	})
	if err != nil {
		return nil, err
	}
	var responseContent strings.Builder
	var responseRole types.Role
	responseContent.WriteString(res.Choices[0].Message.Content)
	responseRole = types.Role(res.Choices[0].Message.Role)
	if responseRole == "" {
		responseRole = types.RoleAssistant
	}
	return &types.ChatMessage{
		Role:    responseRole,
		Content: responseContent.String(),
	}, nil
}

// AGenerate provides an asynchronous Generate.
func (a *AzureLLM) AGenerate(ctx context.Context, messages []types.ChatMessage) (<-chan *types.ChatMessage, <-chan error) {
	msgCh := make(chan *types.ChatMessage, 1)
	errCh := make(chan error, 1)
	go func() {
		defer close(msgCh)
		defer close(errCh)
		msg, err := a.Generate(ctx, messages)
		if err != nil {
			errCh <- err
			return
		}
		msgCh <- msg
	}()
	return msgCh, errCh
}

// Stream streams model output from Azure, forwarding chunks to onToken.
func (a *AzureLLM) Stream(ctx context.Context, messages []types.ChatMessage, onToken func(token string) error) (*types.ChatMessage, error) {
	apiMessages := toOpenAIMessages(messages)
	stream := true
	resStream, err := a.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:    a.deploymentName,
		Messages: apiMessages,
		Stream:   stream,
	})
	if err != nil {
		return nil, err
	}
	var responseContent strings.Builder
	var responseRole types.Role
	for {
		resp, err := resStream.Recv()
		if err != nil {
			// io.EOF => end of stream, treat as normal exit
			break
		}
		if len(resp.Choices) == 0 {
			continue
		}
		choice := resp.Choices[0]
		token := choice.Delta.Content
		role := choice.Delta.Role
		if responseRole == "" && role != "" {
			responseRole = types.Role(role)
		}
		if token != "" {
			if err := onToken(token); err != nil {
				return nil, err
			}
			responseContent.WriteString(token)
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

// AStream provides an asynchronous streaming interface returning tokens.
func (a *AzureLLM) AStream(ctx context.Context, messages []types.ChatMessage) (<-chan string, <-chan error) {
	tokenCh := make(chan string)
	errCh := make(chan error, 1)
	go func() {
		defer close(tokenCh)
		defer close(errCh)
		_, err := a.Stream(ctx, messages, func(token string) error {
			tokenCh <- token
			return nil
		})
		if err != nil {
			errCh <- err
		}
	}()
	return tokenCh, errCh
}
