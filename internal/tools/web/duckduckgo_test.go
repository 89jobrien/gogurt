package web

import (
	"encoding/json"
	"testing"
)

func TestDuckDuckGoSearchTool(t *testing.T) {
	// Test a simple query that should return results
	result, err := DuckDuckGoSearchTool.Call(`{"query":"golang", "num_results": 2}`)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// The result should be a JSON string. Let's unmarshal it.
	var searchResults []SearchResult
	err = json.Unmarshal([]byte(result.(string)), &searchResults)
	if err != nil {
		t.Fatalf("failed to unmarshal search result JSON: %v", err)
	}

	if len(searchResults) == 0 {
		t.Errorf("expected search results, but got none")
	}

	if len(searchResults) > 2 {
		t.Errorf("expected at most 2 results, but got %d", len(searchResults))
	}

	// Check that the first result has content
	if searchResults[0].Title == "" || searchResults[0].URL == "" || searchResults[0].Snippet == "" {
		t.Errorf("expected the first search result to have a title, URL, and snippet")
	}
}