package simple

import (
	"context"
	"gogurt/internal/embeddings"
	"gogurt/internal/types"
	"gogurt/internal/vectorstores"
	"math"
	"sort"
)

type Store struct {
	embedder  embeddings.Embedder
	documents []types.Document
	vectors   [][]float32
}

type searchResult struct {
	document   types.Document
	similarity float64
}

// New creates a simple in-memory vector store.
func New(embedder embeddings.Embedder) vectorstores.VectorStore {
	return &Store{embedder: embedder}
}

// AddDocuments adds documents to the vector store asynchronously.
func (s *Store) AddDocuments(ctx context.Context, docs []types.Document) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		if len(docs) == 0 {
			errCh <- nil
			return
		}
		s.documents = append(s.documents, docs...)
		// Assuming embedder methods are async and return channels
		docEmbeddingsCh, embedErrCh := s.embedder.AEmbedDocuments(ctx, docs)

		select {
		case docEmbeddings := <-docEmbeddingsCh:
			s.vectors = append(s.vectors, docEmbeddings...)
			errCh <- nil
		case err := <-embedErrCh:
			errCh <- err
		case <-ctx.Done():
			errCh <- ctx.Err()
		}
	}()
	return errCh
}

// SimilaritySearch performs a similarity search asynchronously.
func (s *Store) SimilaritySearch(ctx context.Context, query string, k int) (<-chan []types.Document, <-chan error) {
	out := make(chan []types.Document, 1)
	errCh := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errCh)

		queryVectorCh, embedErrCh := s.embedder.AEmbedQuery(ctx, query)
		var queryVector []float32
		select {
		case queryVector = <-queryVectorCh:
		case err := <-embedErrCh:
			errCh <- err
			return
		case <-ctx.Done():
			errCh <- ctx.Err()
			return
		}

		var results []searchResult
		for i, vector := range s.vectors {
			similarity := cosineSimilarity(queryVector, vector)
			results = append(results, searchResult{
				document:   s.documents[i],
				similarity: similarity,
			})
		}
		sort.Slice(results, func(i, j int) bool {
			return results[i].similarity > results[j].similarity
		})

		topK := min(len(results), k)
		var documents []types.Document
		for i := 0; i < topK; i++ {
			documents = append(documents, results[i].document)
		}
		out <- documents
	}()
	return out, errCh
}

// cosineSimilarity is a synchronous helper function.
func cosineSimilarity(a, b []float32) float64 {
	var dotProduct float64
	var normA, normB float64
	for i := range a {
		dotProduct += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}
	if normA == 0 || normB == 0 {
		return 0.0
	}
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}