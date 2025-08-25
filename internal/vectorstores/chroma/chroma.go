package chroma

import (
	"context"
	"fmt"

	chromadb "github.com/amikos-tech/chroma-go/pkg/api/v2"

	"gogurt/internal/config"
	ggtypes "gogurt/internal/types"
)

type Store struct {
	Client chromadb.Client
	Col    chromadb.Collection
}

func New(cfg *config.Config) (*Store, error) {
	store := &Store{}
	client, err := chromadb.NewHTTPClient(
		chromadb.WithBaseURL(cfg.ChromaURL),
	)
	if err != nil {
		return nil, err
	}
	store.Client = client

	col, err := client.GetOrCreateCollection(context.Background(), cfg.ChromaCollection,
		chromadb.WithCollectionMetadataCreate(
			chromadb.NewMetadata(
				chromadb.NewStringAttribute("space", cfg.ChromaSpace),
				chromadb.NewIntAttribute("ef_construction", int64(cfg.ChromaEFConstruction)),
				chromadb.NewIntAttribute("ef_search", int64(cfg.ChromaEFSearch)),
				chromadb.NewIntAttribute("max_neighbors", int64(cfg.ChromaMaxNeighbors)),
			),
		),
	)
	if err != nil {
		return nil, err
	}
	store.Col = col
	return store, nil
}

func (s *Store) AddDocuments(ctx context.Context, docs []ggtypes.Document) error {
	if s.Col == nil {
		return fmt.Errorf("collection not initialized")
	}

	ids := make([]string, len(docs))
	texts := make([]string, len(docs))
	metadatas := make([]chromadb.DocumentMetadata, len(docs))

	for i, d := range docs {
		ids[i] = fmt.Sprintf("doc-%d", i) // you could also use ULID if you want uniqueness
		texts[i] = d.PageContent

		// Convert map[string]any to Chroma DocumentMetadata
		if d.Metadata != nil {
			attrs := []*chromadb.MetaAttribute{}
			for k, v := range d.Metadata {
				switch val := v.(type) {
				case string:
					attrs = append(attrs, chromadb.NewStringAttribute(k, val))
				case int:
					attrs = append(attrs, chromadb.NewIntAttribute(k, int64(val)))
				case float64:
					attrs = append(attrs, chromadb.NewFloatAttribute(k, val))
				default:
					// skip unsupported metadata types or stringify them
					attrs = append(attrs, chromadb.NewStringAttribute(k, fmt.Sprintf("%v", val)))
				}
			}
			metadatas[i] = chromadb.NewDocumentMetadata(attrs...)
		} else {
			metadatas[i] = chromadb.NewDocumentMetadata()
		}
	}

	return s.Col.Add(ctx,
		chromadb.WithIDGenerator(chromadb.NewULIDGenerator()),
		chromadb.WithTexts(texts...),
		chromadb.WithMetadatas(metadatas...),
	)
}

func (s *Store) SimilaritySearch(ctx context.Context, query string, k int) ([]ggtypes.Document, error) {
	if s.Col == nil {
		return nil, fmt.Errorf("collection not initialized")
	}

	resp, err := s.Col.Query(ctx,
		chromadb.WithQueryTexts(query),
		chromadb.WithNResults(k))
	if err != nil {
		return nil, err
	}

	textsGroups := resp.GetDocumentsGroups()
	metadatasGroups := resp.GetMetadatasGroups()

	// Define the keys you expect, or ideally get this from configuration
	metadataKeys := []string{"str", "int", "float"} // <-- set your known/extracted keys here

	docs := []ggtypes.Document{}
	for groupIdx, documents := range textsGroups {
		var metadatas []chromadb.DocumentMetadata
		if groupIdx < len(metadatasGroups) {
			metadatas = metadatasGroups[groupIdx]
		}
		for idx, doc := range documents {
			var metadata map[string]any
			if metadatas != nil && idx < len(metadatas) {
				md := make(map[string]any)
				for _, key := range metadataKeys {
					if val, ok := metadatas[idx].GetRaw(key); ok {
						md[key] = val
					}
				}
				metadata = md
			}
			docs = append(docs, ggtypes.Document{
				PageContent: doc.ContentString(),
				Metadata:    metadata,
			})
		}
	}

	return docs, nil
}
