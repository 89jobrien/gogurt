package interactive

import (
	"bufio"
	"context"
	"fmt"
	"gogurt/internal/config"
	"gogurt/internal/console"
	"gogurt/internal/pipes"
	"gogurt/internal/vectorstores/chroma"
	"os"
	"strconv"
	"strings"
)

// Create a console instance
var c = console.ConsoleInstance()

// Define an interface for RAG functionality
type RAGRunner interface {
	Run(ctx context.Context, prompt string) (string, error)
}

// Updated to handle both ingest and rag workflows
func Run(cfg *config.Config, documentPath string, s *chroma.Store, mode string) {
	if err := configureProviders(cfg); err != nil {
		c.Err("ERROR: Failed to configure providers: %v \n", err)
		return
	}

	switch mode {
	case "ingest":
		runIngestMode(cfg, documentPath, s)
	case "rag":
		runRAGMode(cfg, documentPath, s)
	default:
		runInteractiveMode(cfg, documentPath, s)
	}
}

func runIngestMode(cfg *config.Config, documentPath string, s *chroma.Store) {
	c.Write("\n==================================================================")
	c.Title("\n==================== Document Ingestion Mode =====================\n")

	ingestor, err := pipes.NewIngestPipe(context.Background(), cfg, documentPath)
	if err != nil {
		c.Err("ERROR: Failed to create ingestion pipeline: %v \n", err)
		os.Exit(1)
	}

	c.Info("Starting document ingestion...\n")
	if err := ingestor.Run(context.Background()); err != nil {
		fmt.Printf("ERROR: Ingestion failed: %v \n", err)
		os.Exit(1)
	}

	c.Info("Document ingestion completed successfully\n")

	if cfg.VectorStoreProvider == "chroma" {
		showDBMetrics(cfg, documentPath, s)
	}
}

func runRAGMode(cfg *config.Config, documentPath string, s *chroma.Store) {
	c.Write("\n==================================================================")
	c.Title("\n========================= RAG Mode ===============================\n\n")

	rag, err := pipes.NewRAGPipe(context.Background(), cfg)
	if err != nil {
		c.Err("ERROR: Failed to create RAG query pipeline: %v \n", err)
		os.Exit(1)
	}

	runChatSession(rag, cfg, documentPath, s)
}

func runInteractiveMode(cfg *config.Config, documentPath string, s *chroma.Store) {
	c.Write("\n==================================================================")
	c.Title("\n==================== Interactive Mode ============================\n")

	for {
		topChoice := promptForChoice("\nChoose an action:", []string{"ingest", "rag", "metrics", "exit"})

		switch topChoice {
		case "ingest":
			runIngestMode(cfg, documentPath, s)
		case "rag":
			runRAGMode(cfg, documentPath, s)
		case "metrics":
			showDBMetrics(cfg, documentPath, s)
		case "exit":
			c.Input("Exiting.\n")
			return
		default:
			c.Warn("Unknown option selected\n")
		}
	}
}

func configureProviders(cfg *config.Config) error {

	c.Write("\n==================================================================")
	llmProvider := promptForChoice("\nChoose an LLM Provider:\n", []string{"ollama", "openai", "azure"})
	if llmProvider == "" {
		c.Warn("No LLM provider selected")
	}

	c.Write("\n==================================================================")
	splitterProvider := promptForChoice("\nChoose a Splitter Provider:\n", []string{"recursive", "markdown", "character"})
	if splitterProvider == "" {
		c.Warn("No splitter provider selected")
	}

	c.Write("\n==================================================================")
	vectorStoreProvider := promptForChoice("\nChoose a Vector Store Provider:\n", []string{"simple", "chroma"})
	if vectorStoreProvider == "" {
		c.Warn("No vector store provider selected")
	}

	cfg.LLMProvider = llmProvider
	cfg.SplitterProvider = splitterProvider
	cfg.VectorStoreProvider = vectorStoreProvider

	return nil
}

func runChatSession(rag RAGRunner, cfg *config.Config, documentPath string, s *chroma.Store) {
	c.Write("\n==================================================================")
	c.Title("\n=================== Chat Session Started =========================\n")
	c.Hdr("\nCommands: [ metrics | help | delete-collection | exit ]\n")
	reader := bufio.NewReader(os.Stdin)

	for {
		c.Usr("\n>>> ")
		prompt, err := reader.ReadString('\n')
		if err != nil {
			c.Err("ERROR: Failed to read input: %v \n", err)
			continue
		}

		prompt = strings.TrimSpace(prompt)
		if prompt == "" {
			c.Prompt("Please enter a question or command.")
			continue
		}

		switch strings.ToLower(prompt) {

		case "exit":
			c.Ok("Ending chat session.")
			return

		case "metrics":
			showDBMetrics(cfg, documentPath, s)
			continue

		case "init-collection":
			c.Prompt("Enter the name of the collection to initialize: ")
			collectionName, _ := reader.ReadString('\n')
			collectionName = strings.TrimSpace(collectionName)

			if collectionName == "" {
				c.Warn("Collection name cannot be empty.")
				continue
			}

			initCollection(s, collectionName)
			continue

		case "delete-collection":
			c.Prompt("Enter the name of the collection to delete: ")
			collectionName, _ := reader.ReadString('\n')
			deleteCollection(s, collectionName)
			continue

		case "help":
			showChatHelp()
			continue
		}

		// Process the user's question
		response, err := rag.Run(context.Background(), prompt)
		if err != nil {
			c.Warn("ERROR: Pipeline run failed: %v \n", err)
			c.AI("AI: I'm sorry, I encountered an error processing your question. Please try again.")
			continue
		}

		c.AI("\n🤖 AI: %s\n\n", response)
	}
}

