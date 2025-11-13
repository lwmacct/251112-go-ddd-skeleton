package repository

import (
	"context"
	"errors"
	"time"

	"github.com/yourusername/go-ddd-skeleton/internal/domain/auth"
	"github.com/yourusername/go-ddd-skeleton/internal/infrastructure/persistence/mapper"
	"github.com/yourusername/go-ddd-skeleton/internal/infrastructure/persistence/model"
	"gorm.io/gorm"
)

// TwoFactorRepository 双因素认证仓储实现
type TwoFactorRepository struct {
	db *gorm.DB
}

// NewTwoFactorRepository 创建双因素认证仓储
func NewTwoFactorRepository(db *gorm.DB) auth.TwoFactorRepository {
	return &TwoFactorRepository{db: db}
}

func (r *TwoFactorRepository) Create(ctx context.Context, tf *auth.TwoFactor) error {
	m := mapper.TwoFactorToModel(tf)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *TwoFactorRepository) Update(ctx context.Context, tf *auth.TwoFactor) error {
	m := mapper.TwoFactorToModel(tf)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *TwoFactorRepository) FindByUserID(ctx context.Context, userID string) (*auth.TwoFactor, error) {
	var m model.TwoFactor
	if err := r.db.WithContext(ctx).First(&m, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, auth.ErrTwoFactorNotFound
		}
		return nil, err
	}
	return mapper.TwoFactorToDomain(&m), nil
}

func (r *TwoFactorRepository) Delete(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Delete(&model.TwoFactor{}, "user_id = ?", userID).Error
}

// PATRepository PAT仓储实现
type PATRepository struct {
	db *gorm.DB
}

// NewPATRepository 创建PAT仓储
func NewPATRepository(db *gorm.DB) auth.PATRepository {
	return &PATRepository{db: db}
}

func (r *PATRepository) Create(ctx context.Context, pat *auth.PAT) error {
	m := mapper.PATToModel(pat)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PATRepository) FindByID(ctx context.Context, id string) (*auth.PAT, error) {
	var m model.PersonalAccessToken
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, auth.ErrPATNotFound
		}
		return nil, err
	}
	return mapper.PATToDomain(&m), nil
}

func (r *PATRepository) FindByToken(ctx context.Context, token string) (*auth.PAT, error) {
	var m model.PersonalAccessToken
	if err := r.db.WithContext(ctx).First(&m, "token = ?", token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, auth.ErrPATNotFound
		}
		return nil, err
	}
	return mapper.PATToDomain(&m), nil
}

func (r *PATRepository) ListByUserID(ctx context.Context, userID string) ([]*auth.PAT, error) {
	var models []model.PersonalAccessToken
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil, err
	}

	pats := make([]*auth.PAT, len(models))
	for i, m := range models {
		pats[i] = mapper.PATToDomain(&m)
	}
	return pats, nil
}

func (r *PATRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.PersonalAccessToken{}, "id = ?", id).Error
}

func (r *PATRepository) Update(ctx context.Context, pat *auth.PAT) error {
	m := mapper.PATToModel(pat)
	return r.db.WithContext(ctx).Save(m).Error
}

// SessionRepository Session仓储实现
type SessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository 创建Session仓储
func NewSessionRepository(db *gorm.DB) auth.SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, session *auth.Session) error {
	m := mapper.SessionToModel(session)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *SessionRepository) FindByToken(ctx context.Context, token string) (*auth.Session, error) {
	var m model.Session
	if err := r.db.WithContext(ctx).First(&m, "token = ?", token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, auth.ErrSessionNotFound
		}
		return nil, err
	}
	return mapper.SessionToDomain(&m), nil
}

func (r *SessionRepository) ListByUserID(ctx context.Context, userID string) ([]*auth.Session, error) {
	var models []model.Session
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil, err
	}

	sessions := make([]*auth.Session, len(models))
	for i, m := range models {
		sessions[i] = mapper.SessionToDomain(&m)
	}
	return sessions, nil
}

func (r *SessionRepository) Delete(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Delete(&model.Session{}, "token = ?", token).Error
}

func (r *SessionRepository) DeleteByUserID(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Delete(&model.Session{}, "user_id = ?", userID).Error
}

func (r *SessionRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&model.Session{}).Error
}
