package user

import (
	"github.com/yourusername/go-ddd-skeleton/internal/domain/user"
)

// PasswordHasher 密码哈希接口（端口）
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}

// Service 用户应用服务
type Service struct {
	userRepo       user.Repository
	userService    *user.Service
	passwordHasher PasswordHasher
}

// NewService 创建用户应用服务
func NewService(userRepo user.Repository, userService *user.Service, passwordHasher PasswordHasher) *Service {
	return &Service{
		userRepo:       userRepo,
		userService:    userService,
		passwordHasher: passwordHasher,
	}
}

// domainToDTO 将领域实体转换为DTO
func domainToDTO(u *user.User) *UserDTO {
	return &UserDTO{
		ID:        u.ID,
		Email:     u.Email.String(),
		Username:  u.Username,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
