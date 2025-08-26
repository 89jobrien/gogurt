package handlers

import (
	"encoding/json"
	"net/http"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get("http://localhost:8080/health")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status": "ERROR"}`))
	} else if res != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "OK"}`))
	}
	status := map[string]any{
		"name":    "Gogurt API",
		"status":  "OK",
		"uptime":  "unknown",
		"version": "v0.1.0",
		"message": "Server is healthy and running",
	}
	if res != nil {
		status["uptime"] = res.Header.Get("Uptime")
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}
