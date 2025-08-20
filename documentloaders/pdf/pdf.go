package pdf

import (
	"gogurt/types"
	"strings"

	"rsc.io/pdf"
)

// NewPDFLoader reads a PDF file from a given path and extracts its text.
func NewPDFLoader(filePath string) ([]types.Document, error) {
	reader, err := pdf.Open(filePath)
	if err != nil {
		return nil, err
	}

	var textBuilder strings.Builder
	numPages := reader.NumPage()

	for i := 1; i <= numPages; i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}
		
		var pageText []string
		texts := page.Content().Text
		for _, t := range texts {
			pageText = append(pageText, t.S)
		}
		// Join the text elements with spaces and add a newline for the page break
		textBuilder.WriteString(strings.Join(pageText, " ") + "\n")
	}

	doc := types.Document{
		PageContent: textBuilder.String(),
		Metadata:    map[string]interface{}{"source": filePath},
	}

	return []types.Document{doc}, nil
}