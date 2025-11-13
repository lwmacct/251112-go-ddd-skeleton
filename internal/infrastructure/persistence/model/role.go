package model

import "time"

// Role GORM角色模型
type Role struct {
	ID          string    `gorm:"primaryKey;type:varchar(26)"`
	Name        string    `gorm:"not null;type:varchar(100)"`
	Code        string    `gorm:"uniqueIndex;not null;type:varchar(50)"`
	Description string    `gorm:"type:varchar(500)"`
	IsActive    bool      `gorm:"default:true"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	// 关联关系（用于预加载）
	Users       []User       `gorm:"many2many:user_roles;"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
	Menus       []Menu       `gorm:"many2many:role_menus;"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}
