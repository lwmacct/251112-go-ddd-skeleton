package auth

import (
	"context"
	"math/rand"
	"time"

	"github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/auth"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/user"
	"github.com/oklog/ulid/v2"
)

// Login 登录（命令）
func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// 根据邮箱查找用户
	email, err := user.NewEmail(req.Email)
	if err != nil {
		return nil, err
	}

	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, auth.ErrInvalidCredentials
	}

	// 验证密码
	if err := s.passwordHasher.Compare(u.Password.Hash(), req.Password); err != nil {
		return nil, auth.ErrInvalidCredentials
	}

	// 检查用户是否激活
	if !u.IsActive {
		return nil, user.ErrUserNotActive
	}

	// 生成访问令牌
	accessToken, expiresIn, err := s.tokenIssuer.GenerateAccessToken(u.ID)
	if err != nil {
		return nil, err
	}

	// 生成刷新令牌
	refreshToken, err := s.tokenIssuer.GenerateRefreshToken(u.ID)
	if err != nil {
		return nil, err
	}

	// 创建会话
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	sessionID := ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

	session, err := auth.NewSession(u.ID, refreshToken, "", "", auth.GenerateSessionExpiryDate())
	if err != nil {
		return nil, err
	}
	session.ID = sessionID

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
	}, nil
}

// Logout 登出（命令）
func (s *Service) Logout(ctx context.Context, token string) error {
	return s.sessionRepo.Delete(ctx, token)
}

// RefreshToken 刷新令牌（命令）
func (s *Service) RefreshToken(ctx context.Context, req RefreshTokenRequest) (*LoginResponse, error) {
	// 验证刷新令牌
	userID, err := s.tokenIssuer.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// 验证会话
	_, err = s.authService.ValidateSession(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// 生成新的访问令牌
	accessToken, expiresIn, err := s.tokenIssuer.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	// 生成新的刷新令牌
	refreshToken, err := s.tokenIssuer.GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	// 删除旧会话
	if err := s.sessionRepo.Delete(ctx, req.RefreshToken); err != nil {
		return nil, err
	}

	// 创建新会话
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	sessionID := ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

	session, err := auth.NewSession(userID, refreshToken, "", "", auth.GenerateSessionExpiryDate())
	if err != nil {
		return nil, err
	}
	session.ID = sessionID

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
	}, nil
}

// Enable2FA 启用双因素认证（命令）
func (s *Service) Enable2FA(ctx context.Context, userID, email string) (*Enable2FAResponse, error) {
	// 生成TOTP密钥
	secret, qrCode, err := s.totpGenerator.Generate(email)
	if err != nil {
		return nil, err
	}

	// 创建或更新双因素认证
	tf, err := auth.NewTwoFactor(userID, secret)
	if err != nil {
		return nil, err
	}

	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	tf.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

	// 尝试查找现有的2FA配置
	existing, err := s.tfRepo.FindByUserID(ctx, userID)
	if err == nil {
		// 更新现有配置
		existing.Secret = secret
		if err := s.tfRepo.Update(ctx, existing); err != nil {
			return nil, err
		}
	} else {
		// 创建新配置
		if err := s.tfRepo.Create(ctx, tf); err != nil {
			return nil, err
		}
	}

	return &Enable2FAResponse{
		Secret:  secret,
		QRCode:  qrCode,
		Enabled: false,
	}, nil
}

// Verify2FA 验证双因素认证（命令）
func (s *Service) Verify2FA(ctx context.Context, userID string, req Verify2FARequest) error {
	// 查找双因素认证配置
	tf, err := s.tfRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// 验证TOTP代码
	if !s.totpGenerator.Validate(tf.Secret, req.Code) {
		return auth.ErrInvalidTOTPCode
	}

	// 启用双因素认证
	tf.Enable()
	return s.tfRepo.Update(ctx, tf)
}

// Disable2FA 禁用双因素认证（命令）
func (s *Service) Disable2FA(ctx context.Context, userID string) error {
	return s.tfRepo.Delete(ctx, userID)
}

// CreatePAT 创建个人访问令牌（命令）
func (s *Service) CreatePAT(ctx context.Context, userID string, req CreatePATRequest) (*CreatePATResponse, error) {
	// 生成令牌
	token, err := s.tokenIssuer.GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	// 计算过期时间
	var expiresAt *time.Time
	if req.ExpiresIn > 0 {
		expiry := time.Now().AddDate(0, 0, req.ExpiresIn)
		expiresAt = &expiry
	}

	// 创建PAT
	pat, err := auth.NewPAT(userID, req.Name, token, req.Scopes, expiresAt)
	if err != nil {
		return nil, err
	}

	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	pat.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

	if err := s.patRepo.Create(ctx, pat); err != nil {
		return nil, err
	}

	return &CreatePATResponse{
		ID:        pat.ID,
		Name:      pat.Name,
		Token:     pat.Token,
		Scopes:    pat.Scopes,
		ExpiresAt: pat.ExpiresAt,
		CreatedAt: pat.CreatedAt,
	}, nil
}

// RevokePAT 撤销个人访问令牌（命令）
func (s *Service) RevokePAT(ctx context.Context, patID string) error {
	return s.patRepo.Delete(ctx, patID)
}

// RevokeAllUserSessions 撤销用户所有会话（命令）
func (s *Service) RevokeAllUserSessions(ctx context.Context, userID string) error {
	return s.authService.RevokeAllUserSessions(ctx, userID)
}
