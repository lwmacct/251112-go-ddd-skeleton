package user

import "context"

// Repository 用户仓储接口
type Repository interface {
	// Create 创建用户
	Create(ctx context.Context, user *User) error

	// Update 更新用户
	Update(ctx context.Context, user *User) error

	// Delete 删除用户
	Delete(ctx context.Context, id string) error

	// FindByID 根据ID查找用户
	FindByID(ctx context.Context, id string) (*User, error)

	// FindByEmail 根据邮箱查找用户
	FindByEmail(ctx context.Context, email Email) (*User, error)

	// List 列出用户（分页）
	List(ctx context.Context, offset, limit int) ([]*User, int64, error)

	// ExistsByEmail 检查邮箱是否已存在
	ExistsByEmail(ctx context.Context, email Email) (bool, error)
}
