package web

import (
	"encoding/json"
	"fmt"
	"gogurt/internal/tools"
	"reflect"

	"github.com/sap-nocops/duckduckgogo/client"
)

// SearchResult defines the structure for a single search result.
type SearchResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

// DuckDuckGoSearchArgs defines the input arguments for the search tool.
type DuckDuckGoSearchArgs struct {
	Query      string `json:"query"`
	NumResults int    `json:"num_results"`
}

// DuckDuckGoSearch performs a web search and returns a list of results as a JSON string.
func DuckDuckGoSearch(args DuckDuckGoSearchArgs) (string, error) {
	if args.Query == "" {
		return "", fmt.Errorf("query cannot be empty")
	}

	if args.NumResults <= 0 {
		args.NumResults = 5 // Default to 5 results if not specified or invalid
	}

	ddg := client.NewDuckDuckGoSearchClient()
	results, err := ddg.SearchLimited(args.Query, args.NumResults)
	if err != nil {
		return "", fmt.Errorf("DuckDuckGo search failed: %w", err)
	}

	if len(results) == 0 {
		return "[]", nil // Return an empty JSON array if no results are found
	}

	// Convert the results to our structured format
	var searchResults []SearchResult
	for _, res := range results {
		searchResults = append(searchResults, SearchResult{
			Title:   res.Title,
			URL:     res.FormattedUrl,
			Snippet: res.Snippet,
		})
	}

	// Marshal the structured results into a JSON string
	jsonResult, err := json.Marshal(searchResults)
	if err != nil {
		return "", fmt.Errorf("failed to marshal search results to JSON: %w", err)
	}

	return string(jsonResult), nil
}

var DuckDuckGoSearchTool = &tools.Tool{
	Name:        "duckduckgo_search",
	Description: "Performs a web search using DuckDuckGo and returns a JSON array of results, each with a title, URL, and snippet.",
	Func:        reflect.ValueOf(DuckDuckGoSearch),
	InputSchema: tools.GenInputSchema(reflect.TypeOf(DuckDuckGoSearchArgs{})),
	Example:     `{"query":"what is a gopher?", "num_results": 3}`,
	Metadata:    map[string]any{"category": "web"},
}