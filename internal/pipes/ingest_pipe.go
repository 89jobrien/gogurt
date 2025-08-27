package pipes

import (
	"context"
	"fmt"
	"gogurt/internal/config"
	"gogurt/internal/documentloaders"
	"gogurt/internal/embeddings"
	"gogurt/internal/factories"
	"gogurt/internal/splitters"
	"gogurt/internal/vectorstores"
)

// IngestPipe handles the asynchronous ingestion of documents into a vector store.
type IngestPipe struct {
	VectorStore  vectorstores.VectorStore
	splitter     splitters.Splitter
	embedder     embeddings.Embedder
	documentPath string
}

// NewIngestPipe creates a new document ingestion pipeline.
func NewIngestPipe(ctx context.Context, cfg *config.Config, documentPath string) (*IngestPipe, error) {
	c.Write("Setting up document ingestion pipeline...")

	// Factories are currently synchronous, but are used to construct the async-capable pipe.
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

// Run loads, splits, and embeds documents into the vector store asynchronously.
// It returns a channel that will receive an error if one occurs, or nil on success.
func (i *IngestPipe) Run(ctx context.Context) <-chan error {
	errCh := make(chan error, 1)

	go func() {
		defer close(errCh)
		c.Write("Starting document ingestion", "path", i.documentPath)

		// 1. Load documents from the specified path.
		docs, err := documentloaders.LoadDocuments(i.documentPath)
		if err != nil {
			errCh <- fmt.Errorf("failed to load documents from %s: %w", i.documentPath, err)
			return
		}
		if len(docs) == 0 {
			c.Warn("No documents found at specified path %v", i.documentPath)
			errCh <- fmt.Errorf("no documents found at %s", i.documentPath)
			return
		}

		// 2. Split the documents into smaller chunks.
		chunks := i.splitter.SplitDocuments(docs)
		if len(chunks) == 0 {
			c.Warn("No chunks created from documents")
			errCh <- fmt.Errorf("no chunks created from documents")
			return
		}

		// 3. Add the chunks to the vector store asynchronously.
		addErrCh := i.VectorStore.AddDocuments(ctx, chunks)
		select {
		case err := <-addErrCh:
			if err != nil {
				errCh <- fmt.Errorf("failed to add documents to vector store: %w", err)
				return
			}
			c.Write("Document ingestion completed successfully",
				"documents_loaded", len(docs),
				"chunks_created", len(chunks))
			errCh <- nil // Signal success
		case <-ctx.Done():
			errCh <- ctx.Err()
		}
	}()

	return errCh
}

// GetVectorStore returns the vector store instance.
func (i *IngestPipe) GetVectorStore() vectorstores.VectorStore {
	return i.VectorStore
}