package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/yourusername/go-ddd-skeleton/internal/commands/migrate"
)

func main() {
	cmd := &cli.Command{
		Name:        "migrate",
		Version:     "1.0.0",
		Usage:       "数据库迁移工具",
		Description: "执行数据库迁移操作",
		Commands:    migrate.Command.Commands,
		// 默认执行 migrate up
		Action: migrate.Command.Action,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

