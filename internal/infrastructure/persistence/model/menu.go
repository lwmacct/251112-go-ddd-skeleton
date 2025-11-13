package model

import "time"

// Menu GORM菜单模型
type Menu struct {
	ID          string    `gorm:"primaryKey;type:varchar(26)"`
	Name        string    `gorm:"not null;type:varchar(100)"`
	Path        string    `gorm:"not null;type:varchar(255)"`
	Icon        string    `gorm:"type:varchar(100)"`
	ParentID    *string   `gorm:"index;type:varchar(26)"` // NULL表示根菜单
	SortOrder   int       `gorm:"default:0"`
	Type        string    `gorm:"not null;type:varchar(20)"` // dir, menu, link
	IsVisible   bool      `gorm:"default:true"`
	Component   string    `gorm:"type:varchar(255)"` // 前端组件路径
	Permission  string    `gorm:"type:varchar(100)"` // 关联的权限码
	Description string    `gorm:"type:varchar(500)"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	// 关联关系
	Parent   *Menu  `gorm:"foreignKey:ParentID"`
	Children []Menu `gorm:"foreignKey:ParentID"`
	Roles    []Role `gorm:"many2many:role_menus;"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menus"
}
