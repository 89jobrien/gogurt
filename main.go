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
	"gogurt/documentloaders"
	"gogurt/embeddings"
	embollama "gogurt/embeddings/ollama"
	"gogurt/llm/azure"
	llmollama "gogurt/llm/ollama"
	"gogurt/llm/openai"
	"gogurt/splitters"
	"gogurt/splitters/character"
	"gogurt/splitters/markdown"
	"gogurt/splitters/recursive"
	"gogurt/types"
	"gogurt/vectorstores"
	"gogurt/vectorstores/chroma"
	"gogurt/vectorstores/simple"
)

// llm factory
func getLLM(cfg *config.Config) types.LLM {
	var llm types.LLM
	var err error

	switch cfg.LLMProvider {
	case "azure":
		slog.Info("Using AzureOpenAI for LLM")
		llm, err = azure.New(cfg)
	case "ollama":
		slog.Info("Using Ollama for LLM")
		llm, err = llmollama.New(cfg)
	default:
		slog.Info("Using OpenAI	for LLM")
		llm, err = openai.New(cfg)
	}

	if err != nil {
		slog.Error("failed to create LLM", "error", err)
		os.Exit(1)
	}
	return llm
}

// text splitter factory
func getSplitter(cfg *config.Config) splitters.Splitter {
	switch cfg.SplitterProvider {
	case "character":
		slog.Info("Using character text splitter")
		return character.New(100, 20)
	case "markdown":
		slog.Info("Using markdown text splitter")
		return markdown.New(512, 50)
	default:
		slog.Info("Using recursive text splitter")
		return recursive.New(512, 50)
	}
}

// vector store factory
func getVectorStore(cfg *config.Config, embedder embeddings.Embedder) vectorstores.VectorStore {
	var store vectorstores.VectorStore
	var err error

	switch cfg.VectorStoreProvider {
	case "chroma":
		slog.Info("Using Chroma vector store")
		store, err = chroma.New(context.Background(), cfg.ChromaURL, embedder)
	default:
		slog.Info("Using simple in-memory vector store")
		store = simple.New(embedder)
	}

	if err != nil {
		slog.Error("failed to create vector store", "error", err)
		os.Exit(1)
	}
	return store
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// --- Determine document path ---
	var documentPath string
	if len(os.Args) >= 2 {
		documentPath = os.Args[1]
		slog.Info("Loading document from command-line argument", "path", documentPath)
	} else {
		documentPath = "docs/"
		slog.Info("No document path provided, loading from default 'docs/' directory")
	}
	// ---

	cfg := config.Load()
	llmClient := getLLM(cfg)

	// --- RAG pipeline setup ---
	slog.Info("Setting up RAG pipeline...")
	docs, err := documentloaders.LoadDocuments(documentPath)
	if err != nil {
		slog.Error("failed to load documents", "path", documentPath, "error", err)
		os.Exit(1)
	}
	if len(docs) == 0 {
		slog.Warn("No documents were loaded.")
	}

	splitter := getSplitter(cfg)
	chunks := splitter.SplitDocuments(docs)

	embedder, err := embollama.New(cfg)
	if err != nil {
		slog.Error("failed to create embedder", "error", err)
		os.Exit(1)
	}

	vectorStore := getVectorStore(cfg, embedder)

	if len(chunks) > 0 {
		err = vectorStore.AddDocuments(context.Background(), chunks)
		if err != nil {
			slog.Error("failed to add documents to vector store", "error", err)
			os.Exit(1)
		}
	}
	slog.Info("RAG pipeline setup complete.", "documents_loaded", len(docs), "chunks_created", len(chunks))

	aiAgent := agent.New(llmClient, cfg.AgentMaxIterations)

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

		var augmentedPrompt string
		if len(docs) > 0 {
			relevantDocs, err := vectorStore.SimilaritySearch(context.Background(), prompt, 2)
			if err != nil {
				slog.Error("failed to retrieve documents", "error", err)
				continue
			}

			var contextBuilder strings.Builder
			for _, doc := range relevantDocs {
				contextBuilder.WriteString(doc.PageContent + "\n")
			}

			augmentedPrompt = fmt.Sprintf(`
			Answer the following question based on this context:
			---
			Context:
			%s
			---
			Question: %s`, contextBuilder.String(), prompt)
		} else {
			augmentedPrompt = prompt
		}

		response, err := aiAgent.Run(context.Background(), augmentedPrompt)
		if err != nil {
			slog.Error("agent run failed", "error", err)
			continue
		}

		os.Stdout.WriteString("AI: " + response + "\n")
	}
}