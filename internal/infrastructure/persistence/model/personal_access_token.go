package model

import "time"

// PersonalAccessToken GORM个人访问令牌模型
type PersonalAccessToken struct {
	ID        string `gorm:"primaryKey;type:varchar(26)"`
	UserID    string `gorm:"index;not null;type:varchar(26)"`
	Name      string `gorm:"not null;type:varchar(100)"`
	Token     string `gorm:"uniqueIndex;not null;type:varchar(255)"`
	Scopes    string `gorm:"type:text"` // JSON array
	ExpiresAt *time.Time
	LastUsed  *time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`

	User User `gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (PersonalAccessToken) TableName() string {
	return "personal_access_tokens"
}
