package rbac

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/adapters/http/response"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/application/role"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/rbac"
)

// RoleHandler 角色处理器
type RoleHandler struct {
	roleService *role.Service
}

// NewRoleHandler 创建角色处理器
func NewRoleHandler(roleService *role.Service) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// ========== 管理端点（需要管理员权限）==========

// CreateRole 创建角色
// POST /api/admin/roles
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req role.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.roleService.CreateRole(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// UpdateRole 更新角色
// PUT /api/admin/roles/:id
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, rbac.ErrRoleNotFound)
		return
	}

	var req role.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.roleService.UpdateRole(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// DeleteRole 删除角色
// DELETE /api/admin/roles/:id
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, rbac.ErrRoleNotFound)
		return
	}

	if err := h.roleService.DeleteRole(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "角色删除成功"})
}

// GetRole 获取角色详情
// GET /api/admin/roles/:id
func (h *RoleHandler) GetRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, rbac.ErrRoleNotFound)
		return
	}

	result, err := h.roleService.GetRole(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// ListRoles 列出角色
// GET /api/admin/roles
func (h *RoleHandler) ListRoles(c *gin.Context) {
	// 分页参数
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if limit > 100 {
		limit = 100
	}

	roles, err := h.roleService.ListRoles(c.Request.Context(), offset, limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"roles":  roles,
		"offset": offset,
		"limit":  limit,
	})
}

// AssignRoleToUser 给用户分配角色
// POST /api/admin/users/:userId/roles/:roleId
func (h *RoleHandler) AssignRoleToUser(c *gin.Context) {
	userID := c.Param("userId")
	roleID := c.Param("roleId")

	if userID == "" || roleID == "" {
		response.Error(c, rbac.ErrRoleNotFound)
		return
	}

	if err := h.roleService.AssignRoleToUser(c.Request.Context(), userID, roleID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "角色分配成功"})
}

// RemoveRoleFromUser 从用户移除角色
// DELETE /api/admin/users/:userId/roles/:roleId
func (h *RoleHandler) RemoveRoleFromUser(c *gin.Context) {
	userID := c.Param("userId")
	roleID := c.Param("roleId")

	if userID == "" || roleID == "" {
		response.Error(c, rbac.ErrRoleNotFound)
		return
	}

	if err := h.roleService.RemoveRoleFromUser(c.Request.Context(), userID, roleID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "角色移除成功"})
}

// GetUserRoles 获取用户的所有角色
// GET /api/admin/users/:userId/roles
func (h *RoleHandler) GetUserRoles(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		response.Error(c, rbac.ErrUserRoleNotFound)
		return
	}

	roles, err := h.roleService.GetUserRoles(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"roles": roles})
}

// AssignPermissionsToRole 给角色分配权限
// POST /api/admin/roles/:roleId/permissions
func (h *RoleHandler) AssignPermissionsToRole(c *gin.Context) {
	roleID := c.Param("roleId")
	if roleID == "" {
		response.Error(c, rbac.ErrRoleNotFound)
		return
	}

	var req role.AssignPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	if err := h.roleService.AssignPermissionsToRole(c.Request.Context(), roleID, req.PermissionIDs); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "权限分配成功"})
}

// GetRolePermissions 获取角色的所有权限
// GET /api/admin/roles/:roleId/permissions
func (h *RoleHandler) GetRolePermissions(c *gin.Context) {
	roleID := c.Param("roleId")
	if roleID == "" {
		response.Error(c, rbac.ErrRoleNotFound)
		return
	}

	permissions, err := h.roleService.GetRolePermissions(c.Request.Context(), roleID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"permissions": permissions})
}
