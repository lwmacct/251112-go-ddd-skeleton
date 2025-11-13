package main

import (
	"context"
	"log"
	"os"

	"github.com/lwmacct/251112-go-ddd-skeleton/internal/commands/worker"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:        "worker",
		Version:     "1.0.0",
		Usage:       "后台任务处理器",
		Description: "启动后台 Worker 进程",
		Commands: []*cli.Command{
			worker.Command,
		},
		// 默认执行 Worker 命令
		Action: worker.Command.Action,
		Flags:  worker.Command.Flags,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
