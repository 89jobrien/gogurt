package ollama

import (
	"context"
	"fmt"
	"gogurt/internal/config"
	"gogurt/internal/types"
	"sync"

	"github.com/ollama/ollama/api"
)

type Embedder struct {
	client *api.Client
	model  string
}

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

// Async: AEmbedDocuments
func (e *Embedder) AEmbedDocuments(ctx context.Context, docs []types.Document) (<-chan [][]float32, <-chan error) {
	out := make(chan [][]float32, 1)
	errCh := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errCh)
		embeddings, err := e.EmbedDocuments(ctx, docs)
		if err != nil {
			errCh <- err
			return
		}
		out <- embeddings
	}()
	return out, errCh
}

// Async: AEmbedQuery
func (e *Embedder) AEmbedQuery(ctx context.Context, text string) (<-chan []float32, <-chan error) {
	out := make(chan []float32, 1)
	errCh := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errCh)
		embedding, err := e.EmbedQuery(ctx, text)
		if err != nil {
			errCh <- err
			return
		}
		out <- embedding
	}()
	return out, errCh
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

// Async: AEmbedAll
func (e *Embedder) AEmbedAll(ctx context.Context, docs []types.Document, workers int) (<-chan [][]float32, <-chan error) {
	out := make(chan [][]float32, 1)
	errCh := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errCh)
		result, err := e.EmbedAll(ctx, docs, workers)
		if err != nil {
			errCh <- err
			return
		}
		out <- result
	}()
	return out, errCh
}

func (e *Embedder) EmbedAll(ctx context.Context, docs []types.Document, workers int) ([][]float32, error) {
	work := make(chan types.Document)
	out := make(chan []float32)
	go func() {
		for _, doc := range docs {
			work <- doc
		}
		close(work)
	}()
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for doc := range work {
				emb, err := e.EmbedQuery(ctx, doc.PageContent)
				if err != nil {
					continue // You might want to send the error out
				}
				out <- emb
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	var result [][]float32
	for emb := range out {
		result = append(result, emb)
	}
	return result, nil
}
