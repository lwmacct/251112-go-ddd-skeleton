package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpserver "github.com/lwmacct/251112-go-ddd-skeleton/internal/adapters/http"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/bootstrap"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/config"
	"github.com/urfave/cli/v3"
)

// Command 定义 API 服务命令
var Command = &cli.Command{
	Name:    "api",
	Aliases: []string{"serve", "server"},
	Usage:   "启动 REST API 服务",
	Description: `
   启动 HTTP API 服务器，提供 RESTful API 接口。
   服务器支持优雅关闭，会等待正在处理的请求完成后再退出。
	`,
	Action: runAPIServer,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "host",
			Aliases: []string{"H"},
			Usage:   "服务器监听地址",
			Value:   "0.0.0.0",
		},
		&cli.IntFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Usage:   "服务器监听端口",
			Value:   8080,
		},
	},
}

// runAPIServer 执行 API 服务器启动逻辑
func runAPIServer(ctx context.Context, cmd *cli.Command) error {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 如果命令行指定了 host 和 port，覆盖配置文件
	if cmd.IsSet("host") {
		cfg.Server.Host = cmd.String("host")
	}
	if cmd.IsSet("port") {
		cfg.Server.Port = cmd.Int("port")
	}

	// 初始化容器（依赖注入）
	container, err := bootstrap.NewContainer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}

	// 创建HTTP服务器
	server := httpserver.NewServer(container.Router, cfg.Server.Host, cfg.Server.Port)

	// 启动服务器（在goroutine中）
	go func() {
		log.Printf("Starting API server on %s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down API server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("API server exited")
	return nil
}
