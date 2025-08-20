package documentloaders

import (
	"fmt"
	"gogurt/documentloaders/pdf"
	"gogurt/documentloaders/text"
	"gogurt/types"
	"path/filepath"
)

// detects the file type and uses the appropriate loader
func LoadDocuments(filePath string) ([]types.Document, error) {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".txt":
		return text.NewTextLoader(filePath)
	case ".pdf":
		return pdf.NewPDFLoader(filePath)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}
}