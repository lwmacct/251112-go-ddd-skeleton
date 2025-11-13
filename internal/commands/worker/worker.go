package worker

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli/v3"
	"github.com/yourusername/go-ddd-skeleton/internal/config"
)

// Command 定义后台任务处理器命令
var Command = &cli.Command{
	Name:    "worker",
	Aliases: []string{"w"},
	Usage:   "启动后台任务处理器",
	Description: `
   启动后台 Worker 进程，用于处理异步任务。
   例如：清理过期会话、发送通知邮件、处理订单状态等。
	`,
	Action: runWorker,
	Flags: []cli.Flag{
		&cli.DurationFlag{
			Name:    "interval",
			Aliases: []string{"i"},
			Usage:   "任务执行间隔",
			Value:   1 * time.Minute,
		},
	},
}

// runWorker 执行 Worker 逻辑
func runWorker(ctx context.Context, cmd *cli.Command) error {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	interval := cmd.Duration("interval")
	log.Printf("Starting worker in %s environment (interval: %v)...", cfg.App.Env, interval)

	// 创建context用于优雅关闭
	workerCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动后台任务
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.Println("Worker task executing...")
				// TODO: 这里可以添加定期任务，例如：
				// - 清理过期会话
				// - 发送通知邮件
				// - 处理订单状态
				executeWorkerTasks(workerCtx, cfg)
			case <-workerCtx.Done():
				log.Println("Worker stopping...")
				return
			}
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")
	cancel()

	// 给一些时间让任务完成
	time.Sleep(2 * time.Second)
	log.Println("Worker exited")

	return nil
}

// executeWorkerTasks 执行具体的后台任务
func executeWorkerTasks(ctx context.Context, cfg *config.Config) {
	// 在这里实现具体的后台任务逻辑
	// 例如：
	// 1. 清理过期的会话和令牌
	// 2. 发送待发送的邮件
	// 3. 处理超时的订单
	// 4. 生成统计报表

	log.Println("Executing background tasks...")
}
