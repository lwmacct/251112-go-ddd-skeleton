package menu

import (
	"context"
	"sort"

	"github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/rbac"
)

// Service 菜单应用服务
type Service struct {
	domainService *rbac.Service
	menuRepo      rbac.MenuRepository
}

// NewService 创建菜单应用服务
func NewService(domainService *rbac.Service, menuRepo rbac.MenuRepository) *Service {
	return &Service{
		domainService: domainService,
		menuRepo:      menuRepo,
	}
}

// GetUserMenuTree 获取用户的菜单树（核心API）
func (s *Service) GetUserMenuTree(ctx context.Context, userID string) (*MenuTreeResponse, error) {
	// 使用领域服务获取菜单树
	menuTree, err := s.domainService.GetUserMenuTree(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 转换为DTO并排序
	return s.buildMenuTreeResponse(menuTree), nil
}

// GetAllMenuTree 获取所有菜单树（用于管理后台）
func (s *Service) GetAllMenuTree(ctx context.Context) (*MenuTreeResponse, error) {
	// 获取所有菜单
	allMenus, err := s.menuRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// 手动构建树形结构
	nodeMap := make(map[string]*rbac.MenuNode)
	var rootNodes []*rbac.MenuNode

	// 第一遍：创建所有节点
	for _, menu := range allMenus {
		node := &rbac.MenuNode{
			Menu:     menu,
			Children: make([]*rbac.MenuNode, 0),
		}
		nodeMap[menu.ID] = node

		if menu.ParentID == nil {
			rootNodes = append(rootNodes, node)
		}
	}

	// 第二遍：建立父子关系
	for _, menu := range allMenus {
		if menu.ParentID != nil {
			if parentNode, ok := nodeMap[*menu.ParentID]; ok {
				if childNode, ok := nodeMap[menu.ID]; ok {
					parentNode.Children = append(parentNode.Children, childNode)
				}
			}
		}
	}

	return s.buildMenuTreeResponse(rootNodes), nil
}

// GetRoleMenuTree 获取角色的菜单树
func (s *Service) GetRoleMenuTree(ctx context.Context, roleID string) (*MenuTreeResponse, error) {
	menuTree, err := s.domainService.GetRoleMenuTree(ctx, roleID)
	if err != nil {
		return nil, err
	}

	return s.buildMenuTreeResponse(menuTree), nil
}

// buildMenuTreeResponse 构建菜单树响应
func (s *Service) buildMenuTreeResponse(nodes []*rbac.MenuNode) *MenuTreeResponse {
	items := make([]*MenuTreeItem, 0, len(nodes))
	for _, node := range nodes {
		items = append(items, s.nodeToTreeItem(node))
	}

	// 按SortOrder排序
	sort.Slice(items, func(i, j int) bool {
		if items[i].SortOrder == items[j].SortOrder {
			return items[i].ID < items[j].ID
		}
		return items[i].SortOrder < items[j].SortOrder
	})

	return &MenuTreeResponse{
		Menus: items,
	}
}

// nodeToTreeItem 将领域节点转换为DTO
func (s *Service) nodeToTreeItem(node *rbac.MenuNode) *MenuTreeItem {
	if node == nil {
		return nil
	}

	item := &MenuTreeItem{
		ID:          node.Menu.ID,
		Name:        node.Menu.Name,
		Path:        node.Menu.Path,
		Icon:        node.Menu.Icon,
		Type:        string(node.Menu.Type),
		SortOrder:   node.Menu.SortOrder,
		IsVisible:   node.Menu.IsVisible,
		Component:   node.Menu.Component,
		Permission:  node.Menu.Permission,
		Description: node.Menu.Description,
	}

	// 递归处理子节点
	if len(node.Children) > 0 {
		item.Children = make([]*MenuTreeItem, 0, len(node.Children))
		for _, child := range node.Children {
			item.Children = append(item.Children, s.nodeToTreeItem(child))
		}

		// 子节点排序
		sort.Slice(item.Children, func(i, j int) bool {
			if item.Children[i].SortOrder == item.Children[j].SortOrder {
				return item.Children[i].ID < item.Children[j].ID
			}
			return item.Children[i].SortOrder < item.Children[j].SortOrder
		})
	}

	return item
}
