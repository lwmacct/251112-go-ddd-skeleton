package mapper

import (
	"github.com/yourusername/go-ddd-skeleton/internal/domain/user"
	"github.com/yourusername/go-ddd-skeleton/internal/infrastructure/persistence/model"
)

// UserToModel 将User领域实体转换为GORM模型
func UserToModel(u *user.User) *model.User {
	return &model.User{
		ID:        u.ID,
		Email:     u.Email.String(),
		Password:  u.Password.Hash(),
		Username:  u.Username,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// UserToDomain 将GORM模型转换为User领域实体
func UserToDomain(m *model.User) (*user.User, error) {
	email, err := user.NewEmail(m.Email)
	if err != nil {
		return nil, err
	}

	return &user.User{
		ID:        m.ID,
		Email:     email,
		Password:  user.NewPasswordFromHash(m.Password),
		Username:  m.Username,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}
