package rbac

import (
	"time"
)

// MenuType 菜单类型
type MenuType string

const (
	MenuTypeDir  MenuType = "dir"  // 目录（父菜单）
	MenuTypeMenu MenuType = "menu" // 菜单（可点击跳转）
	MenuTypeLink MenuType = "link" // 外部链接
)

// Menu 表示系统中的菜单实体（聚合根）
type Menu struct {
	ID          string
	Name        string   // 菜单名称，如 "用户管理"
	Path        string   // 路由路径，如 "/users"
	Icon        string   // 图标名称
	ParentID    *string  // 父菜单ID（支持两层结构）
	SortOrder   int      // 排序顺序（越小越靠前）
	Type        MenuType // 菜单类型
	IsVisible   bool     // 是否可见
	Component   string   // 前端组件路径（可选）
	Permission  string   // 关联的权限码（可选）
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewMenu 创建新的菜单
func NewMenu(name, path, icon string, menuType MenuType, parentID *string) (*Menu, error) {
	if name == "" {
		return nil, ErrInvalidMenuName
	}
	if path == "" {
		return nil, ErrInvalidMenuPath
	}

	now := time.Now()
	return &Menu{
		Name:      name,
		Path:      path,
		Icon:      icon,
		Type:      menuType,
		ParentID:  parentID,
		SortOrder: 0,
		IsVisible: true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// UpdateInfo 更新菜单信息
func (m *Menu) UpdateInfo(name, path, icon, component, permission, description string) error {
	if name == "" {
		return ErrInvalidMenuName
	}
	if path == "" {
		return ErrInvalidMenuPath
	}

	m.Name = name
	m.Path = path
	m.Icon = icon
	m.Component = component
	m.Permission = permission
	m.Description = description
	m.UpdatedAt = time.Now()
	return nil
}

// SetParent 设置父菜单
func (m *Menu) SetParent(parentID *string) {
	m.ParentID = parentID
	m.UpdatedAt = time.Now()
}

// UpdateOrder 更新排序
func (m *Menu) UpdateOrder(sortOrder int) {
	m.SortOrder = sortOrder
	m.UpdatedAt = time.Now()
}

// Show 显示菜单
func (m *Menu) Show() {
	m.IsVisible = true
	m.UpdatedAt = time.Now()
}

// Hide 隐藏菜单
func (m *Menu) Hide() {
	m.IsVisible = false
	m.UpdatedAt = time.Now()
}

// IsParent 判断是否为父菜单（目录）
func (m *Menu) IsParent() bool {
	return m.ParentID == nil && m.Type == MenuTypeDir
}

// IsChild 判断是否为子菜单
func (m *Menu) IsChild() bool {
	return m.ParentID != nil
}

// HasPermission 判断是否有关联权限
func (m *Menu) HasPermission() bool {
	return m.Permission != ""
}
