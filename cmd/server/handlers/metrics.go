package handlers

import (
	"encoding/json"
	"net/http"
	"runtime"
)

func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	mem := &runtime.MemStats{}
	runtime.ReadMemStats(mem)
	metrics := map[string]any{
		"goroutines":   runtime.NumGoroutine(),
		"memory_bytes": mem.Alloc,
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
