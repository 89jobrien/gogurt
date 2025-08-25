package pdf

import (
	"bytes"
	"gogurt/internal/types"
	"io"

	"github.com/ledongthuc/pdf"
)

// reads a PDF using the pdftotext utility for higher accuracy.
func NewPDFLoader(filePath string) ([]types.Document, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return nil, err
	}
	// the library requires the reader to be closed
	defer f.Close()

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(&buf, b)
	if err != nil {
		return nil, err
	}

	doc := types.Document{
		PageContent: buf.String(),
		Metadata:    map[string]any{"source": filePath},
	}

	return []types.Document{doc}, nil
}
