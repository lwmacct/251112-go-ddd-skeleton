package model

import "time"

// RoleMenu 角色-菜单关联表
type RoleMenu struct {
	RoleID    string    `gorm:"primaryKey;type:varchar(26)"`
	MenuID    string    `gorm:"primaryKey;type:varchar(26)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	Role Role `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
	Menu Menu `gorm:"foreignKey:MenuID;constraint:OnDelete:CASCADE"`
}

// TableName 指定表名
func (RoleMenu) TableName() string {
	return "role_menus"
}
