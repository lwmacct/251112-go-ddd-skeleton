package model

import "time"

// User GORM用户模型
type User struct {
	ID        string    `gorm:"primaryKey;type:varchar(26)"`
	Email     string    `gorm:"uniqueIndex;not null;type:varchar(255)"`
	Password  string    `gorm:"not null;type:varchar(255)"`
	Username  string    `gorm:"not null;type:varchar(100)"`
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// 关联关系
	TwoFactor *TwoFactor            `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	PATs      []PersonalAccessToken `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Sessions  []Session             `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Orders    []Order               `gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
