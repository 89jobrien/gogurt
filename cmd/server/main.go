package server

import (
	"gogurt/internal/tools"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func Serve() {
	// --- Structured Logger Setup ---
	// This configures a global logger that writes to gogurt.log.
	logFile, err := os.OpenFile("gogurt.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(logFile, nil))
	slog.SetDefault(logger)
	// --- End Logger Setup ---

	mux := http.NewServeMux()
	RegisterRoutes(mux)

	registry := tools.NewRegistry()
	errs := registry.RegisterBatch([]*tools.Tool{
		tools.PalindromeTool,
		tools.ConcatenateTool,
		tools.ReverseTool,
		tools.UppercaseTool,
		tools.AddTool,
		tools.DivideTool,
		tools.SubtractTool,
		tools.MultiplyTool,
	})
	for _, err := range errs {
		if err != nil {
			log.Println("Registration error:", err)
		}
	}

	stats := registry.Stats()
	log.Println("Tool Stats:")
	log.Printf("  Count: %d\n", stats.Count)
	log.Printf("  Available Tools: %v\n", stats.ToolNames)
	log.Printf("  Tool Categories: %v\n", stats.Categories)
	log.Printf("  Duplicates?: %v\n", stats.HasDups)
	log.Printf("  HasCategory: %v\n\n", stats.HasCategory)
	log.Println("Server running at :8080")

	handler := MiddlewareHandler(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
