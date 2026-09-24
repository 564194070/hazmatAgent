package actions

import (
	"agentFrame/cmd/xatu-agent/agent"
	"agentFrame/cmd/xatu-agent/llm"
	"agentFrame/cmd/xatu-agent/memory"
	"agentFrame/cmd/xatu-agent/memory/manager"
	"agentFrame/cmd/xatu-agent/tools"
	"context"
	"log/slog"

	"github.com/urfave/cli/v3"
)

func CommandCliAction(ctx context.Context, c *cli.Command) error {

	slog.Info("启动综合智能体，命令行客户端模式")

	prompt := c.String("prompt")
	slog.Info("用户提问", "prompt", prompt)

	memoryManager, err := manager.NewMemoryStoreManager([]memory.MemoryStoreType{memory.MemoryStoreTypeMySQL}, nil)
	if err != nil {
		slog.Error("创建记忆管理器失败", "error", err)
		return err
	}

	reActAgent := agent.NewReActAgent(llm.NewOpenAILLMClient(), memoryManager)
	registry := tools.NewToolRegistry()
	registry.RegisterTool(tools.NewEchoTool())
	registry.RegisterTool(tools.NewMoveFilesTool())
	registry.RegisterTool(tools.NewReadFileListTool())
	registry.RegisterTool(tools.NewMkdirAllTool())
	reActAgent.SetTools(registry)
	res, err := reActAgent.Run(ctx, "1", "1001", prompt)
	if err != nil {
		slog.Error("执行智能体失败", "error", err)
		return err
	}
	slog.Info("智能体回答", "res", res)
	return nil
}
