package user

import (
	"errors"
	"regexp"
	"time"
)

// User 用户聚合根
type User struct {
	ID        string
	Email     Email
	Password  Password
	Username  string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser 创建新用户
func NewUser(email, password, username string) (*User, error) {
	e, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	p, err := NewPassword(password)
	if err != nil {
		return nil, err
	}

	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	return &User{
		Email:     e,
		Password:  p,
		Username:  username,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// ChangePassword 修改密码
func (u *User) ChangePassword(newPassword string) error {
	p, err := NewPassword(newPassword)
	if err != nil {
		return err
	}
	u.Password = p
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateProfile 更新用户资料
func (u *User) UpdateProfile(username string) error {
	if username == "" {
		return errors.New("username cannot be empty")
	}
	u.Username = username
	u.UpdatedAt = time.Now()
	return nil
}

// Deactivate 停用用户
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// Activate 激活用户
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

// Email 值对象
type Email struct {
	value string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// NewEmail 创建邮箱值对象
func NewEmail(email string) (Email, error) {
	if email == "" {
		return Email{}, errors.New("email cannot be empty")
	}
	if !emailRegex.MatchString(email) {
		return Email{}, errors.New("invalid email format")
	}
	return Email{value: email}, nil
}

// String 返回邮箱字符串
func (e Email) String() string {
	return e.value
}

// Equals 判断两个邮箱是否相等
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// Password 值对象
type Password struct {
	hash string
}

// NewPassword 创建密码值对象（已哈希）
func NewPassword(password string) (Password, error) {
	if len(password) < 8 {
		return Password{}, errors.New("password must be at least 8 characters")
	}
	// 这里存储的应该是哈希值，实际哈希操作在 infrastructure 层
	return Password{hash: password}, nil
}

// NewPasswordFromHash 从哈希创建密码对象（用于从数据库恢复）
func NewPasswordFromHash(hash string) Password {
	return Password{hash: hash}
}

// Hash 返回密码哈希值
func (p Password) Hash() string {
	return p.hash
}
