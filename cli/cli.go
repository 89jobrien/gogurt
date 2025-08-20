package cli

import (
	"bufio"
	"context"
	"flag"
	"gogurt/cli/interactive"
	"gogurt/config"
	"gogurt/pipes"
	"log/slog"
	"os"
	"strings"
)

func Execute() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	var interactiveFlag = flag.Bool("i", false, "Enable interactive mode to select providers.")
	flag.Parse()

	var documentPath string
	if len(flag.Args()) >= 1 {
		documentPath = flag.Args()[0]
		slog.Info("Loading document from command-line argument", "path", documentPath)
	} else {
		documentPath = "docs/"
		slog.Info("No document path provided, loading from default 'docs/' directory")
	}

	cfg := config.Load()

	if *interactiveFlag {
		slog.Info("Starting in interactive mode.")
		interactive.Run(cfg, documentPath)
	} else {
		slog.Info("Starting in default mode.")

		rag, err := pipes.NewRAG(context.Background(), cfg, documentPath)
		if err != nil {
			slog.Error("failed to create RAG pipeline", "error", err)
			os.Exit(1)
		}
		slog.Info("Chat session started. Type 'exit' to end.")
		reader := bufio.NewReader(os.Stdin)
		for {
			os.Stdout.WriteString("You: ")
			prompt, err := reader.ReadString('\n')
			if err != nil {
				slog.Error("failed to read input", "error", err)
				break
			}
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
}