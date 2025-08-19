package simple

import (
	"context"
	"gogurt/embeddings"
	"gogurt/types"
	"math"
	"sort"
)

type VectorStore struct {
	embedder  embeddings.Embedder
	documents []types.Document
	vectors   [][]float32
}

func New(embedder embeddings.Embedder) *VectorStore {
	return &VectorStore{embedder: embedder}
}

func (s *VectorStore) AddDocuments(ctx context.Context, docs []types.Document) error {
	s.documents = append(s.documents, docs...)
	embeddings, err := s.embedder.EmbedDocuments(ctx, docs)
	if err != nil {
		return err
	}
	s.vectors = append(s.vectors, embeddings...)
	return nil
}

type searchResult struct {
	document   types.Document
	similarity float64
}

func (s *VectorStore) SimilaritySearch(ctx context.Context, query string, k int) ([]types.Document, error) {
	queryVector, err := s.embedder.EmbedQuery(ctx, query)
	if err != nil {
		return nil, err
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
	for i := range topK {
		documents = append(documents, results[i].document)
	}
	return documents, nil
}

func cosineSimilarity(a, b []float32) float64 {
	var dotProduct float64
	var normA, normB float64
	for i := range a {
		dotProduct += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}