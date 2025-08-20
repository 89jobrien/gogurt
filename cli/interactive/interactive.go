package interactive

import (
	"bufio"
	"context"
	"fmt"
	"gogurt/config"
	"gogurt/pipes"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

func Run(cfg *config.Config, documentPath string) {
	// Interactive menus for configuration
	llmProvider := promptForChoice("Choose an LLM Provider:", []string{"ollama", "openai", "azure"})
	splitterProvider := promptForChoice("Choose a Splitter Provider:", []string{"recursive", "markdown", "character"})
	vectorStoreProvider := promptForChoice("Choose a Vector Store Provider:", []string{"simple", "chroma"})

	cfg.LLMProvider = llmProvider
	cfg.SplitterProvider = splitterProvider
	cfg.VectorStoreProvider = vectorStoreProvider

	// Create and run the RAG pipeline with the chosen configurations
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

func promptForChoice(question string, options []string) string {
	fmt.Println(question)
	for i, option := range options {
		fmt.Printf("  %d) %s\n", i+1, option)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter your choice: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(options) {
			fmt.Println("Invalid choice. Please enter a number from the list.")
			continue
		}
		return options[choice-1]
	}
}