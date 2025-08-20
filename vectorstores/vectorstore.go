package vectorstores

import (
	"context"
	"gogurt/types"
)

type VectorStore interface {
	AddDocuments(ctx context.Context, docs []types.Document) error
	SimilaritySearch(ctx context.Context, query string, k int) ([]types.Document, error)
}