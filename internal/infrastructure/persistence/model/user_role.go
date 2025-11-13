package model

import "time"

// UserRole 用户-角色关联表
type UserRole struct {
	UserID    string    `gorm:"primaryKey;type:varchar(26)"`
	RoleID    string    `gorm:"primaryKey;type:varchar(26)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Role Role `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}
