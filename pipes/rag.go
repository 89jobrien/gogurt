package pipes

import (
	"context"
	"fmt"
	"gogurt/agent"
	"gogurt/config"
	"gogurt/documentloaders"
	"gogurt/factories"
	"gogurt/prompts"
	"gogurt/vectorstores"
	"log/slog"
	"strings"
)

type RAGPipe struct {
	vectorStore  vectorstores.VectorStore
	agent        *agent.Agent
	prompt       *prompts.PromptTemplate
	hasDocuments bool
}

func NewRAG(ctx context.Context, cfg *config.Config, documentPath string) (*RAGPipe, error) {
	slog.Info("Setting up RAG pipeline...")

	llmClient := factories.GetLLM(cfg)
	aiAgent := agent.New(llmClient, cfg.AgentMaxIterations)

	ragPrompt, err := prompts.NewPromptTemplate(prompts.RagPrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to create prompt template: %w", err)
	}

	docs, err := documentloaders.LoadDocuments(documentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load documents from %s: %w", documentPath, err)
	}
	if len(docs) == 0 {
		slog.Warn("No documents were loaded. The agent will run without retrieval context.")
		return &RAGPipe{
			agent:        aiAgent,
			prompt:       ragPrompt,
			hasDocuments: false,
		}, nil
	}

	splitter := factories.GetSplitter(cfg)
	chunks := splitter.SplitDocuments(docs)
	embedder := factories.GetEmbedder(cfg)
	vectorStore := factories.GetVectorStore(cfg, embedder)

	err = vectorStore.AddDocuments(ctx, chunks)
	if err != nil {
		return nil, fmt.Errorf("failed to add documents to vector store: %w", err)
	}

	slog.Info("RAG pipeline setup complete.", "documents_loaded", len(docs), "chunks_created", len(chunks))

	return &RAGPipe{
		vectorStore:  vectorStore,
		agent:        aiAgent,
		prompt:       ragPrompt,
		hasDocuments: true,
	}, nil
}

func (p *RAGPipe) Run(ctx context.Context, prompt string) (string, error) {
	augmentedPrompt := prompt

	if p.hasDocuments {
		relevantDocs, err := p.vectorStore.SimilaritySearch(ctx, prompt, 2)
		if err != nil {
			return "", fmt.Errorf("failed to retrieve documents: %w", err)
		}

		var contextBuilder strings.Builder
		for _, doc := range relevantDocs {
			contextBuilder.WriteString(doc.PageContent + "\n")
		}

		augmentedPrompt, err = p.prompt.Format(map[string]string{
			"context":  contextBuilder.String(),
			"question": prompt,
		})
		if err != nil {
			return "", fmt.Errorf("failed to format prompt: %w", err)
		}
	}

	response, err := p.agent.Run(ctx, augmentedPrompt)
	if err != nil {
		return "", fmt.Errorf("agent run failed: %w", err)
	}

	return response, nil
}