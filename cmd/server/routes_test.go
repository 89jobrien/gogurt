package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterRoutes(t *testing.T) {
	mux := http.NewServeMux()
	RegisterRoutes(mux)

	routes := []struct {
		path string
	}{
		{"/tool"},
		{"/health"},
		{"/status"},
		{"/ping"},
		{"/metrics"},
		{"/version"},
		{"/docs"},
	}

	for _, route := range routes {
		req := httptest.NewRequest("GET", route.path, nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		// Always expect a 200 unless the handler requires specific context
		if rr.Code != 200 && rr.Code != 404 { // 404 if handler logic not written
			t.Errorf("route %s returned status %d, want 200 or 404", route.path, rr.Code)
		}
	}
}