package user

import "errors"

var (
	// ErrUserNotFound 用户未找到
	ErrUserNotFound = errors.New("user not found")

	// ErrEmailAlreadyExists 邮箱已存在
	ErrEmailAlreadyExists = errors.New("email already exists")

	// ErrInvalidEmail 无效的邮箱格式
	ErrInvalidEmail = errors.New("invalid email format")

	// ErrInvalidPassword 无效的密码
	ErrInvalidPassword = errors.New("invalid password")

	// ErrUserNotActive 用户未激活
	ErrUserNotActive = errors.New("user is not active")

	// ErrPasswordTooShort 密码太短
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
)
