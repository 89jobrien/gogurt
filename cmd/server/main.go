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
	// Configure main logger for gogurt.log
	logFilePath := "logs/gogurt.log"
	logFile, err := logger.OpenLogFile(logFilePath)
	if err != nil {
		panic(err)
	}
	mainLogger := logger.NewLogger(
		os.Stdout,
		logFile,
		types.FormatText,
		types.FormatJSON,
	)

	// Set the main logger as the default for convenience functions
	logger.SetDefaultLogger(mainLogger)

	// Configure a separate logger for websocket events
	wsLogFilePath := "logs/websocket.log"
	wsLogFile, err := logger.OpenLogFile(wsLogFilePath)
	if err != nil {
		panic(err)
	}
	websocketLogger := logger.NewLogger(
		os.Stdout,
		wsLogFile,
		types.FormatText,
		types.FormatJSON,
	)

	// Pass the dedicated websocket logger to the Socket.IO server
	socketServer := NewSocketIOServer(websocketLogger)
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