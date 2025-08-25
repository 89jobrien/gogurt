package documentloaders

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"gogurt/internal/documentloaders/code"
	"gogurt/internal/documentloaders/markdown"
	"gogurt/internal/documentloaders/pdf"
	"gogurt/internal/documentloaders/text"
	"gogurt/internal/types"
)

// detects if the path is a file or a directory and loads accordingly.
func LoadDocuments(path string) ([]types.Document, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("could not access path %s: %w", path, err)
	}

	if fileInfo.IsDir() {
		return loadFromDirectory(path)
	}

	return loadFromFile(path)
}

// loads a single file using the appropriate loader.
func loadFromFile(filePath string) ([]types.Document, error) {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".txt":
		return text.NewTextLoader(filePath)
	case ".pdf":
		return pdf.NewPDFLoader(filePath)
	case ".md":
		return markdown.NewMarkdownLoader(filePath)
	case ".go", ".py", ".js", ".ts", ".java", ".cpp", ".c", ".rs":
		return code.NewCodeLoader(filePath)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}
}

// reads all supported files from a directory.
func loadFromDirectory(dirPath string) ([]types.Document, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("could not read directory %s: %w", dirPath, err)
	}

	var allDocs []types.Document
	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(dirPath, file.Name())
			docs, err := loadFromFile(filePath)
			if err != nil {
				// log the error for the specific file but continue with others
				slog.Warn("could not load file", "path", filePath, "error", err)
				continue
			}
			allDocs = append(allDocs, docs...)
		}
	}
	return allDocs, nil
}