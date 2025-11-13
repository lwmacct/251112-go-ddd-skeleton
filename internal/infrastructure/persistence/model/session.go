package model

import "time"

// Session GORM会话模型
type Session struct {
	ID        string `gorm:"primaryKey;type:varchar(26)"`
	UserID    string `gorm:"index;not null;type:varchar(26)"`
	Token     string `gorm:"uniqueIndex;not null;type:varchar(255)"`
	IP        string `gorm:"type:varchar(45)"`
	UserAgent string `gorm:"type:varchar(500)"`
	ExpiresAt time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`

	User User `gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (Session) TableName() string {
	return "sessions"
}
