// In vectorstores/chroma/chroma.go
package chroma

import (
	"context"
	"fmt"
	"gogurt/embeddings"
	gogurttypes "gogurt/types"
	"gogurt/vectorstores"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
)

// Store implements the vectorstores.VectorStore interface for ChromaDB.
type Store struct {
	client     *chromago.Client
	embedder   embeddings.Embedder
	collection *chromago.Collection // Corrected type from the root chromago package
}

// New creates a new Chroma vector store using the documented API.
func New(ctx context.Context, url string, embedder embeddings.Embedder) (vectorstores.VectorStore, error) {
	// 1. Use WithBasePath as the documented option for setting the URL
	client, err := chromago.NewClient(
		chromago.WithBasePath(url),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create chroma client: %w", err)
	}

	// 2. Use the documented CreateCollection method with the createOrGet flag
	// Setting embedding function to nil and distance function to L2 (default)
	col, err := client.CreateCollection(ctx, "gogurt-collection", nil, true, nil, types.L2)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create collection: %w", err)
	}

	return &Store{
		client:     client,
		embedder:   embedder,
		collection: col,
	}, nil
}

func (s *Store) AddDocuments(ctx context.Context, docs []gogurttypes.Document) error {
	if len(docs) == 0 {
		return nil
	}

	texts := make([]string, len(docs))
	ids := make([]string, len(docs))
	for i, d := range docs {
		texts[i] = d.PageContent
		ids[i] = fmt.Sprintf("doc-%d-%s", i, d.Metadata["source"])
	}

	docEmbeddings, err := s.embedder.EmbedDocuments(ctx, docs)
	if err != nil {
		return err
	}

	// 3. Use the collection.Add method which takes slices directly
	_, err = s.collection.Add(
		ctx,
		types.NewEmbeddingsFromFloat32(docEmbeddings),
		nil, // No metadata
		texts,
		ids,
	)
	if err != nil {
		return fmt.Errorf("failed to add documents to collection: %w", err)
	}

	return nil
}

func (s *Store) SimilaritySearch(ctx context.Context, query string, k int) ([]gogurttypes.Document, error) {
	queryVector, err := s.embedder.EmbedQuery(ctx, query)
	if err != nil {
		return nil, err
	}

	// The Query method takes query texts, so we use QueryWithOptions with our embeddings
	results, err := s.collection.QueryWithOptions(
		ctx,
		types.WithQueryEmbeddings(types.NewEmbeddingsFromFloat32([][]float32{queryVector})),
		types.WithNResults(int32(k)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query collection: %w", err)
	}

	var documents []gogurttypes.Document
	for _, docStr := range results.Documents[0] {
		documents = append(documents, gogurttypes.Document{PageContent: docStr})
	}

	return documents, nil
}