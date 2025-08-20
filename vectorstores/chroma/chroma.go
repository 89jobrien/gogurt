package chroma

import (
	"context"
	"fmt"
	"gogurt/embeddings"
	ggtypes "gogurt/types"
	"gogurt/vectorstores"

	chroma "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
)

type Store struct {
	client     *chroma.Client
	embedder   embeddings.Embedder
	collection *chroma.Collection
}

// creates a new Chroma vector store using the documented API.
func New(ctx context.Context, url string, embedder embeddings.Embedder) (vectorstores.VectorStore, error) {
	client, err := chroma.NewClient(
		chroma.WithBasePath(url),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create chroma client: %w", err)
	}

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

func (s *Store) AddDocuments(ctx context.Context, docs []ggtypes.Document) error {
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

	_, err = s.collection.Add(
		ctx,
		types.NewEmbeddingsFromFloat32(docEmbeddings),
		nil,
		texts,
		ids,
	)
	if err != nil {
		return fmt.Errorf("failed to add documents to collection: %w", err)
	}

	return nil
}

func (s *Store) SimilaritySearch(ctx context.Context, query string, k int) ([]ggtypes.Document, error) {
	queryVector, err := s.embedder.EmbedQuery(ctx, query)
	if err != nil {
		return nil, err
	}

	results, err := s.collection.QueryWithOptions(
		ctx,
		types.WithQueryEmbeddings(types.NewEmbeddingsFromFloat32([][]float32{queryVector})),
		types.WithNResults(int32(k)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query collection: %w", err)
	}

	var documents []ggtypes.Document
	for _, docStr := range results.Documents[0] {
		documents = append(documents, ggtypes.Document{PageContent: docStr})
	}

	return documents, nil
}