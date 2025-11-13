package auth

import "context"

// TwoFactorRepository 双因素认证仓储接口
type TwoFactorRepository interface {
	// Create 创建双因素认证
	Create(ctx context.Context, tf *TwoFactor) error

	// Update 更新双因素认证
	Update(ctx context.Context, tf *TwoFactor) error

	// FindByUserID 根据用户ID查找双因素认证
	FindByUserID(ctx context.Context, userID string) (*TwoFactor, error)

	// Delete 删除双因素认证
	Delete(ctx context.Context, userID string) error
}

// PATRepository 个人访问令牌仓储接口
type PATRepository interface {
	// Create 创建个人访问令牌
	Create(ctx context.Context, pat *PAT) error

	// FindByID 根据ID查找令牌
	FindByID(ctx context.Context, id string) (*PAT, error)

	// FindByToken 根据令牌查找
	FindByToken(ctx context.Context, token string) (*PAT, error)

	// ListByUserID 列出用户的所有令牌
	ListByUserID(ctx context.Context, userID string) ([]*PAT, error)

	// Delete 删除令牌
	Delete(ctx context.Context, id string) error

	// Update 更新令牌
	Update(ctx context.Context, pat *PAT) error
}

// SessionRepository 会话仓储接口
type SessionRepository interface {
	// Create 创建会话
	Create(ctx context.Context, session *Session) error

	// FindByToken 根据令牌查找会话
	FindByToken(ctx context.Context, token string) (*Session, error)

	// ListByUserID 列出用户的所有会话
	ListByUserID(ctx context.Context, userID string) ([]*Session, error)

	// Delete 删除会话
	Delete(ctx context.Context, token string) error

	// DeleteByUserID 删除用户的所有会话
	DeleteByUserID(ctx context.Context, userID string) error

	// DeleteExpired 删除过期会话
	DeleteExpired(ctx context.Context) error
}
