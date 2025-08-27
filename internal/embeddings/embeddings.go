package embeddings

import (
	"context"
	"gogurt/internal/types"
)

type Embedder interface {
	EmbedDocuments(ctx context.Context, docs []types.Document) ([][]float32, error)
	EmbedQuery(ctx context.Context, text string) ([]float32, error)
	EmbedAll(ctx context.Context, docs []types.Document, workers int) ([][]float32, error)
	AEmbedQuery(ctx context.Context, text string) (<-chan []float32, <-chan error)
	AEmbedDocuments(ctx context.Context, docs []types.Document) (<-chan [][]float32, <-chan error)
}
