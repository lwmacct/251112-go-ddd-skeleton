package repository

import (
	"context"
	"errors"

	"github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/rbac"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/persistence/mapper"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/persistence/model"

	"gorm.io/gorm"
)

// RoleRepo 角色仓储实现
type RoleRepo struct {
	db *gorm.DB
}

// NewRoleRepo 创建角色仓储
func NewRoleRepo(db *gorm.DB) rbac.RoleRepository {
	return &RoleRepo{db: db}
}

// Create 创建角色
func (r *RoleRepo) Create(ctx context.Context, role *rbac.Role) error {
	m := mapper.RoleToModel(role)
	return r.db.WithContext(ctx).Create(m).Error
}

// Update 更新角色
func (r *RoleRepo) Update(ctx context.Context, role *rbac.Role) error {
	m := mapper.RoleToModel(role)
	return r.db.WithContext(ctx).Model(m).Updates(m).Error
}

// Delete 删除角色
func (r *RoleRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Role{}, "id = ?", id).Error
}

// FindByID 根据ID查找角色
func (r *RoleRepo) FindByID(ctx context.Context, id string) (*rbac.Role, error) {
	var m model.Role
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, rbac.ErrRoleNotFound
		}
		return nil, err
	}
	return mapper.RoleToDomain(&m), nil
}

// FindByCode 根据Code查找角色
func (r *RoleRepo) FindByCode(ctx context.Context, code string) (*rbac.Role, error) {
	var m model.Role
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, rbac.ErrRoleNotFound
		}
		return nil, err
	}
	return mapper.RoleToDomain(&m), nil
}

// List 获取角色列表
func (r *RoleRepo) List(ctx context.Context, offset, limit int) ([]*rbac.Role, error) {
	var models []*model.Role
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return mapper.RolesToDomain(models), nil
}

// ExistsByCode 检查角色Code是否存在
func (r *RoleRepo) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Role{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}

// GetUserRoles 获取用户的所有角色
func (r *RoleRepo) GetUserRoles(ctx context.Context, userID string) ([]*rbac.Role, error) {
	var models []*model.Role
	err := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return mapper.RolesToDomain(models), nil
}

// AssignRoleToUser 分配角色给用户
func (r *RoleRepo) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	userRole := &model.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return r.db.WithContext(ctx).Create(userRole).Error
}

// RemoveRoleFromUser 从用户移除角色
func (r *RoleRepo) RemoveRoleFromUser(ctx context.Context, userID, roleID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&model.UserRole{}).Error
}

// GetRoleUsers 获取角色下的所有用户ID
func (r *RoleRepo) GetRoleUsers(ctx context.Context, roleID string) ([]string, error) {
	var userIDs []string
	err := r.db.WithContext(ctx).
		Model(&model.UserRole{}).
		Where("role_id = ?", roleID).
		Pluck("user_id", &userIDs).Error
	return userIDs, err
}
