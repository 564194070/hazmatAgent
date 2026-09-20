package commandcli

import (
	"agentFrame/cmd/xatu-agent/agent"
	"agentFrame/cmd/xatu-agent/llm"
	"context"
	"log/slog"

	"github.com/urfave/cli/v3"
)

func CommandCliAction(ctx context.Context, c *cli.Command) error {

	slog.Info("启动综合智能体，命令行客户端模式")

	prompt := c.String("prompt")
	slog.Info("用户提问", "prompt", prompt)
	reActAgent := agent.NewReActAgent(llm.NewOpenAILLMClient())
	res, err := reActAgent.Run(ctx, prompt)
	if err != nil {
		slog.Error("执行智能体失败", "error", err)
		return err
	}
	slog.Info("智能体回答", "res", res)
	return nil
}
