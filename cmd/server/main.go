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
	
	socketServer := NewSocketIOServer()
	go func() {
		if err := socketServer.Serve(); err != nil {
			logger.Fatal("Socket.IO listen error: %s\n", err)
		}
	}()
	defer socketServer.Close()

	mux := http.NewServeMux()

	mux.Handle("/socket.io/", socketServer)

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
		file_tools.SaveToScratchpadTool,
		file_tools.ReadScratchpadTool,
		web.DuckDuckGoSearchTool,
		web.SerpAPISearchTool,
	})
	for _, err := range errs {
		if err != nil {
			logger.Error("Registration error: %v", err)
		}
	}

	stats := registry.Stats()
	logger.Info("Available Tools: %v", stats.ToolNames)
	logger.Info("Available Routes: %v", registeredRoutes)
	logger.Info("Available Pipes: [ /workflow /ddgs /serpapi ]")
	logger.Info("Server running at :8080")

	handler := MiddlewareHandler(mux)
	logger.ErrorCtx(context.Background(), "Server error: %v", http.ListenAndServe(":8080", handler))
}