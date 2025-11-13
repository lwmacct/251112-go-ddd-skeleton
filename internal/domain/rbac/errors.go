package rbac

import "errors"

var (
	// Role相关错误
	ErrInvalidRoleName     = errors.New("角色名称不能为空")
	ErrInvalidRoleCode     = errors.New("角色代码不能为空")
	ErrRoleNotFound        = errors.New("角色不存在")
	ErrRoleCodeDuplicate   = errors.New("角色代码已存在")
	ErrRoleInactive        = errors.New("角色未激活")
	ErrCannotDeleteSystemRole = errors.New("不能删除系统内置角色")

	// Permission相关错误
	ErrInvalidPermissionName     = errors.New("权限名称不能为空")
	ErrInvalidPermissionResource = errors.New("权限资源不能为空")
	ErrInvalidPermissionAction   = errors.New("权限操作不能为空")
	ErrPermissionNotFound        = errors.New("权限不存在")
	ErrPermissionCodeDuplicate   = errors.New("权限代码已存在")
	ErrPermissionDenied          = errors.New("权限不足")

	// Menu相关错误
	ErrInvalidMenuName       = errors.New("菜单名称不能为空")
	ErrInvalidMenuPath       = errors.New("菜单路径不能为空")
	ErrMenuNotFound          = errors.New("菜单不存在")
	ErrMenuHierarchyTooDeep  = errors.New("菜单层级过深，只支持两层")
	ErrParentMenuNotFound    = errors.New("父菜单不存在")
	ErrCannotDeleteMenuWithChildren = errors.New("不能删除有子菜单的菜单")

	// 关联关系错误
	ErrUserRoleNotFound      = errors.New("用户角色关联不存在")
	ErrRolePermissionNotFound = errors.New("角色权限关联不存在")
	ErrRoleMenuNotFound      = errors.New("角色菜单关联不存在")
)
