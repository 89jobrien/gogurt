package cli

import (
	"flag"
	"gogurt/cli/interactive"
	"gogurt/config"
	"gogurt/console"
	"gogurt/vectorstores/chroma"
	"os"
)

var c = console.ConsoleInstance()

func Execute() {
    var (
        interactiveMode = flag.Bool("i", false, "Run in interactive mode")
        ingestMode     = flag.Bool("ingest", false, "Run document ingestion only")
        ragMode        = flag.Bool("rag", false, "Run RAG queries only (requires pre-ingested documents)")
        documentPath   = flag.String("docs", "docs/", "Path to documents directory")
        configPath     = flag.String("config", ".env", "Path to configuration file")
    )
    flag.Parse()

	if *configPath == "" {
		c.Write("Please provide a configuration file path using the -config flag.")
	}

    // Load configuration
    cfg := config.Load()
    if cfg == nil {
		c.Write("Failed to load configuration")
        os.Exit(1)
    }

    var chromaStore *chroma.Store

    switch {
    case *ingestMode:
        c.Write("Starting in ingestion mode")
        interactive.Run(cfg, *documentPath, chromaStore, "ingest")
    case *ragMode:
        c.Write("Starting in RAG query mode")
        interactive.Run(cfg, *documentPath, chromaStore, "rag")
    case *interactiveMode:
        c.Write("Starting in interactive mode")
        interactive.Run(cfg, *documentPath, chromaStore, "interactive")
    default:
        c.Write("Usage:")
        c.Write("  -i              Interactive mode (choose actions from menu)")
        c.Write("  -ingest         Ingest documents only")
        c.Write("  -rag            Run RAG queries only")
        c.Write("  -docs <path>    Document directory path (default: docs/)")
        c.Write("  -config <path>  Configuration file path")
        flag.Usage()
        os.Exit(1)
    }
}
