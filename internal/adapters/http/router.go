package http

import (
	"github.com/gin-gonic/gin"
	authhandler "github.com/lwmacct/251112-go-ddd-skeleton/internal/adapters/http/handler/auth"
	orderhandler "github.com/lwmacct/251112-go-ddd-skeleton/internal/adapters/http/handler/order"
	rbachandler "github.com/lwmacct/251112-go-ddd-skeleton/internal/adapters/http/handler/rbac"
	userhandler "github.com/lwmacct/251112-go-ddd-skeleton/internal/adapters/http/handler/user"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/adapters/http/middleware"
)

// SetupRouter 配置路由
func SetupRouter(
	userHandler *userhandler.Handler,
	authHandler *authhandler.Handler,
	orderHandler *orderhandler.Handler,
	menuHandler *rbachandler.MenuHandler,
	roleHandler *rbachandler.RoleHandler,
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

	// API 路由
	api := r.Group("/api")
	{
		// ========== 公开端点 ==========
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// ========== 需要认证的端点 ==========
		authenticated := api.Group("")
		authenticated.Use(middleware.Auth())
		{
			// 用户个人中心
			user := authenticated.Group("/user")
			{
				user.GET("", userHandler.GetProfile)       // GET /api/v1/user
				user.PUT("", userHandler.UpdateProfile)    // PUT /api/v1/user
				user.PATCH("", userHandler.PatchProfile)   // PATCH /api/v1/user
				user.DELETE("", userHandler.DeleteAccount) // DELETE /api/v1/user
				user.PUT("/password", userHandler.ChangePassword)
				user.PUT("/email", userHandler.ChangeEmail)
				user.POST("/avatar", userHandler.UploadAvatar)

				// 会话管理
				user.GET("/sessions", authHandler.GetSessions)
				user.DELETE("/sessions/:id", authHandler.RevokeSession)

				// 令牌管理
				user.GET("/tokens", authHandler.GetPATs)
				user.POST("/tokens", authHandler.CreatePAT)
				user.DELETE("/tokens/:id", authHandler.RevokePAT)

				// 双因素认证
				user.POST("/2fa/enable", authHandler.Enable2FA)
				user.POST("/2fa/verify", authHandler.Verify2FA)
				user.POST("/2fa/disable", authHandler.Disable2FA)
			}

			// 认证相关（非个人中心）
			authGroup := authenticated.Group("/auth")
			{
				authGroup.POST("/logout", authHandler.Logout)
			}

			// 用户订单
			orders := authenticated.Group("/orders")
			{
				orders.POST("", orderHandler.CreateOrder)
				orders.GET("", orderHandler.ListOrders)
				orders.GET("/:id", orderHandler.GetOrder)
				orders.POST("/:id/cancel", orderHandler.CancelOrder)
				orders.POST("/:id/payment", orderHandler.ProcessPayment)
				orders.GET("/:id/payment", orderHandler.GetPayment)
				orders.POST("/:id/shipment", orderHandler.CreateShipment)
				orders.GET("/:id/shipment", orderHandler.GetShipment)
			}

			// 用户菜单（RBAC）
			authenticated.GET("/menus/user/tree", menuHandler.GetUserMenuTree)
		}

		// ========== 管理员接口 ==========
		admin := api.Group("/admin")
		admin.Use(middleware.Auth(), middleware.Admin())
		{

			// 用户-角色管理
			admin.POST("/users/:userId/roles/:roleId", roleHandler.AssignRoleToUser)
			admin.DELETE("/users/:userId/roles/:roleId", roleHandler.RemoveRoleFromUser)
			admin.GET("/users/:userId/roles", roleHandler.GetUserRoles)

			// 用户管理
			adminUsers := admin.Group("/users")
			{
				adminUsers.GET("", userHandler.ListUsers)
				adminUsers.GET("/:id", userHandler.GetUser)
				adminUsers.PUT("/:id", userHandler.UpdateUser)
				adminUsers.DELETE("/:id", userHandler.DeleteUser)
				adminUsers.POST("/:id/ban", userHandler.BanUser)
				adminUsers.POST("/:id/unban", userHandler.UnbanUser)
			}

			// 订单管理
			adminOrders := admin.Group("/orders")
			{
				adminOrders.GET("", orderHandler.ListAllOrders)
				adminOrders.GET("/:id", orderHandler.GetOrderByID)
				adminOrders.PUT("/:id/status", orderHandler.UpdateOrderStatus)
			}

			// RBAC管理
			// 菜单管理
			adminMenus := admin.Group("/menus")
			{
				adminMenus.POST("", menuHandler.CreateMenu)
				adminMenus.PUT("/:id", menuHandler.UpdateMenu)
				adminMenus.DELETE("/:id", menuHandler.DeleteMenu)
				adminMenus.GET("/tree", menuHandler.GetAllMenuTree)
				adminMenus.PUT("/order", menuHandler.UpdateMenuOrder)
			}

			// 角色管理
			adminRoles := admin.Group("/roles")
			{
				adminRoles.POST("", roleHandler.CreateRole)
				adminRoles.PUT("/:id", roleHandler.UpdateRole)
				adminRoles.DELETE("/:id", roleHandler.DeleteRole)
				adminRoles.GET("/:id", roleHandler.GetRole)
				adminRoles.GET("", roleHandler.ListRoles)

				// 角色-菜单关联
				adminRoles.POST("/:roleId/menus", menuHandler.AssignMenusToRole)
				adminRoles.GET("/:roleId/menus", menuHandler.GetRoleMenuTree)

				// 角色-权限关联
				adminRoles.POST("/:roleId/permissions", roleHandler.AssignPermissionsToRole)
				adminRoles.GET("/:roleId/permissions", roleHandler.GetRolePermissions)
			}

		}
	}

	return r
}
