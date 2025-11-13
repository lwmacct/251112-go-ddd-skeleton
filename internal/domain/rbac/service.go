package rbac

import (
	"context"
)

// Service RBAC领域服务
type Service struct {
	roleRepo       RoleRepository
	permissionRepo PermissionRepository
	menuRepo       MenuRepository
}

// NewService 创建RBAC领域服务
func NewService(
	roleRepo RoleRepository,
	permissionRepo PermissionRepository,
	menuRepo MenuRepository,
) *Service {
	return &Service{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		menuRepo:       menuRepo,
	}
}

// CheckPermission 检查用户是否具有指定权限
func (s *Service) CheckPermission(ctx context.Context, userID, permissionCode string) (bool, error) {
	// 获取用户的所有权限
	permissions, err := s.permissionRepo.GetUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}

	// 检查是否包含指定权限
	for _, perm := range permissions {
		if perm.MatchesCode(permissionCode) {
			return true, nil
		}
	}

	return false, nil
}

// CheckUserHasRole 检查用户是否具有指定角色
func (s *Service) CheckUserHasRole(ctx context.Context, userID, roleCode string) (bool, error) {
	// 获取用户的所有角色
	roles, err := s.roleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return false, err
	}

	// 检查是否包含指定角色
	for _, role := range roles {
		if role.Code == roleCode {
			return true, nil
		}
	}

	return false, nil
}

// GetUserMenuTree 获取用户的菜单树（两层结构）
func (s *Service) GetUserMenuTree(ctx context.Context, userID string) ([]*MenuNode, error) {
	// 获取用户的所有菜单（已按角色过滤）
	menus, err := s.menuRepo.GetUserMenus(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 构建两层树形结构
	return s.buildMenuTree(menus), nil
}

// GetRoleMenuTree 获取角色的菜单树
func (s *Service) GetRoleMenuTree(ctx context.Context, roleID string) ([]*MenuNode, error) {
	// 获取角色的所有菜单
	menus, err := s.menuRepo.GetRoleMenus(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// 构建两层树形结构
	return s.buildMenuTree(menus), nil
}

// ValidateRoleAssignment 验证角色分配是否合法
func (s *Service) ValidateRoleAssignment(ctx context.Context, roleID string) error {
	// 检查角色是否存在
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return err
	}

	// 检查角色是否激活
	if !role.IsActive {
		return ErrRoleInactive
	}

	return nil
}

// ValidateMenuHierarchy 验证菜单层级是否合法（只允许两层）
func (s *Service) ValidateMenuHierarchy(ctx context.Context, parentID *string) error {
	// 如果没有父菜单，说明是一级菜单，合法
	if parentID == nil {
		return nil
	}

	// 检查父菜单是否存在
	parentMenu, err := s.menuRepo.FindByID(ctx, *parentID)
	if err != nil {
		return err
	}

	// 检查父菜单是否已经是子菜单（不允许三层）
	if parentMenu.IsChild() {
		return ErrMenuHierarchyTooDeep
	}

	return nil
}

// MenuNode 菜单树节点（用于构建菜单树）
type MenuNode struct {
	*Menu
	Children []*MenuNode
}

// buildMenuTree 构建两层菜单树
func (s *Service) buildMenuTree(menus []*Menu) []*MenuNode {
	// 创建节点映射
	nodeMap := make(map[string]*MenuNode)
	var rootNodes []*MenuNode

	// 第一遍：创建所有节点
	for _, menu := range menus {
		// 只显示可见的菜单
		if !menu.IsVisible {
			continue
		}

		node := &MenuNode{
			Menu:     menu,
			Children: make([]*MenuNode, 0),
		}
		nodeMap[menu.ID] = node

		// 如果是根节点（没有父菜单），加入根节点列表
		if menu.ParentID == nil {
			rootNodes = append(rootNodes, node)
		}
	}

	// 第二遍：建立父子关系
	for _, menu := range menus {
		if menu.ParentID != nil {
			parentNode, parentExists := nodeMap[*menu.ParentID]
			childNode, childExists := nodeMap[menu.ID]

			if parentExists && childExists {
				parentNode.Children = append(parentNode.Children, childNode)
			}
		}
	}

	// 按 SortOrder 排序（可在应用层实现）
	return rootNodes
}
