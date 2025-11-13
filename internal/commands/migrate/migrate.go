package migrate

import (
	"context"
	"log"

	"github.com/lwmacct/251112-go-ddd-skeleton/internal/config"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/auth"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/persistence"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/persistence/model"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/seed"
	"github.com/urfave/cli/v3"
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
		{
			Name:    "seed",
			Aliases: []string{"sd"},
			Usage:   "执行数据库种子数据初始化",
			Description: `
   初始化数据库种子数据，包括：
   - 默认角色（admin, user, editor, viewer）
   - 默认权限（user:*, role:*, menu:*, order:*）
   - 默认菜单（系统管理、订单管理、个人中心）
   - 默认管理员账户（admin@example.com）

   该命令支持幂等性，可以安全地多次执行。
			`,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "force",
					Aliases: []string{"f"},
					Usage:   "强制重新执行（忽略历史记录）",
				},
			},
			Action: runSeed,
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

// runSeed 执行数据库种子数据初始化
func runSeed(ctx context.Context, cmd *cli.Command) error {
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

	// 创建 seed 管理器
	manager := seed.NewManager(db)

	// 创建密码哈希器
	passwordHasher := auth.NewPasswordHasher()

	// 注册 seeders（顺序很重要：先 RBAC，后 User）
	manager.Register(seed.NewRBACSeeder())
	manager.Register(seed.NewUserSeeder(passwordHasher))

	// 执行所有 seed
	if err := manager.RunAll(ctx); err != nil {
		log.Fatalf("❌ Seed failed: %v", err)
	}

	// 打印成功信息
	log.Println("\n╔══════════════════════════════════════════════════════╗")
	log.Println("║                   🎉 Seed 完成！                     ║")
	log.Println("╚══════════════════════════════════════════════════════╝")
	log.Println("\n📋 默认管理员账户:")
	log.Println("   Email:    admin@example.com")
	log.Println("   Password: Admin@123456")
	log.Println("\n⚠️  请在首次登录后立即修改默认密码！")
	log.Println("\n💡 提示: 你可以通过以下方式登录测试:")
	log.Println("   curl -X POST http://localhost:8080/api/auth/login \\")
	log.Println("     -H \"Content-Type: application/json\" \\")
	log.Println("     -d '{\"email\":\"admin@example.com\",\"password\":\"Admin@123456\"}'")
	log.Println()

	return nil
}
