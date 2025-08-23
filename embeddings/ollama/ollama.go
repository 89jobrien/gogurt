package ollama

import (
	"context"
	"fmt"

	"gogurt/config"
	"gogurt/types"

	"github.com/ollama/ollama/api"
)

type Embedder struct {
	client *api.Client
	model  string
}

a = []float32

a = "string"


func New(cfg *config.Config) (*Embedder, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to create ollama client: %w", err)
	}
	return &Embedder{
		client: client,
		model:  cfg.OllamaEmbedModel,
	}, nil
}

func (e *Embedder) EmbedDocuments(ctx context.Context, docs []types.Document) ([][]float32, error) {
	embeddings := make([][]float32, len(docs))
	for i, doc := range docs {
		embedding, err := e.EmbedQuery(ctx, doc.PageContent)
		if err != nil {
			return nil, err
		}
		embeddings[i] = embedding
	}
	return embeddings, nil
}

func (e *Embedder) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	req := &api.EmbeddingRequest{
		Model:  e.model,
		Prompt: text,
	}
	res, err := e.client.Embeddings(ctx, req)
	if err != nil {
		return nil, err
	}

	// convert []float64 to []float32
	embeddingF32 := make([]float32, len(res.Embedding))
	for i, v := range res.Embedding {
		embeddingF32[i] = float32(v)
	}
	return embeddingF32, nil
}