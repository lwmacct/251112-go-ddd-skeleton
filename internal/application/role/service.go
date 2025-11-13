package role

import (
	"context"

	"github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/rbac"
	"github.com/oklog/ulid/v2"
)

// Service 角色应用服务
type Service struct {
	roleRepo       rbac.RoleRepository
	permissionRepo rbac.PermissionRepository
	domainService  *rbac.Service
}

// NewService 创建角色应用服务
func NewService(
	roleRepo rbac.RoleRepository,
	permissionRepo rbac.PermissionRepository,
	domainService *rbac.Service,
) *Service {
	return &Service{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		domainService:  domainService,
	}
}

// CreateRole 创建角色
func (s *Service) CreateRole(ctx context.Context, req CreateRoleRequest) (*RoleResponse, error) {
	// 检查角色代码是否已存在
	exists, err := s.roleRepo.ExistsByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, rbac.ErrRoleCodeDuplicate
	}

	// 创建领域实体
	role, err := rbac.NewRole(req.Name, req.Code, req.Description)
	if err != nil {
		return nil, err
	}

	// 设置ID
	role.ID = ulid.Make().String()

	// 保存到数据库
	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, err
	}

	return toRoleResponse(role), nil
}

// UpdateRole 更新角色
func (s *Service) UpdateRole(ctx context.Context, id string, req UpdateRoleRequest) (*RoleResponse, error) {
	// 查找角色
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新角色信息
	if err := role.UpdateInfo(req.Name, req.Description); err != nil {
		return nil, err
	}

	// 更新状态
	if req.IsActive != nil {
		if *req.IsActive {
			role.Activate()
		} else {
			role.Deactivate()
		}
	}

	// 保存到数据库
	if err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, err
	}

	return toRoleResponse(role), nil
}

// DeleteRole 删除角色
func (s *Service) DeleteRole(ctx context.Context, id string) error {
	// 查找角色
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 检查是否为系统角色
	if role.IsSystemRole() {
		return rbac.ErrCannotDeleteSystemRole
	}

	// 删除角色
	return s.roleRepo.Delete(ctx, id)
}

// GetRole 获取角色详情
func (s *Service) GetRole(ctx context.Context, id string) (*RoleResponse, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return toRoleResponse(role), nil
}

// ListRoles 列出角色
func (s *Service) ListRoles(ctx context.Context, offset, limit int) ([]*RoleResponse, error) {
	roles, err := s.roleRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]*RoleResponse, 0, len(roles))
	for _, role := range roles {
		responses = append(responses, toRoleResponse(role))
	}

	return responses, nil
}

// AssignRoleToUser 给用户分配角色
func (s *Service) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	// 验证角色分配是否合法
	if err := s.domainService.ValidateRoleAssignment(ctx, roleID); err != nil {
		return err
	}

	return s.roleRepo.AssignRoleToUser(ctx, userID, roleID)
}

// RemoveRoleFromUser 从用户移除角色
func (s *Service) RemoveRoleFromUser(ctx context.Context, userID, roleID string) error {
	return s.roleRepo.RemoveRoleFromUser(ctx, userID, roleID)
}

// GetUserRoles 获取用户的所有角色
func (s *Service) GetUserRoles(ctx context.Context, userID string) ([]*RoleResponse, error) {
	roles, err := s.roleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*RoleResponse, 0, len(roles))
	for _, role := range roles {
		responses = append(responses, toRoleResponse(role))
	}

	return responses, nil
}

// AssignPermissionsToRole 给角色分配权限
func (s *Service) AssignPermissionsToRole(ctx context.Context, roleID string, permissionIDs []string) error {
	// 验证所有权限是否存在
	for _, permID := range permissionIDs {
		_, err := s.permissionRepo.FindByID(ctx, permID)
		if err != nil {
			return err
		}
	}

	// 先移除旧的权限，再添加新的（可以在repository层优化为事务）
	return s.permissionRepo.AssignPermissionsToRole(ctx, roleID, permissionIDs)
}

// GetRolePermissions 获取角色的所有权限
func (s *Service) GetRolePermissions(ctx context.Context, roleID string) ([]*PermissionResponse, error) {
	permissions, err := s.permissionRepo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return nil, err
	}

	responses := make([]*PermissionResponse, 0, len(permissions))
	for _, perm := range permissions {
		responses = append(responses, toPermissionResponse(perm))
	}

	return responses, nil
}

// toRoleResponse 转换为响应DTO
func toRoleResponse(r *rbac.Role) *RoleResponse {
	return &RoleResponse{
		ID:          r.ID,
		Name:        r.Name,
		Code:        r.Code,
		Description: r.Description,
		IsActive:    r.IsActive,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// toPermissionResponse 转换为响应DTO
func toPermissionResponse(p *rbac.Permission) *PermissionResponse {
	return &PermissionResponse{
		ID:          p.ID,
		Name:        p.Name,
		Code:        p.Code,
		Resource:    p.Resource,
		Action:      p.Action,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
