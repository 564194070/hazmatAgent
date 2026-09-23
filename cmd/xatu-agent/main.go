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
		Name:  "xatu-agent",
		Usage: "启动xatu-agent智能体",
		Commands: []*cli.Command{
			{
				Name:   "cmdcli",
				Usage:  "命令行交互模式",
				Action: commandcli.CommandCliAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "prompt",
						Aliases: []string{"p"},
						Usage:   "提问文字",
					},
				},
			},
			{
				Name:   "server",
				Usage:  "服务化交互模式",
				Action: commandcli.ServerAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Usage:   "端口",
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		slog.Error("程序执行失败", "err", err)
		os.Exit(1)
	}
}
