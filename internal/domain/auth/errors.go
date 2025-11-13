package auth

import "errors"

var (
	// ErrTwoFactorNotFound 双因素认证未找到
	ErrTwoFactorNotFound = errors.New("two factor authentication not found")

	// ErrTwoFactorNotEnabled 双因素认证未启用
	ErrTwoFactorNotEnabled = errors.New("two factor authentication not enabled")

	// ErrInvalidTOTPCode 无效的TOTP验证码
	ErrInvalidTOTPCode = errors.New("invalid TOTP code")

	// ErrPATNotFound 个人访问令牌未找到
	ErrPATNotFound = errors.New("personal access token not found")

	// ErrTokenExpired 令牌已过期
	ErrTokenExpired = errors.New("token expired")

	// ErrInvalidToken 无效的令牌
	ErrInvalidToken = errors.New("invalid token")

	// ErrSessionNotFound 会话未找到
	ErrSessionNotFound = errors.New("session not found")

	// ErrSessionExpired 会话已过期
	ErrSessionExpired = errors.New("session expired")

	// ErrInvalidCredentials 无效的凭证
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrUnauthorized 未授权
	ErrUnauthorized = errors.New("unauthorized")
)

