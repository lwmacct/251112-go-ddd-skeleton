package migrate

import (
	"context"
	"log"

	"github.com/urfave/cli/v3"
	"github.com/yourusername/go-ddd-skeleton/internal/config"
	"github.com/yourusername/go-ddd-skeleton/internal/infrastructure/persistence"
	"github.com/yourusername/go-ddd-skeleton/internal/infrastructure/persistence/model"
)

// Command 定义数据库迁移命令
var Command = &cli.Command{
	Name:    "migrate",
	Aliases: []string{"m", "migration"},
	Usage:   "数据库迁移工具",
	Description: `
   执行数据库结构迁移，自动创建或更新数据库表结构。
   使用 GORM 的 AutoMigrate 功能进行迁移。
	`,
	Commands: []*cli.Command{
		{
			Name:    "up",
			Aliases: []string{"u"},
			Usage:   "执行数据库迁移（创建/更新表）",
			Action:  runMigrationUp,
		},
		{
			Name:    "status",
			Aliases: []string{"s"},
			Usage:   "查看数据库迁移状态",
			Action:  runMigrationStatus,
		},
	},
	// 默认执行 up 子命令
	Action: runMigrationUp,
}

// runMigrationUp 执行数据库迁移
func runMigrationUp(ctx context.Context, cmd *cli.Command) error {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Connecting to database: %s@%s:%d/%s",
		cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	// 连接数据库
	db, err := persistence.NewPostgres(persistence.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移
	log.Println("Running database migrations...")
	if err := model.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("✅ Migrations completed successfully!")
	return nil
}

// runMigrationStatus 查看迁移状态
func runMigrationStatus(ctx context.Context, cmd *cli.Command) error {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 连接数据库
	_, err = persistence.NewPostgres(persistence.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("✅ Database connection successful!")
	log.Println("Checking migration status...")

	// 获取所有模型
	models := model.AllModels()
	log.Printf("Total models to migrate: %d", len(models))

	for _, m := range models {
		// 这里可以检查表是否存在
		log.Printf("- Model: %T", m)
	}

	log.Println("Status check completed!")
	return nil
}
