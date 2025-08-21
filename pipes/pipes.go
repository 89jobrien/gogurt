package pipes

import (
	"context"

	"gogurt/console"
)

type Pipe interface {
	Run(ctx context.Context, prompt string) (string, error)
}

var c = console.ConsoleInstance()