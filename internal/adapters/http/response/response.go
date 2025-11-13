package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/shared/errors"
)

// Response 统一响应结构
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo 错误信息
type ErrorInfo struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// Created 创建成功响应
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

// NoContent 无内容响应
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error 错误响应
func Error(c *gin.Context, err error) {
	code := errors.GetCode(err)
	statusCode := getHTTPStatus(code)

	c.JSON(statusCode, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    string(code),
			Message: err.Error(),
		},
	})
}

// ErrorWithDetails 带详情的错误响应
func ErrorWithDetails(c *gin.Context, err error, details interface{}) {
	code := errors.GetCode(err)
	statusCode := getHTTPStatus(code)

	c.JSON(statusCode, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    string(code),
			Message: err.Error(),
			Details: details,
		},
	})
}

// getHTTPStatus 根据错误码获取HTTP状态码
func getHTTPStatus(code errors.ErrorCode) int {
	switch code {
	case errors.CodeBadRequest, errors.CodeValidation:
		return http.StatusBadRequest
	case errors.CodeUnauthorized, errors.CodeInvalidToken, errors.CodeTokenExpired:
		return http.StatusUnauthorized
	case errors.CodeForbidden:
		return http.StatusForbidden
	case errors.CodeNotFound:
		return http.StatusNotFound
	case errors.CodeConflict:
		return http.StatusConflict
	case errors.CodeTooManyRequests:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
