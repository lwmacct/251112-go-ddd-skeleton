package auth

import (
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/auth"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/user"
)

// TokenIssuer 令牌生成器接口（端口）
type TokenIssuer interface {
	GenerateAccessToken(userID string) (string, int, error)
	GenerateRefreshToken(userID string) (string, error)
	ValidateToken(token string) (string, error)
}

// PasswordHasher 密码哈希接口（端口）
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}

// TOTPGenerator TOTP生成器接口（端口）
type TOTPGenerator interface {
	Generate(accountName string) (secret, qrCode string, err error)
	Validate(secret, code string) bool
}

// Service 认证应用服务
type Service struct {
	userRepo    user.Repository
	tfRepo      auth.TwoFactorRepository
	patRepo     auth.PATRepository
	sessionRepo auth.SessionRepository
	authService *auth.Service

	tokenIssuer    TokenIssuer
	passwordHasher PasswordHasher
	totpGenerator  TOTPGenerator
}

// NewService 创建认证应用服务
func NewService(
	userRepo user.Repository,
	tfRepo auth.TwoFactorRepository,
	patRepo auth.PATRepository,
	sessionRepo auth.SessionRepository,
	authService *auth.Service,
	tokenIssuer TokenIssuer,
	passwordHasher PasswordHasher,
	totpGenerator TOTPGenerator,
) *Service {
	return &Service{
		userRepo:       userRepo,
		tfRepo:         tfRepo,
		patRepo:        patRepo,
		sessionRepo:    sessionRepo,
		authService:    authService,
		tokenIssuer:    tokenIssuer,
		passwordHasher: passwordHasher,
		totpGenerator:  totpGenerator,
	}
}
