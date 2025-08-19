package main

import (
	"bufio"
	"context"
	"log/slog"
	"os"
	"strings"

	"gogurt/agent"
	"gogurt/config"
	"gogurt/llm/azure"
	"gogurt/llm/ollama"
	"gogurt/llm/openai"
	"gogurt/types"
)

func getLLM(cfg *config.Config) types.LLM {
	var llm types.LLM
	var err error

	switch cfg.LLMProvider {
	case "azure":
		slog.Info("Using Azure LLM")
		llm, err = azure.New(cfg)
	case "ollama":
		slog.Info("Using Ollama LLM")
		llm, err = ollama.New(cfg)
	default:
		slog.Info("Using OpenAI LLM")
		llm, err = openai.New(cfg)
	}

	if err != nil {
		slog.Error("failed to create LLM", "error", err)
		os.Exit(1)
	}
	return llm
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()
	llmClient := getLLM(cfg)

	aiAgent := agent.New(llmClient)

	slog.Info("Chat session started. Type 'exit' to end.")
	reader := bufio.NewReader(os.Stdin)

	for {
		os.Stdout.WriteString("You: ")
		prompt, _ := reader.ReadString('\n')
		prompt = strings.TrimSpace(prompt)

		if strings.ToLower(prompt) == "exit" {
			slog.Info("Ending chat session.")
			break
		}

		response, err := aiAgent.Run(context.Background(), prompt)
		if err != nil {
			slog.Error("agent run failed", "error", err)
			continue
		}

		os.Stdout.WriteString("AI: " + response + "\n")
	}
}