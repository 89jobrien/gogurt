package embeddings

import (
	"context"
	"gogurt/types"
)

type Embedder interface {
	EmbedDocuments(ctx context.Context, docs []types.Document) ([][]float32, error)
	EmbedQuery(ctx context.Context, text string) ([]float32, error)
}