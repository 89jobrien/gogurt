package code

import (
	"gogurt/splitters/code/golang"
	"gogurt/splitters/code/javascript"
	"gogurt/splitters/code/python"
	"gogurt/splitters/recursive"
	"gogurt/types"
	"path/filepath"
)

type Splitter struct {
	ChunkSize    int
	ChunkOverlap int
}

func New(chunkSize, chunkOverlap int) *Splitter {
	return &Splitter{ChunkSize: chunkSize, ChunkOverlap: chunkOverlap}
}

func (s *Splitter) SplitDocuments(docs []types.Document) []types.Document {
	var finalChunks []types.Document
	fallbackSplitter := recursive.New(s.ChunkSize, s.ChunkOverlap)

	for _, doc := range docs {
		source, ok := doc.Metadata["source"].(string)
		if !ok {
			finalChunks = append(finalChunks, fallbackSplitter.SplitDocuments([]types.Document{doc})...)
			continue
		}

		var chunks []string
		switch filepath.Ext(source) {
		case ".go":
			chunks = golang.Split(doc.PageContent)
		case ".py":
			chunks = python.Split(doc.PageContent)
		case ".js", ".ts":
			chunks = javascript.Split(doc.PageContent)
		default:
			chunks = []string{doc.PageContent}
		}

		for _, chunk := range chunks {
			if len(chunk) > s.ChunkSize {
				subChunks := fallbackSplitter.SplitDocuments([]types.Document{{PageContent: chunk, Metadata: doc.Metadata}})
				finalChunks = append(finalChunks, subChunks...)
			} else {
				finalChunks = append(finalChunks, types.Document{
					PageContent: chunk,
					Metadata:    doc.Metadata,
				})
			}
		}
	}
	return finalChunks
}