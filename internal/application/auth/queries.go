package auth

import (
	"context"
)

// GetUserSessions 获取用户会话列表（查询）
func (s *Service) GetUserSessions(ctx context.Context, userID string) ([]*SessionDTO, error) {
	sessions, err := s.sessionRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	dtos := make([]*SessionDTO, len(sessions))
	for i, session := range sessions {
		dtos[i] = &SessionDTO{
			ID:        session.ID,
			Token:     session.Token,
			IP:        session.IP,
			UserAgent: session.UserAgent,
			ExpiresAt: session.ExpiresAt,
			CreatedAt: session.CreatedAt,
		}
	}

	return dtos, nil
}

// ListUserPATs 列出用户的个人访问令牌（查询）
func (s *Service) ListUserPATs(ctx context.Context, userID string) ([]*PATDTO, error) {
	pats, err := s.patRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	dtos := make([]*PATDTO, len(pats))
	for i, pat := range pats {
		dtos[i] = &PATDTO{
			ID:        pat.ID,
			Name:      pat.Name,
			Scopes:    pat.Scopes,
			LastUsed:  pat.LastUsed,
			ExpiresAt: pat.ExpiresAt,
			CreatedAt: pat.CreatedAt,
		}
	}

	return dtos, nil
}

// Get2FAStatus 获取双因素认证状态（查询）
func (s *Service) Get2FAStatus(ctx context.Context, userID string) (*TwoFactorDTO, error) {
	tf, err := s.tfRepo.FindByUserID(ctx, userID)
	if err != nil {
		return &TwoFactorDTO{
			Enabled: false,
		}, nil
	}

	return &TwoFactorDTO{
		Enabled:   tf.Enabled,
		CreatedAt: tf.CreatedAt,
	}, nil
}

// ValidatePAT 验证个人访问令牌（查询）
func (s *Service) ValidatePAT(ctx context.Context, token string) (string, error) {
	pat, err := s.authService.ValidatePAT(ctx, token)
	if err != nil {
		return "", err
	}
	return pat.UserID, nil
}
