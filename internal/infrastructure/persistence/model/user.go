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
