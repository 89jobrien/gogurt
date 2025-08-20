package pipes

import "context"

type Pipe interface {
	Run(ctx context.Context, prompt string) (string, error)
}