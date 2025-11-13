package model

import "time"

// RolePermission 角色-权限关联表
type RolePermission struct {
	RoleID       string    `gorm:"primaryKey;type:varchar(26)"`
	PermissionID string    `gorm:"primaryKey;type:varchar(26)"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`

	Role       Role       `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
	Permission Permission `gorm:"foreignKey:PermissionID;constraint:OnDelete:CASCADE"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}
