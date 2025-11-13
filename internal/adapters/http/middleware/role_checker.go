package middleware

import (
	"context"

	"github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/rbac"
)

// rbacRoleChecker 基于 RBAC 领域服务的角色检查器实现
type rbacRoleChecker struct {
	rbacService *rbac.Service
}

// NewRBACRoleChecker 创建基于 RBAC 的角色检查器
func NewRBACRoleChecker(rbacService *rbac.Service) RoleChecker {
	return &rbacRoleChecker{
		rbacService: rbacService,
	}
}

// IsAdmin 检查用户是否是管理员
// 通过检查用户是否具有 "admin" 角色来判断
func (r *rbacRoleChecker) IsAdmin(userID string) (bool, error) {
	// 使用 RBAC 领域服务检查用户是否具有 admin 角色
	return r.rbacService.CheckUserHasRole(context.Background(), userID, "admin")
}
