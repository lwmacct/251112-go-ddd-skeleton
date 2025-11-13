package model

import "time"

// TwoFactor GORM双因素认证模型
type TwoFactor struct {
	ID        string    `gorm:"primaryKey;type:varchar(26)"`
	UserID    string    `gorm:"uniqueIndex;not null;type:varchar(26)"`
	Secret    string    `gorm:"not null;type:varchar(255)"`
	Enabled   bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	User User `gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (TwoFactor) TableName() string {
	return "two_factor_auth"
}
