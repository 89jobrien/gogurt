package splitters

import "gogurt/internal/types"

type Splitter interface {
	SplitDocuments(docs []types.Document) []types.Document
}