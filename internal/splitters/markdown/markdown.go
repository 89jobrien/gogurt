package markdown

import (
	"gogurt/internal/splitters/recursive"
	"gogurt/internal/types"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type MarkdownSplitter struct {
	ChunkSize    int
	ChunkOverlap int
}

func New(chunkSize, chunkOverlap int) *MarkdownSplitter {
	return &MarkdownSplitter{ChunkSize: chunkSize, ChunkOverlap: chunkOverlap}
}

// splits a markdown document into smaller chunks using the recursive splitter
func (s *MarkdownSplitter) SplitDocuments(docs []types.Document) []types.Document {
	var finalChunks []types.Document
	mdParser := goldmark.DefaultParser()
	recursiveMarkdownSplitter := recursive.New(s.ChunkSize, s.ChunkOverlap)

	for _, doc := range docs {
		if strings.TrimSpace(doc.PageContent) == "" {
			continue
		}

		reader := text.NewReader([]byte(doc.PageContent))
		rootNode := mdParser.Parse(reader)

		var chunks []string
		var currentChunk strings.Builder

		for n := rootNode.FirstChild(); n != nil; n = n.NextSibling() {
			if n.Kind() == ast.KindHeading && currentChunk.Len() > 0 {
				chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
				currentChunk.Reset()
			}

			lines := n.Lines()
			for i := 0; i < lines.Len(); i++ {
				line := lines.At(i)
				currentChunk.Write(line.Value(reader.Source()))
				currentChunk.WriteString("\n")
			}

			if n.Kind() == ast.KindHeading {
				// normalize header text spacing but keep original internal whitespace
				headerText := strings.TrimRight(currentChunk.String(), "\n")
				// headerText2 := strings.SplitAfter(headerText, "# ")[1]
				// headerText3 := strings.TrimSpace(headerText2)

				currentChunk.Reset()
				currentChunk.WriteString(headerText)
				if n.NextSibling() != nil {
					currentChunk.WriteString("\n\n")
				}
			}
		}

		if currentChunk.Len() > 0 {
			chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
		}

		for _, chunk := range chunks {
			if len(chunk) > s.ChunkSize {
				subChunks := recursiveMarkdownSplitter.SplitDocuments([]types.Document{{PageContent: chunk, Metadata: doc.Metadata}})
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
