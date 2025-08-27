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

// NewRAGPipe creates a new RAG query pipeline.
func NewRAGPipe(ctx context.Context, cfg *config.Config) (*RAGPipe, error) {
	c.Write("Setting up RAG query pipeline...")
	// Assuming agent.NewAgent returns a compliant agent.
	// The empty AgentConfig might need to be populated depending on the agent type.
	aiAgent, err := agent.NewAgent(types.AgentConfig{Name: "ResearchAgent"}) // Example agent name
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

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

// Run executes a RAG query asynchronously.
func (r *RAGPipe) Run(ctx context.Context, query string) (<-chan string, <-chan error) {
	resultCh := make(chan string, 1)
	errorCh := make(chan error, 1)

	go func() {
		defer close(resultCh)
		defer close(errorCh)

		if strings.TrimSpace(query) == "" {
			errorCh <- fmt.Errorf("query cannot be empty")
			return
		}

		// 1. Retrieve relevant documents asynchronously
		docsCh, docsErrCh := r.vectorStore.SimilaritySearch(ctx, query, 3)
		var relevantDocs []types.Document
		select {
		case relevantDocs = <-docsCh:
		case err := <-docsErrCh:
			errorCh <- fmt.Errorf("failed to retrieve documents: %w", err)
			return
		case <-ctx.Done():
			errorCh <- ctx.Err()
			return
		}

		// 2. Build context from retrieved documents
		var contextBuilder strings.Builder
		for i, doc := range relevantDocs {
			if i > 0 {
				contextBuilder.WriteString("\n---\n")
			}
			contextBuilder.WriteString(doc.PageContent)
		}

		// 3. Format the prompt
		augmentedPrompt, err := r.prompt.Format(map[string]string{
			"context":  contextBuilder.String(),
			"question": query,
		})
		if err != nil {
			errorCh <- fmt.Errorf("failed to format prompt: %w", err)
			return
		}

		// 4. Generate response using the agent asynchronously
		responseCh, agentErrCh := r.Agent.Invoke(ctx, augmentedPrompt)
		select {
		case response := <-responseCh:
			if s, ok := response.(string); ok {
				resultCh <- s
			} else {
				resultCh <- fmt.Sprintf("%v", response)
			}
		case err := <-agentErrCh:
			errorCh <- fmt.Errorf("agent invocation failed: %w", err)
			return
		case <-ctx.Done():
			errorCh <- ctx.Err()
			return
		}
	}()

	return resultCh, errorCh
}

// GetVectorStore returns the vector store instance.
func (r *RAGPipe) GetVectorStore() vectorstores.VectorStore {
	return r.vectorStore
}

// HasDocuments checks if the vector store contains any documents asynchronously.
func (r *RAGPipe) HasDocuments(ctx context.Context) (<-chan bool, <-chan error) {
	resultCh := make(chan bool, 1)
	errorCh := make(chan error, 1)

	go func() {
		defer close(resultCh)
		defer close(errorCh)

		docsCh, errCh := r.vectorStore.SimilaritySearch(ctx, "test", 1)
		select {
		case docs := <-docsCh:
			resultCh <- len(docs) > 0
		case err := <-errCh:
			errorCh <- err
		case <-ctx.Done():
			errorCh <- ctx.Err()
		}
	}()

	return resultCh, errorCh
}