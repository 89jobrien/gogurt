package vectorstores

import (
	"context"
	"gogurt/internal/types"
)

// VectorStore is the interface for a vector database.
// All methods are non-blocking and return results via channels.
type VectorStore interface {
	// AddDocuments adds documents to the vector store asynchronously.
	AddDocuments(ctx context.Context, docs []types.Document) <-chan error
	// SimilaritySearch performs a similarity search asynchronously.
	SimilaritySearch(ctx context.Context, query string, k int) (<-chan []types.Document, <-chan error)
}