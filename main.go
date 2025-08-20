package main

import (
	"bufio"
	"context"
	"gogurt/config"
	"gogurt/pipes"
	"log/slog"
	"os"
	"strings"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	var documentPath string
	if len(os.Args) >= 2 {
		documentPath = os.Args[1]
		slog.Info("Loading document from command-line argument", "path", documentPath)
	} else {
		documentPath = "docs/"
		slog.Info("No document path provided, loading from default 'docs/' directory")
	}

	cfg := config.Load()

	var rag pipes.Pipe
	rag, err := pipes.NewRAG(context.Background(), cfg, documentPath)
	if err != nil {
		slog.Error("failed to create RAG pipeline", "error", err)
		os.Exit(1)
	}

	slog.Info("Chat session started. Type 'exit' to end.")
	reader := bufio.NewReader(os.Stdin)

	for {
		os.Stdout.WriteString("You: ")
		prompt, _ := reader.ReadString('\n')
		prompt = strings.TrimSpace(prompt)

		if strings.ToLower(prompt) == "exit" {
			slog.Info("Ending chat session.")
			break
		}

		response, err := rag.Run(context.Background(), prompt)
		if err != nil {
			slog.Error("pipeline run failed", "error", err)
			continue
		}

		os.Stdout.WriteString("AI: " + response + "\n")
	}
}