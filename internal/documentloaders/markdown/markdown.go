package markdown

import (
	"gogurt/internal/types"
	"os"
)

// reads a markdown file from a given path
func NewMarkdownLoader(filePath string) ([]types.Document, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	doc := types.Document{
		PageContent: string(content),
		Metadata:    map[string]any{"source": filePath},
	}

	return []types.Document{doc}, nil
}
