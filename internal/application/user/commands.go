package user

import (
	"context"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/yourusername/go-ddd-skeleton/internal/domain/user"
)

// CreateUser 创建用户（命令）
func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (*UserDTO, error) {
	// 验证邮箱唯一性
	email, err := user.NewEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if err := s.userService.ValidateUserUniqueness(ctx, email); err != nil {
		return nil, err
	}

	// 哈希密码
	hashedPassword, err := s.passwordHasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	// 创建用户
	u, err := user.NewUser(req.Email, hashedPassword, req.Username)
	if err != nil {
		return nil, err
	}

	// 生成ID
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	u.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

	// 保存到数据库
	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}

	return domainToDTO(u), nil
}

// UpdateUser 更新用户（命令）
func (s *Service) UpdateUser(ctx context.Context, userID string, req UpdateUserRequest) (*UserDTO, error) {
	// 查找用户
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 更新用户信息
	if err := u.UpdateProfile(req.Username); err != nil {
		return nil, err
	}

	// 保存更新
	if err := s.userRepo.Update(ctx, u); err != nil {
		return nil, err
	}

	return domainToDTO(u), nil
}

// DeleteUser 删除用户（命令）
func (s *Service) DeleteUser(ctx context.Context, userID string) error {
	// 检查用户是否存在
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// 删除用户
	return s.userRepo.Delete(ctx, userID)
}

// ChangePassword 修改密码（命令）
func (s *Service) ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error {
	// 查找用户
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// 验证旧密码
	if err := s.passwordHasher.Compare(u.Password.Hash(), req.OldPassword); err != nil {
		return user.ErrInvalidPassword
	}

	// 哈希新密码
	hashedPassword, err := s.passwordHasher.Hash(req.NewPassword)
	if err != nil {
		return err
	}

	// 修改密码
	if err := u.ChangePassword(hashedPassword); err != nil {
		return err
	}

	// 保存更新
	return s.userRepo.Update(ctx, u)
}

// DeactivateUser 停用用户（命令）
func (s *Service) DeactivateUser(ctx context.Context, userID string) error {
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	u.Deactivate()
	return s.userRepo.Update(ctx, u)
}

// ActivateUser 激活用户（命令）
func (s *Service) ActivateUser(ctx context.Context, userID string) error {
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	u.Activate()
	return s.userRepo.Update(ctx, u)
}
