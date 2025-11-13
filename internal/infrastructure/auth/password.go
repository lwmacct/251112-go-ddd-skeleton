package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher Bcrypt密码哈希器
type PasswordHasher struct {
	cost int
}

// NewPasswordHasher 创建密码哈希器
func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{
		cost: bcrypt.DefaultCost,
	}
}

// Hash 哈希密码
func (h *PasswordHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	return string(bytes), err
}

// Compare 比较密码
func (h *PasswordHasher) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

