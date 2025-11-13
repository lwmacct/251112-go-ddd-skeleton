package model

import "time"

// Permission GORM权限模型
type Permission struct {
	ID          string    `gorm:"primaryKey;type:varchar(26)"`
	Name        string    `gorm:"not null;type:varchar(100)"`
	Code        string    `gorm:"uniqueIndex;not null;type:varchar(100)"` // 如 user:create
	Resource    string    `gorm:"not null;index;type:varchar(50)"`         // 如 user
	Action      string    `gorm:"not null;index;type:varchar(50)"`         // 如 create
	Description string    `gorm:"type:varchar(500)"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	// 关联关系
	Roles []Role `gorm:"many2many:role_permissions;"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}
