package repository

import (
	"context"
	"errors"

	"github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/rbac"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/persistence/mapper"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/persistence/model"

	"gorm.io/gorm"
)

// MenuRepo 菜单仓储实现
type MenuRepo struct {
	db *gorm.DB
}

// NewMenuRepo 创建菜单仓储
func NewMenuRepo(db *gorm.DB) rbac.MenuRepository {
	return &MenuRepo{db: db}
}

// Create 创建菜单
func (r *MenuRepo) Create(ctx context.Context, menu *rbac.Menu) error {
	m := mapper.MenuToModel(menu)
	return r.db.WithContext(ctx).Create(m).Error
}

// Update 更新菜单
func (r *MenuRepo) Update(ctx context.Context, menu *rbac.Menu) error {
	m := mapper.MenuToModel(menu)
	return r.db.WithContext(ctx).Model(m).Updates(m).Error
}

// Delete 删除菜单
func (r *MenuRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Menu{}, "id = ?", id).Error
}

// FindByID 根据ID查找菜单
func (r *MenuRepo) FindByID(ctx context.Context, id string) (*rbac.Menu, error) {
	var m model.Menu
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, rbac.ErrMenuNotFound
		}
		return nil, err
	}
	return mapper.MenuToDomain(&m), nil
}

// List 获取菜单列表（可过滤父菜单）
func (r *MenuRepo) List(ctx context.Context, parentID *string, offset, limit int) ([]*rbac.Menu, error) {
	query := r.db.WithContext(ctx).Order("sort_order ASC, id ASC")

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	var models []*model.Menu
	err := query.Offset(offset).Limit(limit).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return mapper.MenusToDomain(models), nil
}

// GetAll 获取所有菜单（用于构建树）
func (r *MenuRepo) GetAll(ctx context.Context) ([]*rbac.Menu, error) {
	var models []*model.Menu
	err := r.db.WithContext(ctx).Order("sort_order ASC, id ASC").Find(&models).Error
	if err != nil {
		return nil, err
	}
	return mapper.MenusToDomain(models), nil
}

// GetChildren 获取指定菜单的子菜单
func (r *MenuRepo) GetChildren(ctx context.Context, parentID string) ([]*rbac.Menu, error) {
	var models []*model.Menu
	err := r.db.WithContext(ctx).
		Where("parent_id = ?", parentID).
		Order("sort_order ASC, id ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return mapper.MenusToDomain(models), nil
}

// GetRoleMenus 获取角色的所有菜单
func (r *MenuRepo) GetRoleMenus(ctx context.Context, roleID string) ([]*rbac.Menu, error) {
	var models []*model.Menu
	err := r.db.WithContext(ctx).
		Joins("JOIN role_menus ON role_menus.menu_id = menus.id").
		Where("role_menus.role_id = ?", roleID).
		Order("menus.sort_order ASC, menus.id ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return mapper.MenusToDomain(models), nil
}

// GetUserMenus 获取用户的所有菜单（通过角色，去重）
func (r *MenuRepo) GetUserMenus(ctx context.Context, userID string) ([]*rbac.Menu, error) {
	var models []*model.Menu
	err := r.db.WithContext(ctx).
		Distinct("menus.*").
		Joins("JOIN role_menus ON role_menus.menu_id = menus.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_menus.role_id").
		Where("user_roles.user_id = ? AND menus.is_visible = ?", userID, true).
		Order("menus.sort_order ASC, menus.id ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return mapper.MenusToDomain(models), nil
}

// AssignMenusToRole 分配菜单给角色
func (r *MenuRepo) AssignMenusToRole(ctx context.Context, roleID string, menuIDs []string) error {
	roleMenus := make([]model.RoleMenu, 0, len(menuIDs))
	for _, menuID := range menuIDs {
		roleMenus = append(roleMenus, model.RoleMenu{
			RoleID: roleID,
			MenuID: menuID,
		})
	}
	return r.db.WithContext(ctx).Create(&roleMenus).Error
}

// RemoveMenusFromRole 从角色移除菜单
func (r *MenuRepo) RemoveMenusFromRole(ctx context.Context, roleID string, menuIDs []string) error {
	return r.db.WithContext(ctx).
		Where("role_id = ? AND menu_id IN ?", roleID, menuIDs).
		Delete(&model.RoleMenu{}).Error
}

// SyncRoleMenus 同步角色的菜单（删除旧的，添加新的）
func (r *MenuRepo) SyncRoleMenus(ctx context.Context, roleID string, menuIDs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除旧的关联
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error; err != nil {
			return err
		}

		// 如果没有新的菜单，直接返回
		if len(menuIDs) == 0 {
			return nil
		}

		// 添加新的关联
		roleMenus := make([]model.RoleMenu, 0, len(menuIDs))
		for _, menuID := range menuIDs {
			roleMenus = append(roleMenus, model.RoleMenu{
				RoleID: roleID,
				MenuID: menuID,
			})
		}
		return tx.Create(&roleMenus).Error
	})
}
