package menu

import "time"

// ========== 请求 DTO ==========

// CreateMenuRequest 创建菜单请求
type CreateMenuRequest struct {
	Name        string  `json:"name" binding:"required"`
	Path        string  `json:"path" binding:"required"`
	Icon        string  `json:"icon"`
	ParentID    *string `json:"parentId"`
	SortOrder   int     `json:"sortOrder"`
	Type        string  `json:"type" binding:"required,oneof=dir menu link"` // dir, menu, link
	Component   string  `json:"component"`
	Permission  string  `json:"permission"`
	Description string  `json:"description"`
}

// UpdateMenuRequest 更新菜单请求
type UpdateMenuRequest struct {
	Name        string  `json:"name" binding:"required"`
	Path        string  `json:"path" binding:"required"`
	Icon        string  `json:"icon"`
	ParentID    *string `json:"parentId"`
	SortOrder   *int    `json:"sortOrder"`
	Component   string  `json:"component"`
	Permission  string  `json:"permission"`
	Description string  `json:"description"`
	IsVisible   *bool   `json:"isVisible"`
}

// AssignMenusRequest 分配菜单给角色请求
type AssignMenusRequest struct {
	MenuIDs []string `json:"menuIds" binding:"required"`
}

// MenuOrderRequest 菜单排序请求
type MenuOrderRequest struct {
	MenuID    string `json:"menuId" binding:"required"`
	SortOrder int    `json:"sortOrder"`
}

// UpdateMenuOrderRequest 批量更新菜单排序请求
type UpdateMenuOrderRequest struct {
	Orders []MenuOrderRequest `json:"orders" binding:"required"`
}

// ========== 响应 DTO ==========

// MenuResponse 菜单响应
type MenuResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Icon        string    `json:"icon,omitempty"`
	ParentID    *string   `json:"parentId,omitempty"`
	SortOrder   int       `json:"sortOrder"`
	Type        string    `json:"type"`
	IsVisible   bool      `json:"isVisible"`
	Component   string    `json:"component,omitempty"`
	Permission  string    `json:"permission,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// MenuTreeResponse 菜单树响应
type MenuTreeResponse struct {
	Menus []*MenuTreeItem `json:"menus"`
}

// MenuTreeItem 菜单树项（支持递归）
type MenuTreeItem struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Path        string          `json:"path"`
	Icon        string          `json:"icon,omitempty"`
	Type        string          `json:"type"` // dir, menu, link
	SortOrder   int             `json:"sortOrder"`
	IsVisible   bool            `json:"isVisible"`
	Component   string          `json:"component,omitempty"`
	Permission  string          `json:"permission,omitempty"`
	Description string          `json:"description,omitempty"`
	Children    []*MenuTreeItem `json:"children,omitempty"` // 子菜单（递归）
}
