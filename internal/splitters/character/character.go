package character

import (
	"gogurt/internal/types"
)

type CharSplitter struct {
	ChunkSize   int
	ChunkOverlap int
}

func New(chunkSize, chunkOverlap int) *CharSplitter {
	return &CharSplitter{ChunkSize: chunkSize, ChunkOverlap: chunkOverlap}
}

func (s *CharSplitter) SplitDocuments(docs []types.Document) []types.Document {
	var chunks []types.Document
	for _, doc := range docs {
		content := doc.PageContent
		for i := 0; i < len(content); i += s.ChunkSize - s.ChunkOverlap {
			end := i + s.ChunkSize
			if end > len(content) {
				end = len(content)
			}
			chunks = append(chunks, types.Document{
				PageContent: content[i:end],
				Metadata:    doc.Metadata,
			})
			if end == len(content) {
				break
			}
		}
	}
	return chunks
}