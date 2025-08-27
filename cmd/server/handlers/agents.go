package handlers

import (
	"encoding/json"
	"net/http"
)

// AgentsHandler returns a list of available agents and pipes.
func AgentsHandler(w http.ResponseWriter, r *http.Request) {
	agents := []string{"workflow", "ddgs", "serpapi"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agents)
}