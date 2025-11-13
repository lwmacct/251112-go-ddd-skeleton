package rbac

import "context"

// RoleRepository 角色仓储接口
type RoleRepository interface {
	// Create 创建角色
	Create(ctx context.Context, role *Role) error

	// Update 更新角色
	Update(ctx context.Context, role *Role) error

	// Delete 删除角色
	Delete(ctx context.Context, id string) error

	// FindByID 根据ID查找角色
	FindByID(ctx context.Context, id string) (*Role, error)

	// FindByCode 根据Code查找角色
	FindByCode(ctx context.Context, code string) (*Role, error)

	// List 获取角色列表
	List(ctx context.Context, offset, limit int) ([]*Role, error)

	// ExistsByCode 检查角色Code是否存在
	ExistsByCode(ctx context.Context, code string) (bool, error)

	// GetUserRoles 获取用户的所有角色
	GetUserRoles(ctx context.Context, userID string) ([]*Role, error)

	// AssignRoleToUser 分配角色给用户
	AssignRoleToUser(ctx context.Context, userID, roleID string) error

	// RemoveRoleFromUser 从用户移除角色
	RemoveRoleFromUser(ctx context.Context, userID, roleID string) error

	// GetRoleUsers 获取角色下的所有用户ID
	GetRoleUsers(ctx context.Context, roleID string) ([]string, error)
}

// PermissionRepository 权限仓储接口
type PermissionRepository interface {
	// Create 创建权限
	Create(ctx context.Context, permission *Permission) error

	// Update 更新权限
	Update(ctx context.Context, permission *Permission) error

	// Delete 删除权限
	Delete(ctx context.Context, id string) error

	// FindByID 根据ID查找权限
	FindByID(ctx context.Context, id string) (*Permission, error)

	// FindByCode 根据Code查找权限
	FindByCode(ctx context.Context, code string) (*Permission, error)

	// List 获取权限列表
	List(ctx context.Context, offset, limit int) ([]*Permission, error)

	// ExistsByCode 检查权限Code是否存在
	ExistsByCode(ctx context.Context, code string) (bool, error)

	// GetRolePermissions 获取角色的所有权限
	GetRolePermissions(ctx context.Context, roleID string) ([]*Permission, error)

	// AssignPermissionsToRole 分配权限给角色
	AssignPermissionsToRole(ctx context.Context, roleID string, permissionIDs []string) error

	// RemovePermissionsFromRole 从角色移除权限
	RemovePermissionsFromRole(ctx context.Context, roleID string, permissionIDs []string) error

	// GetUserPermissions 获取用户的所有权限（通过角色）
	GetUserPermissions(ctx context.Context, userID string) ([]*Permission, error)
}

// MenuRepository 菜单仓储接口
type MenuRepository interface {
	// Create 创建菜单
	Create(ctx context.Context, menu *Menu) error

	// Update 更新菜单
	Update(ctx context.Context, menu *Menu) error

	// Delete 删除菜单
	Delete(ctx context.Context, id string) error

	// FindByID 根据ID查找菜单
	FindByID(ctx context.Context, id string) (*Menu, error)

	// List 获取菜单列表（可过滤父菜单）
	List(ctx context.Context, parentID *string, offset, limit int) ([]*Menu, error)

	// GetAll 获取所有菜单（用于构建树）
	GetAll(ctx context.Context) ([]*Menu, error)

	// GetChildren 获取指定菜单的子菜单
	GetChildren(ctx context.Context, parentID string) ([]*Menu, error)

	// GetRoleMenus 获取角色的所有菜单
	GetRoleMenus(ctx context.Context, roleID string) ([]*Menu, error)

	// GetUserMenus 获取用户的所有菜单（通过角色）
	GetUserMenus(ctx context.Context, userID string) ([]*Menu, error)

	// AssignMenusToRole 分配菜单给角色
	AssignMenusToRole(ctx context.Context, roleID string, menuIDs []string) error

	// RemoveMenusFromRole 从角色移除菜单
	RemoveMenusFromRole(ctx context.Context, roleID string, menuIDs []string) error

	// SyncRoleMenus 同步角色的菜单（删除旧的，添加新的）
	SyncRoleMenus(ctx context.Context, roleID string, menuIDs []string) error
}
