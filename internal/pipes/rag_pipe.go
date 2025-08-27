package pipes

import (
	"context"
	"fmt"
	"gogurt/internal/agent"
	"gogurt/internal/config"
	"gogurt/internal/factories"
	"gogurt/internal/prompts"
	"gogurt/internal/prompts/rag"
	"gogurt/internal/types"
	"gogurt/internal/vectorstores"
	"strings"
)

type RAGPipe struct {
	Agent       agent.Agent
	prompt      *prompts.PromptTemplate
	vectorStore vectorstores.VectorStore
}

// NewRAGPipe creates a new RAG query pipeline (assumes documents are already ingested)
func NewRAGPipe(ctx context.Context, cfg *config.Config) (*RAGPipe, error) {
	c.Write("Setting up RAG query pipeline...")
	aiAgent, _ := agent.NewAgent(types.AgentConfig{})
	embedder := factories.GetEmbedder(cfg)
	vectorStore := factories.GetVectorStore(cfg, embedder)
	ragPrompt, err := prompts.NewPromptTemplate(rag.BasicRagPrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to create prompt template: %w", err)
	}

	c.Write("RAG query pipeline setup complete")

	return &RAGPipe{
		Agent:       aiAgent,
		prompt:      ragPrompt,
		vectorStore: vectorStore,
	}, nil
}

// Run executes a RAG query
func (r *RAGPipe) Run(ctx context.Context, query string) (string, error) {
	if query = strings.TrimSpace(query); query == "" {
		return "", fmt.Errorf("query cannot be empty")
	}

	// Retrieve relevant documents
	relevantDocs, err := r.vectorStore.SimilaritySearch(ctx, query, 3)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve documents: %w", err)
	}

	// Build context from retrieved documents
	var contextBuilder strings.Builder
	for i, doc := range relevantDocs {
		if i > 0 {
			contextBuilder.WriteString("\n---\n")
		}
		contextBuilder.WriteString(doc.PageContent)
	}

	// Format the prompt with context and question
	augmentedPrompt, err := r.prompt.Format(map[string]string{
		"context":  contextBuilder.String(),
		"question": query,
	})
	if err != nil {
		return "", fmt.Errorf("failed to format prompt: %w", err)
	}

	// Generate response using the agent
	response, err := r.Agent.Invoke(ctx, augmentedPrompt)
	if err != nil {
		return "", fmt.Errorf("agent invocation failed: %w", err)
	}

	// FIX: Use response string, not response.Output
	if s, ok := response.(string); ok {
		return s, nil
	}
	return fmt.Sprintf("%v", response), nil
}

// GetVectorStore returns the vector store instance (useful for metrics)
func (r *RAGPipe) GetVectorStore() vectorstores.VectorStore {
	return r.vectorStore
}

// HasDocuments checks if the vector store contains any documents
func (r *RAGPipe) HasDocuments(ctx context.Context) (bool, error) {
	docs, err := r.vectorStore.SimilaritySearch(ctx, "test", 1)
	if err != nil {
		return false, err
	}
	return len(docs) > 0, nil
}
