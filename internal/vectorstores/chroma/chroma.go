package chroma

import (
	"context"
	"fmt"
	"gogurt/internal/config"
	ggtypes "gogurt/internal/types"
	"strings"

	chromadb "github.com/amikos-tech/chroma-go/pkg/api/v2"
)

type Store struct {
	Client chromadb.Client
	Col    chromadb.Collection
}

// New creates a new ChromaDB client and gets or creates a collection.
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

// AddDocuments adds documents to the collection asynchronously.
func (s *Store) AddDocuments(ctx context.Context, docs []ggtypes.Document) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		if s.Col == nil {
			errCh <- fmt.Errorf("collection not initialized")
			return
		}

		ids := make([]string, len(docs))
		texts := make([]string, len(docs))
		metadatas := make([]chromadb.DocumentMetadata, len(docs))

		for i, d := range docs {
			ids[i] = fmt.Sprintf("doc-%d", i)
			texts[i] = d.PageContent
			if d.Metadata != nil {
				var attrs []*chromadb.MetaAttribute
				var keys []string
				for k, v := range d.Metadata {
					keys = append(keys, k)
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
				attrs = append(attrs, chromadb.NewStringAttribute("__keys__", strings.Join(keys, ",")))
				metadatas[i] = chromadb.NewDocumentMetadata(attrs...)
			} else {
				metadatas[i] = chromadb.NewDocumentMetadata()
			}
		}

		err := s.Col.Add(ctx,
			chromadb.WithIDGenerator(chromadb.NewULIDGenerator()),
			chromadb.WithTexts(texts...),
			chromadb.WithMetadatas(metadatas...),
		)
		errCh <- err
	}()
	return errCh
}

// SimilaritySearch performs a query asynchronously.
func (s *Store) SimilaritySearch(ctx context.Context, query string, k int) (<-chan []ggtypes.Document, <-chan error) {
	out := make(chan []ggtypes.Document, 1)
	errCh := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errCh)
		if s.Col == nil {
			errCh <- fmt.Errorf("collection not initialized")
			return
		}

		resp, err := s.Col.Query(ctx,
			chromadb.WithQueryTexts(query),
			chromadb.WithNResults(k))
		if err != nil {
			errCh <- err
			return
		}

		textsGroups := resp.GetDocumentsGroups()
		metadatasGroups := resp.GetMetadatasGroups()
		var docs []ggtypes.Document

		for groupIdx, documents := range textsGroups {
			var metadatas []chromadb.DocumentMetadata
			if groupIdx < len(metadatasGroups) {
				metadatas = metadatasGroups[groupIdx]
			}
			for idx, doc := range documents {
				var metadata map[string]any
				if metadatas != nil && idx < len(metadatas) {
					md := make(map[string]any)
					keysStr, ok := metadatas[idx].GetString("__keys__")
					if ok {
						keys := strings.Split(keysStr, ",")
						for _, key := range keys {
							if val, ok := metadatas[idx].GetRaw(key); ok {
								md[key] = val
							}
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
		out <- docs
	}()

	return out, errCh
}