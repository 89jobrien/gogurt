package server

import (
	"context"
	"gogurt/internal/logger"
	"gogurt/internal/tools"
	"gogurt/internal/types"
	"net/http"
	"os"
)

func Serve() {
	logFilePath := "logs/gogurt.log"
	logFile, err := logger.OpenLogFile(logFilePath)
	if err != nil {
		panic(err)
	}

	// Use NewLogger to set the custom file writer
	// Console will use os.Stdout; file will use gogurt.log, with text/json formats as you prefer
	log := logger.NewLogger(
		os.Stdout,
		logFile,
		types.FormatText,    // Console format
		types.FormatJSON,    // File format
	)
	
	// Set defaultLogger so convenience methods use our logger
	logger.SetDefaultLogger(log)
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
			logger.Error("Registration error: %v", err)
		}
	}

	stats := registry.Stats()
	logger.Info(
		"Available Tools: %v", stats.ToolNames)
	logger.Info("Server running at :8080")

	handler := MiddlewareHandler(mux)
	logger.ErrorCtx(context.Background(), "Server error: %v", http.ListenAndServe(":8080", handler))
}
