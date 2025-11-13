package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/adapters/http/response"
	apperrors "github.com/lwmacct/251112-go-ddd-skeleton/internal/shared/errors"
)

// RoleChecker 角色检查器接口
type RoleChecker interface {
	IsAdmin(userID string) (bool, error)
}

var roleChecker RoleChecker

// SetRoleChecker 设置角色检查器
func SetRoleChecker(checker RoleChecker) {
	roleChecker = checker
}

// Admin 管理员权限中间件
// 必须在 Auth() 中间件之后使用
func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 context 获取用户 ID（由 Auth 中间件设置）
		userID, exists := c.Get("userID")
		if !exists {
			response.Error(c, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		id, ok := userID.(string)
		if !ok {
			response.Error(c, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		// 检查是否是管理员
		if roleChecker != nil {
			isAdmin, err := roleChecker.IsAdmin(id)
			if err != nil || !isAdmin {
				response.Error(c, apperrors.ErrForbidden)
				c.Abort()
				return
			}
		} else {
			// 如果没有设置角色检查器，默认拒绝访问
			response.Error(c, apperrors.ErrForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}
