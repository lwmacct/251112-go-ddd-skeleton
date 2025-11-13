package role

import "time"

// ========== 请求 DTO ==========

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsActive    *bool  `json:"isActive"`
}

// AssignUserRoleRequest 分配角色给用户请求
type AssignUserRoleRequest struct {
	UserID string `json:"userId" binding:"required"`
}

// AssignPermissionsRequest 分配权限给角色请求
type AssignPermissionsRequest struct {
	PermissionIDs []string `json:"permissionIds" binding:"required"`
}

// ========== 响应 DTO ==========

// RoleResponse 角色响应
type RoleResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// PermissionResponse 权限响应
type PermissionResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
