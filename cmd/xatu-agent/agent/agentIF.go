package agent

import "context"

type AgentIF interface {
	Run(ctx context.Context, prompt string) (string, error)
}
