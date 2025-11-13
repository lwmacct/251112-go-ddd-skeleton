package mapper

import (
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/rbac"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/persistence/model"
)

// RoleToDomain 将GORM模型转换为领域实体
func RoleToDomain(m *model.Role) *rbac.Role {
	if m == nil {
		return nil
	}

	return &rbac.Role{
		ID:          m.ID,
		Name:        m.Name,
		Code:        m.Code,
		Description: m.Description,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// RoleToModel 将领域实体转换为GORM模型
func RoleToModel(d *rbac.Role) *model.Role {
	if d == nil {
		return nil
	}

	return &model.Role{
		ID:          d.ID,
		Name:        d.Name,
		Code:        d.Code,
		Description: d.Description,
		IsActive:    d.IsActive,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

// RolesToDomain 批量转换
func RolesToDomain(models []*model.Role) []*rbac.Role {
	domains := make([]*rbac.Role, 0, len(models))
	for _, m := range models {
		domains = append(domains, RoleToDomain(m))
	}
	return domains
}

// PermissionToDomain 将GORM模型转换为领域实体
func PermissionToDomain(m *model.Permission) *rbac.Permission {
	if m == nil {
		return nil
	}

	return &rbac.Permission{
		ID:          m.ID,
		Name:        m.Name,
		Code:        m.Code,
		Resource:    m.Resource,
		Action:      m.Action,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// PermissionToModel 将领域实体转换为GORM模型
func PermissionToModel(d *rbac.Permission) *model.Permission {
	if d == nil {
		return nil
	}

	return &model.Permission{
		ID:          d.ID,
		Name:        d.Name,
		Code:        d.Code,
		Resource:    d.Resource,
		Action:      d.Action,
		Description: d.Description,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

// PermissionsToDomain 批量转换
func PermissionsToDomain(models []*model.Permission) []*rbac.Permission {
	domains := make([]*rbac.Permission, 0, len(models))
	for _, m := range models {
		domains = append(domains, PermissionToDomain(m))
	}
	return domains
}

// MenuToDomain 将GORM模型转换为领域实体
func MenuToDomain(m *model.Menu) *rbac.Menu {
	if m == nil {
		return nil
	}

	return &rbac.Menu{
		ID:          m.ID,
		Name:        m.Name,
		Path:        m.Path,
		Icon:        m.Icon,
		ParentID:    m.ParentID,
		SortOrder:   m.SortOrder,
		Type:        rbac.MenuType(m.Type),
		IsVisible:   m.IsVisible,
		Component:   m.Component,
		Permission:  m.Permission,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// MenuToModel 将领域实体转换为GORM模型
func MenuToModel(d *rbac.Menu) *model.Menu {
	if d == nil {
		return nil
	}

	return &model.Menu{
		ID:          d.ID,
		Name:        d.Name,
		Path:        d.Path,
		Icon:        d.Icon,
		ParentID:    d.ParentID,
		SortOrder:   d.SortOrder,
		Type:        string(d.Type),
		IsVisible:   d.IsVisible,
		Component:   d.Component,
		Permission:  d.Permission,
		Description: d.Description,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

// MenusToDomain 批量转换
func MenusToDomain(models []*model.Menu) []*rbac.Menu {
	domains := make([]*rbac.Menu, 0, len(models))
	for _, m := range models {
		domains = append(domains, MenuToDomain(m))
	}
	return domains
}
