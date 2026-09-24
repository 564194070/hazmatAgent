package agent

import "context"

type AgentIF interface {
	Run(ctx context.Context, userID, sessionID, prompt string) (string, error)
}
