package auth

import (
	"context"
	"time"
)

// Service 认证领域服务
type Service struct {
	tfRepo      TwoFactorRepository
	patRepo     PATRepository
	sessionRepo SessionRepository
}

// NewService 创建认证领域服务
func NewService(tfRepo TwoFactorRepository, patRepo PATRepository, sessionRepo SessionRepository) *Service {
	return &Service{
		tfRepo:      tfRepo,
		patRepo:     patRepo,
		sessionRepo: sessionRepo,
	}
}

// ValidatePAT 验证个人访问令牌
func (s *Service) ValidatePAT(ctx context.Context, token string) (*PAT, error) {
	pat, err := s.patRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if pat.IsExpired() {
		return nil, ErrTokenExpired
	}

	// 标记令牌已使用
	pat.MarkUsed()
	if err := s.patRepo.Update(ctx, pat); err != nil {
		return nil, err
	}

	return pat, nil
}

// ValidateSession 验证会话
func (s *Service) ValidateSession(ctx context.Context, token string) (*Session, error) {
	session, err := s.sessionRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if !session.IsValid() {
		return nil, ErrSessionExpired
	}

	return session, nil
}

// CleanupExpiredSessions 清理过期会话
func (s *Service) CleanupExpiredSessions(ctx context.Context) error {
	return s.sessionRepo.DeleteExpired(ctx)
}

// RevokeAllUserSessions 撤销用户的所有会话
func (s *Service) RevokeAllUserSessions(ctx context.Context, userID string) error {
	return s.sessionRepo.DeleteByUserID(ctx, userID)
}

// GeneratePATExpiryDate 生成PAT过期时间（默认90天）
func GeneratePATExpiryDate() *time.Time {
	expiry := time.Now().AddDate(0, 0, 90)
	return &expiry
}

// GenerateSessionExpiryDate 生成会话过期时间（默认7天）
func GenerateSessionExpiryDate() time.Time {
	return time.Now().AddDate(0, 0, 7)
}
