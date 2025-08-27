package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"gogurt/internal/config"
	"gogurt/internal/pipes"
	"net/http"
	"time"
)

// SerpApiRequest defines the structure for the incoming JSON request.
type SerpApiRequest struct {
	Prompt string `json:"prompt"`
}

// SerpApiResponse defines the structure for the JSON response.
type SerpApiResponse struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

// SerpApiHandler handles HTTP requests to execute the plan-and-execute workflow.
func SerpApiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SerpApiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request: could not decode JSON", http.StatusBadRequest)
		return
	}

	if req.Prompt == "" {
		http.Error(w, "Bad request: prompt cannot be empty", http.StatusBadRequest)
		return
	}

	// Use a context with a timeout to prevent long-running requests.
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	// Load configuration.
	cfg := config.Load()

	// Create and run the SerpApi pipe.
	serpApiPipe, err := pipes.NewSerpApiPipe(ctx, cfg)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create SerpApi pipe: %v", err), http.StatusInternalServerError)
		return
	}

	result, err := serpApiPipe.Run(ctx, req.Prompt)

	// Prepare the response
	w.Header().Set("Content-Type", "application/json")
	resp := SerpApiResponse{}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = err.Error()
	} else {
		w.WriteHeader(http.StatusOK)
		resp.Result = result
	}

	json.NewEncoder(w).Encode(resp)
}