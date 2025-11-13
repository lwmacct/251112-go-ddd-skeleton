package auth

import (
	"errors"
	"time"
)

// TwoFactor 双因素认证实体
type TwoFactor struct {
	ID        string
	UserID    string
	Secret    string
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewTwoFactor 创建双因素认证
func NewTwoFactor(userID, secret string) (*TwoFactor, error) {
	if userID == "" {
		return nil, errors.New("userID cannot be empty")
	}
	if secret == "" {
		return nil, errors.New("secret cannot be empty")
	}

	return &TwoFactor{
		UserID:    userID,
		Secret:    secret,
		Enabled:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// Enable 启用双因素认证
func (tf *TwoFactor) Enable() {
	tf.Enabled = true
	tf.UpdatedAt = time.Now()
}

// Disable 禁用双因素认证
func (tf *TwoFactor) Disable() {
	tf.Enabled = false
	tf.UpdatedAt = time.Now()
}

// PAT Personal Access Token 个人访问令牌实体
type PAT struct {
	ID        string
	UserID    string
	Name      string
	Token     string
	Scopes    []string
	ExpiresAt *time.Time
	LastUsed  *time.Time
	CreatedAt time.Time
}

// NewPAT 创建个人访问令牌
func NewPAT(userID, name, token string, scopes []string, expiresAt *time.Time) (*PAT, error) {
	if userID == "" {
		return nil, errors.New("userID cannot be empty")
	}
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	if token == "" {
		return nil, errors.New("token cannot be empty")
	}

	return &PAT{
		UserID:    userID,
		Name:      name,
		Token:     token,
		Scopes:    scopes,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}, nil
}

// IsExpired 判断令牌是否过期
func (p *PAT) IsExpired() bool {
	if p.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*p.ExpiresAt)
}

// MarkUsed 标记令牌已使用
func (p *PAT) MarkUsed() {
	now := time.Now()
	p.LastUsed = &now
}

// Session 会话实体
type Session struct {
	ID        string
	UserID    string
	Token     string
	IP        string
	UserAgent string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// NewSession 创建会话
func NewSession(userID, token, ip, userAgent string, expiresAt time.Time) (*Session, error) {
	if userID == "" {
		return nil, errors.New("userID cannot be empty")
	}
	if token == "" {
		return nil, errors.New("token cannot be empty")
	}

	return &Session{
		UserID:    userID,
		Token:     token,
		IP:        ip,
		UserAgent: userAgent,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}, nil
}

// IsExpired 判断会话是否过期
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsValid 判断会话是否有效
func (s *Session) IsValid() bool {
	return !s.IsExpired()
}

