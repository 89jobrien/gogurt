package server

import (
	"gogurt/cmd/server/handlers"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) []string {
	routes := map[string]http.HandlerFunc{
		"/tool":    handlers.ToolHandler,
		"/health":  handlers.HealthHandler,
		"/status":  handlers.StatusHandler,
		"/ping":    handlers.PingHandler,
		"/metrics": handlers.MetricsHandler,
		"/version": handlers.VersionHandler,
		"/docs":    handlers.DocsHandler,
		"/workflow":handlers.WorkflowHandler,
		"/ddgs":    handlers.DDGSHandler,
		"/serpapi":    handlers.SerpApiHandler,
	}
	var routePaths []string
	for path, handler := range routes {
		mux.HandleFunc(path, handler)
		routePaths = append(routePaths, path)
	}

	return routePaths
}