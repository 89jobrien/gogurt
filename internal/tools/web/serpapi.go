package web

import (
	"encoding/json"
	"fmt"
	"gogurt/internal/tools"
	"os"
	"reflect"

	g "github.com/serpapi/google-search-results-golang"
)

// SerpAPISearchArgs defines the input arguments for the SerpApi search tool.
type SerpAPISearchArgs struct {
	Query string `json:"query"`
}

// SerpAPISearch performs a web search using SerpApi and returns a list of results as a JSON string.
func SerpAPISearch(args SerpAPISearchArgs) (string, error) {
	if args.Query == "" {
		return "", fmt.Errorf("query cannot be empty")
	}

	apiKey := os.Getenv("SERPAPI_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("SERPAPI_KEY environment variable not set")
	}

	parameter := map[string]string{
		"engine": "google",
		"q":      args.Query,
	}

	search := g.NewGoogleSearch(parameter, apiKey)
	results, err := search.GetJSON()
	if err != nil {
		return "", fmt.Errorf("SerpApi search failed: %w", err)
	}

	organicResults, ok := results["organic_results"].([]interface{})
	if !ok || len(organicResults) == 0 {
		return "[]", nil
	}

	jsonResult, err := json.Marshal(organicResults)
	if err != nil {
		return "", fmt.Errorf("failed to marshal search results to JSON: %w", err)
	}

	return string(jsonResult), nil
}

var SerpAPISearchTool = &tools.Tool{
	Name:        "serpapi_search",
	Description: "Performs a web search using SerpApi and Google, returning a JSON array of results.",
	Func:        reflect.ValueOf(SerpAPISearch),
	InputSchema: tools.GenInputSchema(reflect.TypeOf(SerpAPISearchArgs{})),
	Example:     `{"query":"what is a gopher?"}`,
	Metadata:    map[string]any{"category": "web"},
}