package pdf

import (
	"gogurt/types"
	"os"

	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

// reads a PDF file from a given path and extracts its text
func NewPDFLoader(filePath string) ([]types.Document, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return nil, err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return nil, err
	}

	var allText string
	for i := 1; i <= numPages; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			return nil, err
		}

		ex, err := extractor.New(page)
		if err != nil {
			return nil, err
		}

		text, err := ex.ExtractText()
		if err != nil {
			return nil, err
		}
		allText += text + "\n"
	}

	doc := types.Document{
		PageContent: allText,
		Metadata:    map[string]interface{}{"source": filePath},
	}

	return []types.Document{doc}, nil
}