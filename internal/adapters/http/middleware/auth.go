package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/adapters/http/response"
	apperrors "github.com/lwmacct/251112-go-ddd-skeleton/internal/shared/errors"
)

// TokenValidator 令牌验证器接口
type TokenValidator interface {
	ValidateToken(token string) (string, error)
}

var tokenValidator TokenValidator

// SetTokenValidator 设置令牌验证器
func SetTokenValidator(validator TokenValidator) {
	tokenValidator = validator
}

// Auth 认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		// 提取Bearer令牌
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		token := parts[1]
		userID, err := tokenValidator.ValidateToken(token)
		if err != nil {
			response.Error(c, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		// 将用户ID存入context
		c.Set("userID", userID)
		c.Next()
	}
}
