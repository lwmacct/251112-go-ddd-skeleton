package errors

import (
	"errors"
	"fmt"
)

// ErrorCode 错误码类型
type ErrorCode string

const (
	// 通用错误码
	CodeInternal        ErrorCode = "INTERNAL_ERROR"
	CodeBadRequest      ErrorCode = "BAD_REQUEST"
	CodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	CodeForbidden       ErrorCode = "FORBIDDEN"
	CodeNotFound        ErrorCode = "NOT_FOUND"
	CodeConflict        ErrorCode = "CONFLICT"
	CodeValidation      ErrorCode = "VALIDATION_ERROR"
	CodeTooManyRequests ErrorCode = "TOO_MANY_REQUESTS"

	// 用户相关错误码
	CodeUserNotFound       ErrorCode = "USER_NOT_FOUND"
	CodeUserAlreadyExists  ErrorCode = "USER_ALREADY_EXISTS"
	CodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	CodeInvalidEmail       ErrorCode = "INVALID_EMAIL"
	CodeInvalidPassword    ErrorCode = "INVALID_PASSWORD"
	CodeUserNotActive      ErrorCode = "USER_NOT_ACTIVE"

	// 认证相关错误码
	CodeInvalidToken    ErrorCode = "INVALID_TOKEN"
	CodeTokenExpired    ErrorCode = "TOKEN_EXPIRED"
	CodeSessionNotFound ErrorCode = "SESSION_NOT_FOUND"
	CodeSessionExpired  ErrorCode = "SESSION_EXPIRED"
	CodeInvalidTOTPCode ErrorCode = "INVALID_TOTP_CODE"
	Code2FANotEnabled   ErrorCode = "2FA_NOT_ENABLED"

	// 订单相关错误码
	CodeOrderNotFound      ErrorCode = "ORDER_NOT_FOUND"
	CodeInvalidOrderStatus ErrorCode = "INVALID_ORDER_STATUS"
	CodePaymentFailed      ErrorCode = "PAYMENT_FAILED"
	CodeShipmentNotFound   ErrorCode = "SHIPMENT_NOT_FOUND"
)

// AppError 应用错误结构
type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap 支持 errors.Unwrap
func (e *AppError) Unwrap() error {
	return e.Err
}

// New 创建新的应用错误
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap 包装错误
func Wrap(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Is 判断错误是否匹配特定错误码
func Is(err error, code ErrorCode) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}

// GetCode 获取错误码
func GetCode(err error) ErrorCode {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return CodeInternal
}

// 预定义的常用错误
var (
	ErrInternal        = New(CodeInternal, "Internal server error")
	ErrBadRequest      = New(CodeBadRequest, "Bad request")
	ErrUnauthorized    = New(CodeUnauthorized, "Unauthorized")
	ErrForbidden       = New(CodeForbidden, "Forbidden")
	ErrNotFound        = New(CodeNotFound, "Resource not found")
	ErrConflict        = New(CodeConflict, "Resource conflict")
	ErrValidation      = New(CodeValidation, "Validation error")
	ErrTooManyRequests = New(CodeTooManyRequests, "Too many requests")
)
