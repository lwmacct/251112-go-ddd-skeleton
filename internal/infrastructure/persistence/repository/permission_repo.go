package repository

import (
	"context"
	"errors"

	"github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/rbac"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/persistence/mapper"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/persistence/model"

	"gorm.io/gorm"
)

// PermissionRepo 权限仓储实现
type PermissionRepo struct {
	db *gorm.DB
}

// NewPermissionRepo 创建权限仓储
func NewPermissionRepo(db *gorm.DB) rbac.PermissionRepository {
	return &PermissionRepo{db: db}
}

// Create 创建权限
func (r *PermissionRepo) Create(ctx context.Context, permission *rbac.Permission) error {
	m := mapper.PermissionToModel(permission)
	return r.db.WithContext(ctx).Create(m).Error
}

// Update 更新权限
func (r *PermissionRepo) Update(ctx context.Context, permission *rbac.Permission) error {
	m := mapper.PermissionToModel(permission)
	return r.db.WithContext(ctx).Model(m).Updates(m).Error
}

// Delete 删除权限
func (r *PermissionRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Permission{}, "id = ?", id).Error
}

// FindByID 根据ID查找权限
func (r *PermissionRepo) FindByID(ctx context.Context, id string) (*rbac.Permission, error) {
	var m model.Permission
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, rbac.ErrPermissionNotFound
		}
		return nil, err
	}
	return mapper.PermissionToDomain(&m), nil
}

// FindByCode 根据Code查找权限
func (r *PermissionRepo) FindByCode(ctx context.Context, code string) (*rbac.Permission, error) {
	var m model.Permission
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, rbac.ErrPermissionNotFound
		}
		return nil, err
	}
	return mapper.PermissionToDomain(&m), nil
}

// List 获取权限列表
func (r *PermissionRepo) List(ctx context.Context, offset, limit int) ([]*rbac.Permission, error) {
	var models []*model.Permission
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return mapper.PermissionsToDomain(models), nil
}

// ExistsByCode 检查权限Code是否存在
func (r *PermissionRepo) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Permission{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}

// GetRolePermissions 获取角色的所有权限
func (r *PermissionRepo) GetRolePermissions(ctx context.Context, roleID string) ([]*rbac.Permission, error) {
	var models []*model.Permission
	err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return mapper.PermissionsToDomain(models), nil
}

// AssignPermissionsToRole 分配权限给角色
func (r *PermissionRepo) AssignPermissionsToRole(ctx context.Context, roleID string, permissionIDs []string) error {
	rolePermissions := make([]model.RolePermission, 0, len(permissionIDs))
	for _, permID := range permissionIDs {
		rolePermissions = append(rolePermissions, model.RolePermission{
			RoleID:       roleID,
			PermissionID: permID,
		})
	}
	return r.db.WithContext(ctx).Create(&rolePermissions).Error
}

// RemovePermissionsFromRole 从角色移除权限
func (r *PermissionRepo) RemovePermissionsFromRole(ctx context.Context, roleID string, permissionIDs []string) error {
	return r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id IN ?", roleID, permissionIDs).
		Delete(&model.RolePermission{}).Error
}

// GetUserPermissions 获取用户的所有权限（通过角色）
func (r *PermissionRepo) GetUserPermissions(ctx context.Context, userID string) ([]*rbac.Permission, error) {
	var models []*model.Permission
	err := r.db.WithContext(ctx).
		Distinct("permissions.*").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return mapper.PermissionsToDomain(models), nil
}
