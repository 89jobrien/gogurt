package text

import (
	"gogurt/internal/types"
	"os"
)

// reads a plain text file from a given path
func NewTextLoader(filePath string) ([]types.Document, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	doc := types.Document{
		PageContent: string(content),
		Metadata:    map[string]interface{}{"source": filePath},
	}

	return []types.Document{doc}, nil
}