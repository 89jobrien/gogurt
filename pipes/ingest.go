package pipes

import (
	"context"
	"fmt"
	"gogurt/config"
	"gogurt/documentloaders"
	"gogurt/embeddings"
	"gogurt/factories"
	"gogurt/splitters"
	"gogurt/vectorstores"
)

type IngestPipe struct {
	VectorStore  vectorstores.VectorStore
	splitter     splitters.Splitter
	embedder     embeddings.Embedder
	documentPath string
}

func NewIngestPipe(ctx context.Context, cfg *config.Config, documentPath string) (*IngestPipe, error) {
	c.Write("Setting up document ingestion pipeline...")

	splitter := factories.GetSplitter(cfg)
	embedder := factories.GetEmbedder(cfg)
	vectorStore := factories.GetVectorStore(cfg, embedder)

	return &IngestPipe{
		VectorStore:  vectorStore,
		splitter:     splitter,
		embedder:     embedder,
		documentPath: documentPath,
	}, nil
}

// Run loads, splits, and embeds documents into the vector store
func (i *IngestPipe) Run(ctx context.Context) error {
	c.Write("Starting document ingestion", "path", i.documentPath)

	docs, err := documentloaders.LoadDocuments(i.documentPath)
	if err != nil {
		return fmt.Errorf("failed to load documents from %s: %w", i.documentPath, err)
	}

	if len(docs) == 0 {
		c.Warn("No documents found at specified path %v", i.documentPath)
		return fmt.Errorf("no documents found at %s", i.documentPath)
	}

	chunks := i.splitter.SplitDocuments(docs)

	if len(chunks) == 0 {
		c.Warn("No chunks created from documents")
		return fmt.Errorf("no chunks created from documents")
	}

	err = i.VectorStore.AddDocuments(ctx, chunks)
	if err != nil {
		return fmt.Errorf("failed to add documents to vector store: %w", err)
	}

	c.Write("Document ingestion completed successfully", 
		"documents_loaded", len(docs), 
		"chunks_created", len(chunks))

	return nil
}

// GetVectorStore returns the vector store instance (useful for metrics)
func (i *IngestPipe) GetVectorStore() vectorstores.VectorStore {
	return i.VectorStore
}