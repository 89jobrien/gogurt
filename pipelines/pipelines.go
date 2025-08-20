package pipelines

import "context"

type Pipeline interface {
	Run(ctx context.Context, prompt string) (string, error)
}