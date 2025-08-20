package markdown

import (
	"bytes"
	"gogurt/splitters/recursive"
	"gogurt/types"
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
		reader := text.NewReader([]byte(doc.PageContent))
		rootNode := mdParser.Parse(reader)

		var headerChunks []string
		var currentChunk bytes.Buffer

		ast.Walk(rootNode, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
			if entering {
				if n.Kind() == ast.KindHeading {
					if currentChunk.Len() > 0 {
						headerChunks = append(headerChunks, strings.TrimSpace(currentChunk.String()))
						currentChunk.Reset()
					}
				}
				
				if n.Type() == ast.TypeBlock || n.Type() == ast.TypeInline {
					lines := n.Lines()
					for i := 0; i < lines.Len(); i++ {
						line := lines.At(i)
						currentChunk.Write(line.Value(reader.Source()))
					}
					if n.Type() == ast.TypeBlock && n.Kind() != ast.KindDocument {
						currentChunk.WriteString("\n")
					}
				}
			}
			return ast.WalkContinue, nil
		})

		if currentChunk.Len() > 0 {
			headerChunks = append(headerChunks, strings.TrimSpace(currentChunk.String()))
		}

		for _, chunk := range headerChunks {
			if len(chunk) == 0 {
				continue
			}
			if len(chunk) > s.ChunkSize {
				subChunks := recursiveMarkdownSplitter.SplitDocuments([]types.Document{{PageContent: chunk}})
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