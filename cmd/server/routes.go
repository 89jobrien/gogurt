package server

import (
	"gogurt/cmd/server/handlers"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/tool", handlers.ToolHandler)
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/status", handlers.StatusHandler)
	mux.HandleFunc("/ping", handlers.PingHandler)
	mux.HandleFunc("/metrics", handlers.MetricsHandler)
	mux.HandleFunc("/version", handlers.VersionHandler)
	mux.HandleFunc("/docs", handlers.DocsHandler)
	mux.HandleFunc("/workflow", handlers.WorkflowHandler)
}
