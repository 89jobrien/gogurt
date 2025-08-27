package pipes

import (
	"context"

	"gogurt/internal/console"
)

// Pipe defines the interface for a non-blocking, asynchronous pipeline.
type Pipe interface {
	// Run executes the pipe's workflow asynchronously, returning channels for the final result and any errors.
	Run(ctx context.Context, prompt string) (<-chan string, <-chan error)
}

var c = console.ConsoleInstance()