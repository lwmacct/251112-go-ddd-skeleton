package user

import (
	"context"
	"errors"
)

// Service 用户领域服务
type Service struct {
	repo Repository
}

// NewService 创建用户领域服务
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// ValidateUserUniqueness 验证用户唯一性（邮箱不重复）
func (s *Service) ValidateUserUniqueness(ctx context.Context, email Email) error {
	exists, err := s.repo.ExistsByEmail(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return ErrEmailAlreadyExists
	}
	return nil
}

// CanChangePassword 验证是否可以修改密码
func (s *Service) CanChangePassword(ctx context.Context, userID string, oldPasswordHash string) (bool, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}

	if !user.IsActive {
		return false, errors.New("user is not active")
	}

	// 验证旧密码（这里简化处理，实际应该在 application 层做）
	if user.Password.Hash() != oldPasswordHash {
		return false, errors.New("old password is incorrect")
	}

	return true, nil
}
