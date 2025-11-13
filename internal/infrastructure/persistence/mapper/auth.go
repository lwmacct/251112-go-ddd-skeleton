package mapper

import (
	"encoding/json"

	"github.com/yourusername/go-ddd-skeleton/internal/domain/auth"
	"github.com/yourusername/go-ddd-skeleton/internal/infrastructure/persistence/model"
)

// TwoFactorToModel 转换双因素认证到模型
func TwoFactorToModel(tf *auth.TwoFactor) *model.TwoFactor {
	return &model.TwoFactor{
		ID:        tf.ID,
		UserID:    tf.UserID,
		Secret:    tf.Secret,
		Enabled:   tf.Enabled,
		CreatedAt: tf.CreatedAt,
		UpdatedAt: tf.UpdatedAt,
	}
}

// TwoFactorToDomain 转换模型到双因素认证
func TwoFactorToDomain(m *model.TwoFactor) *auth.TwoFactor {
	return &auth.TwoFactor{
		ID:        m.ID,
		UserID:    m.UserID,
		Secret:    m.Secret,
		Enabled:   m.Enabled,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// PATToModel 转换PAT到模型
func PATToModel(pat *auth.PAT) *model.PersonalAccessToken {
	scopesJSON, _ := json.Marshal(pat.Scopes)
	return &model.PersonalAccessToken{
		ID:        pat.ID,
		UserID:    pat.UserID,
		Name:      pat.Name,
		Token:     pat.Token,
		Scopes:    string(scopesJSON),
		ExpiresAt: pat.ExpiresAt,
		LastUsed:  pat.LastUsed,
		CreatedAt: pat.CreatedAt,
	}
}

// PATToDomain 转换模型到PAT
func PATToDomain(m *model.PersonalAccessToken) *auth.PAT {
	var scopes []string
	if m.Scopes != "" {
		json.Unmarshal([]byte(m.Scopes), &scopes)
	}

	return &auth.PAT{
		ID:        m.ID,
		UserID:    m.UserID,
		Name:      m.Name,
		Token:     m.Token,
		Scopes:    scopes,
		ExpiresAt: m.ExpiresAt,
		LastUsed:  m.LastUsed,
		CreatedAt: m.CreatedAt,
	}
}

// SessionToModel 转换会话到模型
func SessionToModel(s *auth.Session) *model.Session {
	return &model.Session{
		ID:        s.ID,
		UserID:    s.UserID,
		Token:     s.Token,
		IP:        s.IP,
		UserAgent: s.UserAgent,
		ExpiresAt: s.ExpiresAt,
		CreatedAt: s.CreatedAt,
	}
}

// SessionToDomain 转换模型到会话
func SessionToDomain(m *model.Session) *auth.Session {
	return &auth.Session{
		ID:        m.ID,
		UserID:    m.UserID,
		Token:     m.Token,
		IP:        m.IP,
		UserAgent: m.UserAgent,
		ExpiresAt: m.ExpiresAt,
		CreatedAt: m.CreatedAt,
	}
}
