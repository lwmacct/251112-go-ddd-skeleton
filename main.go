package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/yourusername/go-ddd-skeleton/internal/commands/api"
	"github.com/yourusername/go-ddd-skeleton/internal/commands/migrate"
	"github.com/yourusername/go-ddd-skeleton/internal/commands/worker"
)

// @title           BBiz KSO API
// @version         1.0
// @description     BBiz KSO API 文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @BasePath  /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description JWT 认证，格式: Bearer {token}

// buildCommands 根据环境变量条件性构建命令列表
func buildCommands() []*cli.Command {
	commands := []*cli.Command{
		api.Command,     // 🟢 API Service - REST API 服务
		worker.Command,  // 🔧 Worker - 后台任务处理器
		migrate.Command, // 🗄️  Migrate - 数据库迁移工具
	}

	if os.Getenv("SHOW_CLI_ITEM") == "1" {
		// 可以在这里添加额外的调试或开发命令
		commands = append([]*cli.Command{}, commands...)
	}

	return commands
}

func main() {
	app := &cli.Command{
		Name:        "go-ddd-skeleton",
		Version:     "1.0.3",
		Usage:       "DDD 架构的 Golang 应用示例",
		Description: `这是一个基于 Domain-Driven Design (DDD) 的 Golang 应用程序。包含用户认证、订单管理等核心功能。`,
		Commands:    buildCommands(),
		Authors: []any{
			map[string]string{
				"name":  "Your Name",
				"email": "your.email@example.com",
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
