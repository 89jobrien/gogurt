package chroma

import (
	"context"
	"fmt"
	"gogurt/internal/config"
	ggtypes "gogurt/internal/types"

	chromadb "github.com/amikos-tech/chroma-go/pkg/api/v2"
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

func ANew(ctx context.Context, cfg *config.Config) (<-chan *Store, <-chan error) {
	out := make(chan *Store, 1)
	errCh := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errCh)
		store, err := New(cfg)
		if err != nil {
			errCh <- err
			return
		}
		out <- store
	}()
	return out, errCh
}

func (s *Store) AddDocuments(ctx context.Context, docs []ggtypes.Document) error {
	if s.Col == nil {
		return fmt.Errorf("collection not initialized")
	}
	ids := make([]string, len(docs))
	texts := make([]string, len(docs))
	metadatas := make([]chromadb.DocumentMetadata, len(docs))
	for i, d := range docs {
		ids[i] = fmt.Sprintf("doc-%d", i)
		texts[i] = d.PageContent
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

func (s *Store) AAddDocuments(ctx context.Context, docs []ggtypes.Document) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		errCh <- s.AddDocuments(ctx, docs)
	}()
	return errCh
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
	metadataKeys := []string{"str", "int", "float"}
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

func (s *Store) ASimilaritySearch(ctx context.Context, query string, k int) (<-chan []ggtypes.Document, <-chan error) {
	out := make(chan []ggtypes.Document, 1)
	errCh := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errCh)
		result, err := s.SimilaritySearch(ctx, query, k)
		if err != nil {
			errCh <- err
			return
		}
		out <- result
	}()
	return out, errCh
}