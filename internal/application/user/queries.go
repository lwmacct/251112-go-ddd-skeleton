package user

import (
	"context"

	"github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/user"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/shared/pagination"
)

// GetUser 获取用户（查询）
func (s *Service) GetUser(ctx context.Context, userID string) (*UserDTO, error) {
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return domainToDTO(u), nil
}

// GetUserByEmail 根据邮箱获取用户（查询）
func (s *Service) GetUserByEmail(ctx context.Context, email string) (*UserDTO, error) {
	e, err := user.NewEmail(email)
	if err != nil {
		return nil, err
	}

	u, err := s.userRepo.FindByEmail(ctx, e)
	if err != nil {
		return nil, err
	}

	return domainToDTO(u), nil
}

// ListUsers 列出用户（查询）
func (s *Service) ListUsers(ctx context.Context, req ListUsersRequest) (*ListUsersResponse, error) {
	// 解析分页参数
	offset, limit := pagination.ParsePaginationParams(req.Page, req.PageSize)

	// 查询用户列表
	users, total, err := s.userRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	userDTOs := make([]*UserDTO, len(users))
	for i, u := range users {
		userDTOs[i] = domainToDTO(u)
	}

	// 创建分页对象
	pg := pagination.NewPagination(req.Page, req.PageSize, total)

	return &ListUsersResponse{
		Users:      userDTOs,
		Total:      total,
		Page:       pg.Page,
		PageSize:   pg.PageSize,
		TotalPages: pg.TotalPages,
	}, nil
}

// UserExists 检查用户是否存在（查询）
func (s *Service) UserExists(ctx context.Context, email string) (bool, error) {
	e, err := user.NewEmail(email)
	if err != nil {
		return false, err
	}

	return s.userRepo.ExistsByEmail(ctx, e)
}
