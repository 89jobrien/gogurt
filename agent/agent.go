package agent

import "context"

type Agent interface {
	Invoke(ctx context.Context, prompt string) (string, error)
}