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

// WorkflowRequest defines the structure for the incoming JSON request.
type WorkflowRequest struct {
	Prompt string `json:"prompt"`
}

// WorkflowResponse defines the structure for the JSON response.
type WorkflowResponse struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

// WorkflowHandler handles HTTP requests to execute the plan-and-execute workflow.
func WorkflowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req WorkflowRequest
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

	// Load configuration. In a larger application, you might inject this
	// as a dependency into the handler instead of loading it each time.
	cfg := config.Load()

	// Create the workflow pipe.
	workflowPipe, err := pipes.NewWorkflowPipe(ctx, cfg)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create workflow pipe: %v", err), http.StatusInternalServerError)
		return
	}

	// Run the pipe asynchronously and wait for the result from the channels.
	resultCh, errCh := workflowPipe.Run(ctx, req.Prompt)

	// Prepare the response
	w.Header().Set("Content-Type", "application/json")
	resp := WorkflowResponse{}

	select {
	case result := <-resultCh:
		w.WriteHeader(http.StatusOK)
		resp.Result = result
	case err := <-errCh:
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = err.Error()
	case <-ctx.Done():
		w.WriteHeader(http.StatusRequestTimeout)
		resp.Error = "Request timed out or was canceled."
	}

	json.NewEncoder(w).Encode(resp)
}