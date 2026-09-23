package actions

import (
	"agentFrame/cmd/xatu-agent/web"
	"context"
	"log/slog"

	"github.com/urfave/cli/v3"
)

func ServerAction(ctx context.Context, c *cli.Command) error {

	slog.Info("启动xatu-agent智能体,服务化交互模式")

	return web.RunHttpServer()
}
