package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"gogurt/agent"
	"gogurt/config"
	"gogurt/documentloaders/text"
	"gogurt/embeddings/ollama"
	"gogurt/llm/azure"
	llmollama "gogurt/llm/ollama"
	"gogurt/llm/openai"
	"gogurt/textsplitter/character"
	"gogurt/tools"
	"gogurt/types"
	"gogurt/vectorstores/simple"
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
		llm, err = llmollama.New(cfg)
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

	// --- RAG Pipeline Setup ---
	slog.Info("Setting up RAG pipeline...")
	docs, err := text.NewTextLoader("docs.txt") // Replace with the path to your document
	if err != nil {
		slog.Error("failed to load documents", "error", err)
		os.Exit(1)
	}

	splitter := character.New(100, 20)
	chunks := splitter.SplitDocuments(docs)

	embedder, err := ollama.New(cfg)
	if err != nil {
		slog.Error("failed to create embedder", "error", err)
		os.Exit(1)
	}

	vectorStore := simple.New(embedder)
	err = vectorStore.AddDocuments(context.Background(), chunks)
	if err != nil {
		slog.Error("failed to add documents to vector store", "error", err)
		os.Exit(1)
	}
	slog.Info("RAG pipeline setup complete.")
	// --- End of RAG Setup ---

	weatherTool, err := tools.New(tools.GetWeather, "Get the weather for a city")
	if err != nil {
		slog.Error("failed to create weather tool", "error", err)
		os.Exit(1)
	}

	aiAgent := agent.New(llmClient, weatherTool)

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

		// --- Augment the Prompt using the Retriever ---
		relevantDocs, err := vectorStore.SimilaritySearch(context.Background(), prompt, 2)
		if err != nil {
			slog.Error("failed to retrieve documents", "error", err)
			continue
		}

		var contextBuilder strings.Builder
		for _, doc := range relevantDocs {
			contextBuilder.WriteString(doc.PageContent + "\n")
		}

		augmentedPrompt := fmt.Sprintf(`
		Answer the following question based on this context:
		---
		Context:
		%s
		---
		Question: %s`, contextBuilder.String(), prompt)
		// --- End of Augmentation ---

		response, err := aiAgent.Run(context.Background(), augmentedPrompt)
		if err != nil {
			slog.Error("agent run failed", "error", err)
			continue
		}

		os.Stdout.WriteString("AI: " + response + "\n")
	}
}