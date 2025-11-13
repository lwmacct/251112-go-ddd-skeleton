package validator

import (
	"github.com/go-playground/validator/v10"
	"sync"
)

var (
	once     sync.Once
	validate *validator.Validate
)

// GetValidator 获取验证器单例
func GetValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
		// 可以在这里注册自定义验证规则
		registerCustomValidations()
	})
	return validate
}

// registerCustomValidations 注册自定义验证规则
func registerCustomValidations() {
	// 示例：注册自定义验证规则
	// validate.RegisterValidation("custom", customValidationFunc)
}

// Validate 验证结构体
func Validate(s interface{}) error {
	return GetValidator().Struct(s)
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// FormatValidationErrors 格式化验证错误
func FormatValidationErrors(err error) []ValidationError {
	var errors []ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors = append(errors, ValidationError{
				Field:   e.Field(),
				Message: formatErrorMessage(e),
			})
		}
	}

	return errors
}

// formatErrorMessage 格式化错误消息
func formatErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short"
	case "max":
		return "Value is too long"
	case "gte":
		return "Value must be greater than or equal to " + e.Param()
	case "lte":
		return "Value must be less than or equal to " + e.Param()
	default:
		return "Invalid value"
	}
}

