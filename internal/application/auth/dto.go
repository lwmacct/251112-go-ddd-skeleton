package auth

import "time"

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Enable2FAResponse 启用2FA响应
type Enable2FAResponse struct {
	Secret  string `json:"secret"`
	QRCode  string `json:"qr_code"`
	Enabled bool   `json:"enabled"`
}

// Verify2FARequest 验证2FA请求
type Verify2FARequest struct {
	Code string `json:"code" validate:"required,len=6"`
}

// CreatePATRequest 创建个人访问令牌请求
type CreatePATRequest struct {
	Name      string   `json:"name" validate:"required"`
	Scopes    []string `json:"scopes"`
	ExpiresIn int      `json:"expires_in"` // 天数，0表示永不过期
}

// CreatePATResponse 创建个人访问令牌响应
type CreatePATResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Token     string     `json:"token"`
	Scopes    []string   `json:"scopes"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// PATDTO 个人访问令牌DTO
type PATDTO struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Scopes    []string   `json:"scopes"`
	LastUsed  *time.Time `json:"last_used,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// SessionDTO 会话DTO
type SessionDTO struct {
	ID        string    `json:"id"`
	Token     string    `json:"token"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// TwoFactorDTO 双因素认证DTO
type TwoFactorDTO struct {
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}
