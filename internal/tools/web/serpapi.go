package web

import (
	"encoding/json"
	"fmt"
	"gogurt/internal/tools"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
)

// SerpAPIResult defines the structure for a single search result from SerpApi.
type SerpAPIResult struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Snippet string `json:"snippet"`
}

// SerpAPISearchArgs defines the input arguments for the SerpApi search tool.
type SerpAPISearchArgs struct {
	Query      string `json:"query"`
	NumResults int    `json:"num_results"`
}

// SerpAPISearch performs a web search using SerpApi and returns a list of results as a JSON string.
func SerpAPISearch(args SerpAPISearchArgs) (string, error) {
	if args.Query == "" {
		return "", fmt.Errorf("query cannot be empty")
	}

	if args.NumResults <= 0 {
		args.NumResults = 5 // Default to 5 results
	}

	apiKey := os.Getenv("SERPAPI_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("SERPAPI_KEY environment variable not set")
	}

	url := fmt.Sprintf("https://serpapi.com/search.json?engine=duckduckgo&q=%s&num=%d&api_key=%s", args.Query, args.NumResults, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to make request to SerpApi: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("SerpApi request failed with status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read SerpApi response body: %w", err)
	}

	var response struct {
		OrganicResults []SerpAPIResult `json:"organic_results"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal SerpApi response: %w", err)
	}

	if len(response.OrganicResults) == 0 {
		return "[]", nil
	}

	jsonResult, err := json.Marshal(response.OrganicResults)
	if err != nil {
		return "", fmt.Errorf("failed to marshal search results to JSON: %w", err)
	}

	return string(jsonResult), nil
}

var SerpAPISearchTool = &tools.Tool{
	Name:        "serpapi_search",
	Description: "Performs a web search using SerpApi and DuckDuckGo, returning a JSON array of results.",
	Func:        reflect.ValueOf(SerpAPISearch),
	InputSchema: tools.GenInputSchema(reflect.TypeOf(SerpAPISearchArgs{})),
	Example:     `{"query":"what is a gopher?", "num_results": 3}`,
	Metadata:    map[string]any{"category": "web"},
}