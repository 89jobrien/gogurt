package server

import (
	"gogurt/cmd/server/handlers"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) []string {
	routes := map[string]http.Handler{
		"/tool":       http.HandlerFunc(handlers.ToolHandler),
		"/health":     http.HandlerFunc(handlers.HealthHandler),
		"/status":     http.HandlerFunc(handlers.StatusHandler),
		"/ping":       http.HandlerFunc(handlers.PingHandler),
		"/metrics":    http.HandlerFunc(handlers.MetricsHandler),
		"/version":    http.HandlerFunc(handlers.VersionHandler),
		"/docs":       http.HandlerFunc(handlers.DocsHandler),
		"/workflow":   http.HandlerFunc(handlers.WorkflowHandler),
		"/ddgs":       http.HandlerFunc(handlers.DDGSHandler),
		"/serpapi":    http.HandlerFunc(handlers.SerpApiHandler),
		"/agents":     http.HandlerFunc(handlers.AgentsHandler),
	}
	var routePaths []string
	for path, handler := range routes {
		mux.Handle(path, handler)
		routePaths = append(routePaths, path)
	}

	return routePaths
}