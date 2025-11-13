package rbac

import (
	"time"
)

// Role 表示系统中的角色实体（聚合根）
type Role struct {
	ID          string
	Name        string
	Code        string // 唯一标识码，如 "admin", "user", "editor"
	Description string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewRole 创建新的角色
func NewRole(name, code, description string) (*Role, error) {
	if name == "" {
		return nil, ErrInvalidRoleName
	}
	if code == "" {
		return nil, ErrInvalidRoleCode
	}

	now := time.Now()
	return &Role{
		Name:        name,
		Code:        code,
		Description: description,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Activate 激活角色
func (r *Role) Activate() {
	r.IsActive = true
	r.UpdatedAt = time.Now()
}

// Deactivate 停用角色
func (r *Role) Deactivate() {
	r.IsActive = false
	r.UpdatedAt = time.Now()
}

// UpdateInfo 更新角色信息
func (r *Role) UpdateInfo(name, description string) error {
	if name == "" {
		return ErrInvalidRoleName
	}

	r.Name = name
	r.Description = description
	r.UpdatedAt = time.Now()
	return nil
}

// IsSystemRole 判断是否为系统内置角色（不可删除）
func (r *Role) IsSystemRole() bool {
	// 系统内置角色：admin, user
	return r.Code == "admin" || r.Code == "user"
}
