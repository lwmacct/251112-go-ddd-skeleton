package seed

import (
	"context"
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
	"gorm.io/gorm"

	"github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/user"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/auth"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/persistence/mapper"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/persistence/model"
)

// UserData 用户数据结构
type UserData struct {
	Email    string   `yaml:"email"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	IsActive bool     `yaml:"is_active"`
	Roles    []string `yaml:"roles"`
}

// UserSeeder 用户 seed 实现
type UserSeeder struct {
	passwordHasher *auth.PasswordHasher
}

// NewUserSeeder 创建用户 seeder
func NewUserSeeder(passwordHasher *auth.PasswordHasher) *UserSeeder {
	return &UserSeeder{
		passwordHasher: passwordHasher,
	}
}

// Name 返回 seeder 名称
func (s *UserSeeder) Name() string {
	return "DEFAULT_ADMIN_USER"
}

// ShouldRun 检查是否应该执行
func (s *UserSeeder) ShouldRun(ctx context.Context, db *gorm.DB) (bool, error) {
	// 检查是否已有 admin@example.com 用户
	var count int64
	err := db.Model(&model.User{}).
		Where("email = ?", "admin@example.com").
		Count(&count).Error
	return count == 0, err
}

// Run 执行 seed
func (s *UserSeeder) Run(ctx context.Context, db *gorm.DB) error {
	log.Println("  📦 Loading user data...")

	// 加载用户数据
	usersData, err := s.loadUsers()
	if err != nil {
		return fmt.Errorf("failed to load users: %w", err)
	}

	for _, uData := range usersData {
		// 哈希密码
		hashedPassword, err := s.passwordHasher.Hash(uData.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password for %s: %w", uData.Email, err)
		}

		// 创建用户实体
		u, err := user.NewUser(uData.Email, hashedPassword, uData.Username)
		if err != nil {
			return fmt.Errorf("failed to create user %s: %w", uData.Email, err)
		}
		u.ID = generateULID()
		u.IsActive = uData.IsActive

		// 转换为 GORM 模型并插入
		userModel := mapper.UserToModel(u)
		if err := db.Create(userModel).Error; err != nil {
			return fmt.Errorf("failed to insert user %s: %w", uData.Email, err)
		}

		log.Printf("  ✓ Created user: %s", uData.Email)

		// 分配角色
		for _, roleCode := range uData.Roles {
			var role model.Role
			if err := db.Where("code = ?", roleCode).First(&role).Error; err != nil {
				return fmt.Errorf("role not found: %s", roleCode)
			}

			userRole := model.UserRole{
				UserID: u.ID,
				RoleID: role.ID,
			}
			if err := db.Create(&userRole).Error; err != nil {
				return fmt.Errorf("failed to assign role %s to user %s: %w", roleCode, uData.Email, err)
			}

			log.Printf("  ✓ Assigned role: %s", roleCode)
		}
	}

	log.Printf("  ✅ Created %d default user(s)", len(usersData))
	return nil
}

// loadUsers 从 YAML 加载用户数据
func (s *UserSeeder) loadUsers() ([]UserData, error) {
	data, err := seedDataFS.ReadFile("data/users.yaml")
	if err != nil {
		return nil, err
	}

	var result struct {
		Users []UserData `yaml:"users"`
	}

	err = yaml.Unmarshal(data, &result)
	return result.Users, err
}
