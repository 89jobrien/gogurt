package recursive

import (
	"gogurt/internal/types"
	"strings"
)

type RecursiveSplitter struct {
	ChunkSize    int
	ChunkOverlap int
	Separators   []string
}

func New(chunkSize, chunkOverlap int) *RecursiveSplitter {
	if chunkOverlap > chunkSize {
		chunkOverlap = chunkSize / 2 // safe default if overlap is too large
	}
	return &RecursiveSplitter{
		ChunkSize:    chunkSize,
		ChunkOverlap: chunkOverlap,
		Separators:   []string{"\n\n", "\n", " ", ""}, // default separators from most to least significant
	}
}

// applies the recursive logic to a list of documents
func (s *RecursiveSplitter) SplitDocuments(docs []types.Document) []types.Document {
	var finalChunks []types.Document
	for _, doc := range docs {
		// start recursive splitting with the initial list of separators.
		chunks := s.splitText(doc.PageContent, s.Separators)
		for _, chunk := range chunks {
			finalChunks = append(finalChunks, types.Document{
				PageContent: chunk,
				Metadata:    doc.Metadata,
			})
		}
	}
	return finalChunks
}

// recursively breaks down text based on the provided separators
func (s *RecursiveSplitter) splitText(text string, separators []string) []string {
	var finalChunks []string

	if len(text) < s.ChunkSize {
		return []string{text}
	}

	separator := ""
	for _, s := range separators {
		if strings.Contains(text, s) {
			separator = s
			break
		}
	}

	splits := strings.Split(text, separator)
	var goodSplits []string

	for _, split := range splits {
		if len(split) > s.ChunkSize {
			if len(goodSplits) > 0 {
				merged := s.mergeSplits(goodSplits, separator)
				finalChunks = append(finalChunks, merged...)
				goodSplits = []string{}
			}
			if len(separators) > 1 {
				recursiveChunks := s.splitText(split, separators[1:])
				finalChunks = append(finalChunks, recursiveChunks...)
			} else {
				finalChunks = append(finalChunks, s.splitByCharacter(split)...)
			}
		} else {
			goodSplits = append(goodSplits, split)
		}
	}

	if len(goodSplits) > 0 {
		merged := s.mergeSplits(goodSplits, separator)
		finalChunks = append(finalChunks, merged...)
	}

	return finalChunks
}

// intelligently groups smaller text splits into final chunks of the desired size and overlap.
func (s *RecursiveSplitter) mergeSplits(splits []string, separator string) []string {
	var docs []string
	var currentDoc []string
	totalLen := 0

	for _, d := range splits {
		splitLen := len(d)
		if totalLen+splitLen+len(separator) > s.ChunkSize && len(currentDoc) > 0 {
			docs = append(docs, strings.Join(currentDoc, separator))

			overlapLen := 0
			var newDoc []string
			for i := len(currentDoc) - 1; i >= 0; i-- {
				part := currentDoc[i]
				if overlapLen+len(part)+len(separator) > s.ChunkOverlap {
					break
				}
				overlapLen += len(part) + len(separator)
				newDoc = append([]string{part}, newDoc...)
			}
			currentDoc = newDoc
			totalLen = overlapLen
		}

		currentDoc = append(currentDoc, d)
		totalLen += splitLen + len(separator)
	}

	if len(currentDoc) > 0 {
		docs = append(docs, strings.Join(currentDoc, separator))
	}

	return docs
}

// fallback for when a chunk of text has no separators but is still too long
func (s *RecursiveSplitter) splitByCharacter(text string) []string {
	var chunks []string
	for i := 0; i < len(text); i += s.ChunkSize - s.ChunkOverlap {
		end := min(i + s.ChunkSize, len(text))
		chunks = append(chunks, text[i:end])
	}
	return chunks
}