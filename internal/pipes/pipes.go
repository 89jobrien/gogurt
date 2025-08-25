package pipes

import (
	"context"

	"gogurt/internal/console"
)

type Pipe interface {
	Run(ctx context.Context, prompt string) (string, error)
}

var c = console.ConsoleInstance()