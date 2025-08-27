package server

import (
	"context"
	"encoding/json"
	"gogurt/internal/config"
	"gogurt/internal/logger"
	"gogurt/internal/pipes"
	"strings"
	"time"

	socketio "github.com/googollee/go-socket.io"
)

// SocketRequest defines the structure for incoming chat messages from the client.
type SocketRequest struct {
	Endpoint string `json:"endpoint"`
	Prompt   string `json:"prompt"`
}

// SocketResponse defines the structure for outgoing responses to the client.
type SocketResponse struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

// NewSocketIOServer creates and configures a new Socket.IO server.
func NewSocketIOServer(wsLogger *logger.Logger) *socketio.Server {
	server := socketio.NewServer(nil)

	// Handles new client connections.
	server.OnConnect("/", func(s socketio.Conn) error {
		logger.Info("Socket connected: %s", s.ID())
		return nil
	})

	// Listens for the 'pipe-request' event from the client.
	server.OnEvent("/", "pipe-request", func(s socketio.Conn, data string) {
		var req SocketRequest
		if err := json.Unmarshal([]byte(data), &req); err != nil {
			logger.Error("Socket request unmarshal error: %v", err)
			s.Emit("pipe-response", SocketResponse{Error: "Invalid request format"})
			return
		}

		logger.Info("Received pipe-request for endpoint '%s' with prompt: '%s'", req.Endpoint, req.Prompt)

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		cfg := config.Load()
		var resultCh <-chan string
		var errCh <-chan error

		// Route the request to the correct pipe.
		switch strings.TrimPrefix(req.Endpoint, "/") {
		case "workflow":
			pipe, err := pipes.NewWorkflowPipe(ctx, cfg)
			if err != nil {
				s.Emit("pipe-response", SocketResponse{Error: err.Error()})
				return
			}
			resultCh, errCh = pipe.Run(ctx, req.Prompt)
		case "ddgs":
			pipe, err := pipes.NewDDGSPipe(ctx, cfg)
			if err != nil {
				s.Emit("pipe-response", SocketResponse{Error: err.Error()})
				return
			}
			resultCh, errCh = pipe.ARun(ctx, req.Prompt)
		case "serpapi":
			pipe, err := pipes.NewSerpApiPipe(ctx, cfg)
			if err != nil {
				s.Emit("pipe-response", SocketResponse{Error: err.Error()})
				return
			}
			resultCh, errCh = pipe.Run(ctx, req.Prompt)
		default:
			s.Emit("pipe-response", SocketResponse{Error: "Unknown endpoint: " + req.Endpoint})
			return
		}

		// Wait for the pipe to finish and emit the response back to the client.
		select {
		case result := <-resultCh:
			s.Emit("pipe-response", SocketResponse{Result: result})
		case err := <-errCh:
			s.Emit("pipe-response", SocketResponse{Error: err.Error()})
		case <-ctx.Done():
			s.Emit("pipe-response", SocketResponse{Error: "Request timed out"})
		}
	})

	// Handles client disconnections.
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		logger.Info("Socket disconnected: %s, Reason: %s", s.ID(), reason)
	})

	return server
}