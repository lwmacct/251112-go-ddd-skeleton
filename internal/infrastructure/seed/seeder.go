package seed

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// SeedHistory 记录 seed 执行历史
type SeedHistory struct {
	ID         string    `gorm:"primaryKey;type:varchar(26)"`
	Name       string    `gorm:"uniqueIndex;type:varchar(100)"` // seed 名称
	Version    string    `gorm:"type:varchar(50)"`              // 版本号
	Status     string    `gorm:"type:varchar(20)"`              // success/failed
	Error      string    `gorm:"type:text"`
	ExecutedAt time.Time `gorm:"autoCreateTime"`
}

// TableName 指定表名
func (SeedHistory) TableName() string {
	return "seed_history"
}

// Seeder 接口定义 seed 的标准行为
type Seeder interface {
	// Name 返回 seed 的唯一名称
	Name() string
	// Run 执行 seed 逻辑
	Run(ctx context.Context, db *gorm.DB) error
	// ShouldRun 检查是否应该执行（幂等性检查）
	ShouldRun(ctx context.Context, db *gorm.DB) (bool, error)
}

// Manager seed 管理器，负责协调多个 seeder 的执行
type Manager struct {
	db      *gorm.DB
	seeders []Seeder
}

// NewManager 创建 seed 管理器
func NewManager(db *gorm.DB) *Manager {
	return &Manager{
		db:      db,
		seeders: []Seeder{},
	}
}

// Register 注册一个 seeder
func (m *Manager) Register(seeder Seeder) {
	m.seeders = append(m.seeders, seeder)
}

// RunAll 执行所有已注册的 seeder
// 每个 seeder 在独立的事务中执行，失败时自动回滚
func (m *Manager) RunAll(ctx context.Context) error {
	// 确保 seed_history 表存在
	if err := m.db.AutoMigrate(&SeedHistory{}); err != nil {
		return fmt.Errorf("failed to create seed_history table: %w", err)
	}

	log.Println("🌱 Running database seeders...")

	// 遍历执行每个 seeder
	for _, seeder := range m.seeders {
		// 检查是否应该执行
		shouldRun, err := seeder.ShouldRun(ctx, m.db)
		if err != nil {
			return fmt.Errorf("failed to check if seeder should run: %w", err)
		}

		if !shouldRun {
			log.Printf("⏭️  Skipping seed: %s (already executed)", seeder.Name())
			continue
		}

		log.Printf("▶️  Running seed: %s", seeder.Name())

		// 在事务中执行 seed
		err = m.db.Transaction(func(tx *gorm.DB) error {
			if err := seeder.Run(ctx, tx); err != nil {
				// 记录失败
				m.recordSeedHistory(tx, seeder.Name(), "failed", err.Error())
				return err
			}

			// 记录成功
			m.recordSeedHistory(tx, seeder.Name(), "success", "")
			return nil
		})

		if err != nil {
			log.Printf("❌ Seed failed: %s - %v", seeder.Name(), err)
			return fmt.Errorf("seed %s failed: %w", seeder.Name(), err)
		}

		log.Printf("✅ Seed completed: %s", seeder.Name())
	}

	log.Println("\n✅ All seeds completed successfully!")
	return nil
}

// recordSeedHistory 记录 seed 执行历史
func (m *Manager) recordSeedHistory(tx *gorm.DB, name, status, errMsg string) {
	history := &SeedHistory{
		ID:      ulid.Make().String(),
		Name:    name,
		Version: "1.0",
		Status:  status,
		Error:   errMsg,
	}
	tx.Create(history)
}

// GetHistory 获取 seed 执行历史
func (m *Manager) GetHistory(ctx context.Context) ([]SeedHistory, error) {
	var history []SeedHistory
	err := m.db.WithContext(ctx).
		Order("executed_at DESC").
		Find(&history).Error
	return history, err
}
