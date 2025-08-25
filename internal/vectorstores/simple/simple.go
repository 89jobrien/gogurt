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

// creates a simple in-memory vector store.
func New(embedder embeddings.Embedder) vectorstores.VectorStore {
	return &Store{embedder: embedder}
}

func (s *Store) AddDocuments(ctx context.Context, docs []types.Document) error {
	if len(docs) == 0 {
		return nil
	}
	s.documents = append(s.documents, docs...)
	docEmbeddings, err := s.embedder.EmbedDocuments(ctx, docs)
	if err != nil {
		return err
	}
	s.vectors = append(s.vectors, docEmbeddings...)
	return nil
}



func (s *Store) SimilaritySearch(ctx context.Context, query string, k int) ([]types.Document, error) {
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
	if normA == 0 || normB == 0 {
		return 0.0
	}
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}