func initCollection(s *chroma.Store, collectionName string) {
	if s == nil {
		c.Warn("Vector store is not available or not a Chroma store.")
		return
	}

	if collectionName == "" {
		c.Warn("No collection name provided.")
		return
	}

	// Trim whitespace from collection name (important since it comes from user input)
	collectionName = strings.TrimSpace(collectionName)

	ctx := context.Background()

	// Create collection with default options - the API appears to use different option types
	col, err := s.Client.CreateCollection(ctx, collectionName)
	if err != nil {
		c.Err("Error creating collection '%s': %v\n", collectionName, err)
		return
	}

	c.Info("Collection '%s' created successfully\n", collectionName)

	// Optionally update the store's current collection reference
	if s.Col == nil {
		s.Col = col
		c.Info("Set '%s' as the current active collection\n", collectionName)
	}
}

func deleteCollection(s *chroma.Store, collection string) {
	if s == nil {
		c.Warn("Vector store is not available or not a Chroma store.")
		return
	}

	ctx := context.Background()
	err := s.Client.DeleteCollection(ctx, collection)
	if err != nil {
		c.Err("Error deleting collection: %v\n", err)
	} else {
		c.Info("Collection deleted successfully\n")
	}
}

func promptForChoice(question string, options []string) string {
	if len(options) == 0 {
		return ""
	}

	c.Hdr("%v", question)
	for i, option := range options {
		c.Sys("  %d) %s\n", i+1, option)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		c.Input("\nEnter your choice (1-%v): ", strconv.Itoa(len(options)))
		input, err := reader.ReadString('\n')
		if err != nil {
			c.Err("ERROR: Error reading input: %v \n", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			c.Warn("Please enter a valid choice.\n")
			continue
		}

		choice, err := strconv.Atoi(input)
		if err != nil {
			c.Warn("Please enter a valid number.\n")
			continue
		}

		if choice < 1 || choice > len(options) {
			c.Warn("Invalid choice. Please enter a number between 1 and %d.\n", len(options))
			continue
		}

		return options[choice-1]
	}
}

func showDBMetrics(cfg *config.Config, documentPath string, s *chroma.Store) {

	if s == nil && cfg.VectorStoreProvider == "chroma" {
		c.Write("\n==================================================================")
		c.Info("\nCreating new Chroma connection for metrics...\n")
		c.Title("\n===================== Database Metrics ===========================\n\n")
		c.Info("Ingestable Docs Path: %v\n", documentPath)
		c.Info("Chroma Database: %v\n", cfg.ChromaDatabase)
		c.Info("Chroma Collection: %v\n", cfg.ChromaCollection)
		c.Info("Chroma Tenant: %v\n", cfg.ChromaTenant)
		c.Info("Chroma URL: %v\n", cfg.ChromaURL)

		newStore, err := chroma.New(cfg)
		if err != nil {
			c.Err("Error creating Chroma connection: %v\n", err)
			return
		}
		s = newStore
	}

	if s == nil {
		c.Warn("Vector store is not available or not a Chroma store.")
		return
	}

	ctx := context.Background()

	colCount, err := s.Client.CountCollections(ctx)
	if err != nil {
		c.Err("Error getting collection count: %v\n", err)
	} else {
		c.Info("Number of collections: %d\n", colCount)
	}

	if s.Col != nil {
		docCount, err := s.Col.Count(ctx)
		if err != nil {
			c.Err("Error getting document count: %v\n", err)
		} else {
			c.Info("Documents in current collection: %d\n", docCount)
		}
	} else {
		c.Warn("No collection available\n")
	}

	c.Write("\n==================================================================")
}

func showChatHelp() {
	c.Write("\n==================================================================")
	c.Title("\n===================== Chat Help ==================================\n")
	c.Hdr("\nAvailable commands:\n")
	c.Info("  metrics - Show database metrics\n")
	c.Info("  help    - Show this help message\n")
	c.Info("  back    - Return to main menu\n")
	c.Info("  exit    - End the chat session\n")
	c.Info("  Or ask any question about your documents\n")
	c.Write("\n==================================================================")
}
