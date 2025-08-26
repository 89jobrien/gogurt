package factories

import (
	"context"
	"gogurt/internal/config"
	"gogurt/internal/embeddings"
	embollama "gogurt/internal/embeddings/ollama"
	"gogurt/internal/llm"
	"gogurt/internal/llm/azure"
	llmollama "gogurt/internal/llm/ollama"
	"gogurt/internal/llm/openai"
	"gogurt/internal/splitters"
	"gogurt/internal/splitters/character"
	"gogurt/internal/splitters/markdown"
	"gogurt/internal/splitters/recursive"
	"gogurt/internal/vectorstores"
	"gogurt/internal/vectorstores/chroma"
	"gogurt/internal/vectorstores/simple"
	"log/slog"
	"os"
)

// llm factory (synchronous)
func GetLLM(cfg *config.Config) llm.LLM {
	var llmModel llm.LLM
	var err error
	switch cfg.LLMProvider {
	case "azure":
		slog.Info("Using AzureOpenAI for LLM")
		llmModel, err = azure.New(cfg)
	case "ollama":
		slog.Info("Using Ollama for LLM")
		llmModel, err = llmollama.New(cfg)
	default:
		slog.Info("Using OpenAI for LLM")
		llmModel, err = openai.New(cfg)
	}
	if err != nil {
		slog.Error("failed to create LLM", "error", err)
		os.Exit(1)
	}
	return llmModel
}

// llm factory (async)
func AGetLLM(ctx context.Context, cfg *config.Config) (<-chan llm.LLM, <-chan error) {
	out := make(chan llm.LLM, 1)
	errCh := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errCh)
		model := GetLLM(cfg)
		out <- model
	}()
	return out, errCh
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

// async splitter factory
func AGetSplitter(ctx context.Context, cfg *config.Config) <-chan splitters.Splitter {
	out := make(chan splitters.Splitter, 1)
	go func() {
		defer close(out)
		splitter := GetSplitter(cfg)
		out <- splitter
	}()
	return out
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

// async vector store factory
func AGetVectorStore(ctx context.Context, cfg *config.Config, embedder embeddings.Embedder) (<-chan vectorstores.VectorStore, <-chan error) {
	out := make(chan vectorstores.VectorStore, 1)
	errCh := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errCh)
		store := GetVectorStore(cfg, embedder)
		out <- store
	}()
	return out, errCh
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

// async embedder factory
func AGetEmbedder(ctx context.Context, cfg *config.Config) (<-chan embeddings.Embedder, <-chan error) {
	out := make(chan embeddings.Embedder, 1)
	errCh := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errCh)
		embedder := GetEmbedder(cfg)
		out <- embedder
	}()
	return out, errCh
}
