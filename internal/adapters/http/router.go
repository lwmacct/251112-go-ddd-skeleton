package http

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-ddd-skeleton/internal/adapters/http/handler"
	"github.com/yourusername/go-ddd-skeleton/internal/adapters/http/middleware"
)

// SetupRouter 配置路由
func SetupRouter(
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
	orderHandler *handler.OrderHandler,
) *gin.Engine {
	r := gin.New()

	// 全局中间件
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// 公开端点
		v1.POST("/auth/register", userHandler.Register)
		v1.POST("/auth/login", authHandler.Login)

		// 需要认证的端点
		authenticated := v1.Group("")
		authenticated.Use(middleware.Auth())
		{
			// 用户相关
			users := authenticated.Group("/users")
			{
				users.GET("/:id", userHandler.GetUser)
				users.PUT("/:id", userHandler.UpdateUser)
				users.DELETE("/:id", userHandler.DeleteUser)
				users.GET("", userHandler.ListUsers)
			}

			// 认证相关
			auth := authenticated.Group("/auth")
			{
				auth.POST("/logout", authHandler.Logout)
				auth.POST("/refresh", authHandler.RefreshToken)
				auth.POST("/2fa/enable", authHandler.Enable2FA)
				auth.POST("/2fa/verify", authHandler.Verify2FA)
				auth.POST("/pat", authHandler.CreatePAT)
				auth.GET("/sessions", authHandler.GetSessions)
			}

			// 订单相关
			orders := authenticated.Group("/orders")
			{
				orders.POST("", orderHandler.CreateOrder)
				orders.GET("/:id", orderHandler.GetOrder)
				orders.GET("", orderHandler.ListOrders)
				orders.POST("/:id/cancel", orderHandler.CancelOrder)
				orders.POST("/:id/payment", orderHandler.ProcessPayment)
				orders.POST("/:id/shipment", orderHandler.CreateShipment)
			}
		}
	}

	return r
}

