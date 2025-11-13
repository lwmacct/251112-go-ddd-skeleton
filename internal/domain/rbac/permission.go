package rbac

import (
	"fmt"
	"time"
)

// Permission 表示系统中的权限实体
type Permission struct {
	ID          string
	Name        string
	Code        string // 权限码，如 "user:create", "order:read"
	Resource    string // 资源，如 "user", "order", "menu"
	Action      string // 操作，如 "create", "read", "update", "delete"
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewPermission 创建新的权限
func NewPermission(name, resource, action, description string) (*Permission, error) {
	if name == "" {
		return nil, ErrInvalidPermissionName
	}
	if resource == "" {
		return nil, ErrInvalidPermissionResource
	}
	if action == "" {
		return nil, ErrInvalidPermissionAction
	}

	code := fmt.Sprintf("%s:%s", resource, action)
	now := time.Now()

	return &Permission{
		Name:        name,
		Code:        code,
		Resource:    resource,
		Action:      action,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// UpdateInfo 更新权限信息
func (p *Permission) UpdateInfo(name, description string) error {
	if name == "" {
		return ErrInvalidPermissionName
	}

	p.Name = name
	p.Description = description
	p.UpdatedAt = time.Now()
	return nil
}

// Matches 检查是否匹配指定的资源和操作
func (p *Permission) Matches(resource, action string) bool {
	return p.Resource == resource && p.Action == action
}

// MatchesCode 检查是否匹配指定的权限码
func (p *Permission) MatchesCode(code string) bool {
	return p.Code == code
}
