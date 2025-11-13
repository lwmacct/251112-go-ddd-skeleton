package repository

import (
	"context"
	"errors"

	"github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/user"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/persistence/mapper"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/persistence/model"
	"gorm.io/gorm"
)

// UserRepository 用户仓储实现
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) user.Repository {
	return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	m := mapper.UserToModel(u)
	return r.db.WithContext(ctx).Create(m).Error
}

// Update 更新用户
func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	m := mapper.UserToModel(u)
	return r.db.WithContext(ctx).Save(m).Error
}

// Delete 删除用户
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id).Error
}

// FindByID 根据ID查找用户
func (r *UserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	var m model.User
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return mapper.UserToDomain(&m)
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepository) FindByEmail(ctx context.Context, email user.Email) (*user.User, error) {
	var m model.User
	if err := r.db.WithContext(ctx).First(&m, "email = ?", email.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return mapper.UserToDomain(&m)
}

// List 列出用户
func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]*user.User, int64, error) {
	var models []model.User
	var total int64

	// 查询总数
	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询列表
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	users := make([]*user.User, len(models))
	for i, m := range models {
		u, err := mapper.UserToDomain(&m)
		if err != nil {
			return nil, 0, err
		}
		users[i] = u
	}

	return users, total, nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *UserRepository) ExistsByEmail(ctx context.Context, email user.Email) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email.String()).Count(&count).Error
	return count > 0, err
}
