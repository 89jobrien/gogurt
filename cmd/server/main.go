package server

import (
	"context"
	"gogurt/internal/logger"
	"gogurt/internal/tools"
	"gogurt/internal/tools/file_tools"
	"gogurt/internal/tools/web"
	"gogurt/internal/types"
	"net/http"
	"os"
	"sort"
)

func Serve() {
	logFilePath := "logs/gogurt.log"
	logFile, err := logger.OpenLogFile(logFilePath)
	if err != nil {
		panic(err)
	}

	log := logger.NewLogger(
		os.Stdout,
		logFile,
		types.FormatText,
		types.FormatJSON,
	)

	logger.SetDefaultLogger(log)
	mux := http.NewServeMux()
	registeredRoutes := RegisterRoutes(mux)
	sort.Strings(registeredRoutes)

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
		file_tools.ReadFileTool,
		file_tools.WriteFileTool,
		file_tools.ListFilesTool,
		web.DuckDuckGoSearchTool,
	})
	for _, err := range errs {
		if err != nil {
			logger.Error("Registration error: %v", err)
		}
	}

	stats := registry.Stats()
	logger.Info("Available Tools: %v", stats.ToolNames)
	logger.Info("Available Routes: %v", registeredRoutes)
	logger.Info("Available Pipes: [ /workflow /ddgs ]")
	logger.Info("Server running at :8080")

	handler := MiddlewareHandler(mux)
	logger.ErrorCtx(context.Background(), "Server error: %v", http.ListenAndServe(":8080", handler))
}