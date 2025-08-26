package handlers

import (
	"encoding/json"
	"gogurt/internal/tools"
	"net/http"
)

func ToolHandler(w http.ResponseWriter, r *http.Request) {
		registry := tools.NewRegistry()
		var req struct {
			Name string `json:"name"`
			Args []byte `json:"args"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		result, err := registry.Call(req.Name, string(req.Args))
		resp := map[string]any{"result": result, "error": ""}
		if err != nil {
			resp["error"] = err.Error()
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}

