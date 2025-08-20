package splitter

import "gogurt/types"

type Splitter interface {
	SplitDocuments(docs []types.Document) []types.Document
}