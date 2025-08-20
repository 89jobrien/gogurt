package factories

import (
	"gogurt/config"
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
	"log/slog"
	"os"
)

// llm factory
func GetLLM(cfg *config.Config) types.LLM {
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
func GetSplitter(cfg *config.Config) splitters.Splitter {
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
func GetVectorStore(cfg *config.Config, embedder embeddings.Embedder) vectorstores.VectorStore {
	var store vectorstores.VectorStore
	var err error

	switch cfg.VectorStoreProvider {
	case "chroma":
		slog.Info("Using Chroma vector store")
		store, err = chroma.New(cfg)
	default:
		slog.Info("Using in-memory vector store")
		store = simple.New(embedder)
	}

	if err != nil {
		slog.Error("failed to create vector store", "error", err)
		os.Exit(1)
	}
	return store
}

// embedder factory
func GetEmbedder(cfg *config.Config) embeddings.Embedder {
	embedder, err := embollama.New(cfg)
	if err != nil {
		slog.Error("failed to create embedder", "error", err)
		os.Exit(1)
	}
	return embedder
}