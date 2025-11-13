package menu

import (
	"context"

	"github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/rbac"
	"github.com/oklog/ulid/v2"
)

// CreateMenu 创建菜单
func (s *Service) CreateMenu(ctx context.Context, req CreateMenuRequest) (*MenuResponse, error) {
	// 验证父菜单层级（如果有父菜单）
	if req.ParentID != nil {
		if err := s.domainService.ValidateMenuHierarchy(ctx, req.ParentID); err != nil {
			return nil, err
		}
	}

	// 创建领域实体
	menu, err := rbac.NewMenu(req.Name, req.Path, req.Icon, rbac.MenuType(req.Type), req.ParentID)
	if err != nil {
		return nil, err
	}

	// 设置ID
	menu.ID = ulid.Make().String()

	// 设置其他属性
	menu.SortOrder = req.SortOrder
	menu.Component = req.Component
	menu.Permission = req.Permission
	menu.Description = req.Description

	// 保存到数据库
	if err := s.menuRepo.Create(ctx, menu); err != nil {
		return nil, err
	}

	return toMenuResponse(menu), nil
}

// UpdateMenu 更新菜单
func (s *Service) UpdateMenu(ctx context.Context, id string, req UpdateMenuRequest) (*MenuResponse, error) {
	// 查找菜单
	menu, err := s.menuRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 如果更新了父菜单，验证层级
	if req.ParentID != nil && (menu.ParentID == nil || *menu.ParentID != *req.ParentID) {
		if err := s.domainService.ValidateMenuHierarchy(ctx, req.ParentID); err != nil {
			return nil, err
		}
		menu.SetParent(req.ParentID)
	}

	// 更新菜单信息
	if err := menu.UpdateInfo(req.Name, req.Path, req.Icon, req.Component, req.Permission, req.Description); err != nil {
		return nil, err
	}

	// 更新排序
	if req.SortOrder != nil {
		menu.UpdateOrder(*req.SortOrder)
	}

	// 更新可见性
	if req.IsVisible != nil {
		if *req.IsVisible {
			menu.Show()
		} else {
			menu.Hide()
		}
	}

	// 保存到数据库
	if err := s.menuRepo.Update(ctx, menu); err != nil {
		return nil, err
	}

	return toMenuResponse(menu), nil
}

// DeleteMenu 删除菜单
func (s *Service) DeleteMenu(ctx context.Context, id string) error {
	// 检查是否有子菜单
	children, err := s.menuRepo.GetChildren(ctx, id)
	if err != nil {
		return err
	}

	if len(children) > 0 {
		return rbac.ErrCannotDeleteMenuWithChildren
	}

	// 删除菜单
	return s.menuRepo.Delete(ctx, id)
}

// AssignMenusToRole 给角色分配菜单
func (s *Service) AssignMenusToRole(ctx context.Context, roleID string, menuIDs []string) error {
	// 验证所有菜单是否存在
	for _, menuID := range menuIDs {
		_, err := s.menuRepo.FindByID(ctx, menuID)
		if err != nil {
			return err
		}
	}

	// 同步角色菜单（会删除旧的，添加新的）
	return s.menuRepo.SyncRoleMenus(ctx, roleID, menuIDs)
}

// UpdateMenuOrder 批量更新菜单排序
func (s *Service) UpdateMenuOrder(ctx context.Context, orders []MenuOrderRequest) error {
	for _, order := range orders {
		menu, err := s.menuRepo.FindByID(ctx, order.MenuID)
		if err != nil {
			return err
		}

		menu.UpdateOrder(order.SortOrder)

		if err := s.menuRepo.Update(ctx, menu); err != nil {
			return err
		}
	}

	return nil
}

// toMenuResponse 转换为响应DTO
func toMenuResponse(m *rbac.Menu) *MenuResponse {
	return &MenuResponse{
		ID:          m.ID,
		Name:        m.Name,
		Path:        m.Path,
		Icon:        m.Icon,
		ParentID:    m.ParentID,
		SortOrder:   m.SortOrder,
		Type:        string(m.Type),
		IsVisible:   m.IsVisible,
		Component:   m.Component,
		Permission:  m.Permission,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
