package pipes

import (
	"context"

	"gogurt/internal/console"
)

type Pipe interface {
	Run(ctx context.Context, prompt string) (string, error)
	ARun(ctx context.Context, prompt string) (<-chan string, <-chan error)
}

var c = console.ConsoleInstance()
