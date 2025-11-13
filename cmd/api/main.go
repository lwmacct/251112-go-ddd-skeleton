package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/yourusername/go-ddd-skeleton/internal/commands/api"
)

func main() {
	cmd := &cli.Command{
		Name:        "api",
		Version:     "1.0.0",
		Usage:       "REST API 服务",
		Description: "启动 HTTP API 服务器",
		Commands: []*cli.Command{
			api.Command,
		},
		// 默认执行 API 命令
		Action: api.Command.Action,
		Flags:  api.Command.Flags,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

