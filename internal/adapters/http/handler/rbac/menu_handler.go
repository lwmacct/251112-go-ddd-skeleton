package rbac

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/adapters/http/response"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/application/menu"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/rbac"
)

// MenuHandler 菜单处理器
type MenuHandler struct {
	menuService *menu.Service
}

// NewMenuHandler 创建菜单处理器
func NewMenuHandler(menuService *menu.Service) *MenuHandler {
	return &MenuHandler{
		menuService: menuService,
	}
}

// ========== 用户端点（需要认证）==========

// GetUserMenuTree 获取当前用户的菜单树
// GET /api/menus/user/tree
func (h *MenuHandler) GetUserMenuTree(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	menuTree, err := h.menuService.GetUserMenuTree(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, menuTree)
}

// ========== 管理端点（需要管理员权限）==========

// CreateMenu 创建菜单
// POST /api/admin/menus
func (h *MenuHandler) CreateMenu(c *gin.Context) {
	var req menu.CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.menuService.CreateMenu(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// UpdateMenu 更新菜单
// PUT /api/admin/menus/:id
func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, rbac.ErrMenuNotFound)
		return
	}

	var req menu.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.menuService.UpdateMenu(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// DeleteMenu 删除菜单
// DELETE /api/admin/menus/:id
func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, rbac.ErrMenuNotFound)
		return
	}

	if err := h.menuService.DeleteMenu(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "菜单删除成功"})
}

// GetAllMenuTree 获取所有菜单树（管理员查看所有菜单）
// GET /api/admin/menus/tree
func (h *MenuHandler) GetAllMenuTree(c *gin.Context) {
	menuTree, err := h.menuService.GetAllMenuTree(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, menuTree)
}

// UpdateMenuOrder 批量更新菜单排序
// PUT /api/admin/menus/order
func (h *MenuHandler) UpdateMenuOrder(c *gin.Context) {
	var req menu.UpdateMenuOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	if err := h.menuService.UpdateMenuOrder(c.Request.Context(), req.Orders); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "菜单排序更新成功"})
}

// AssignMenusToRole 给角色分配菜单
// POST /api/admin/roles/:roleId/menus
func (h *MenuHandler) AssignMenusToRole(c *gin.Context) {
	roleID := c.Param("roleId")
	if roleID == "" {
		response.Error(c, rbac.ErrRoleNotFound)
		return
	}

	var req menu.AssignMenusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	if err := h.menuService.AssignMenusToRole(c.Request.Context(), roleID, req.MenuIDs); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "角色菜单分配成功"})
}

// GetRoleMenuTree 获取角色的菜单树
// GET /api/admin/roles/:roleId/menus
func (h *MenuHandler) GetRoleMenuTree(c *gin.Context) {
	roleID := c.Param("roleId")
	if roleID == "" {
		response.Error(c, rbac.ErrRoleNotFound)
		return
	}

	menuTree, err := h.menuService.GetRoleMenuTree(c.Request.Context(), roleID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, menuTree)
}

// getCurrentUserID 从 context 获取当前用户 ID
func getCurrentUserID(c *gin.Context) (string, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return "", errors.New("unauthorized")
	}

	id, ok := userID.(string)
	if !ok {
		return "", errors.New("unauthorized")
	}

	return id, nil
}
