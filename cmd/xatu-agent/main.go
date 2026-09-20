package main

import (
	commandcli "agentFrame/cmd/xatu-agent/actions"
	"context"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	// 加载 .env
	_ = godotenv.Load()

	cmd := &cli.Command{
		Name:  "commandCli",
		Usage: "调用综合智能体处理请求",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "prompt",
				Aliases: []string{"p"},
				Usage:   "提问文字",
			},
		},
		Action: commandcli.CommandCliAction,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		slog.Error("程序执行失败", "err", err)
		os.Exit(1)
	}
}
