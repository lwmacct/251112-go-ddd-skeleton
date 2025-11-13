# Admin 中间件 RoleChecker 实现文档

## 概述

本文档描述了 Admin 中间件的 RoleChecker 实现，用于保护管理员端点，确保只有具有 `admin` 角色的用户才能访问 `/api/admin/*` 路由。

## 架构设计

### 依赖倒置原则

```
Middleware (适配器层)
    ↓ 依赖接口
RoleChecker 接口
    ↑ 实现接口
rbacRoleChecker (适配器)
    ↓ 使用
RBAC Domain Service (领域层)
    ↓ 使用
RoleRepository (基础设施层)
```

### 组件说明

1. **RoleChecker 接口** (`internal/adapters/http/middleware/admin.go`)
   - 定义角色检查的抽象接口
   - 与领域层解耦

2. **rbacRoleChecker 实现** (`internal/adapters/http/middleware/role_checker.go`)
   - 实现 RoleChecker 接口
   - 桥接中间件和 RBAC 领域服务
   - 将 userID 转换为角色检查逻辑

3. **RBAC Domain Service** (`internal/domain/rbac/service.go`)
   - 提供 `CheckUserHasRole()` 方法
   - 纯粹的领域逻辑，无框架依赖

## 实现细节

### 1. RoleChecker 接口定义

```go
// internal/adapters/http/middleware/admin.go
type RoleChecker interface {
    IsAdmin(userID string) (bool, error)
}
```

### 2. RBAC RoleChecker 实现

```go
// internal/adapters/http/middleware/role_checker.go
type rbacRoleChecker struct {
    rbacService *rbac.Service
}

func NewRBACRoleChecker(rbacService *rbac.Service) RoleChecker {
    return &rbacRoleChecker{
        rbacService: rbacService,
    }
}

func (r *rbacRoleChecker) IsAdmin(userID string) (bool, error) {
    return r.rbacService.CheckUserHasRole(context.Background(), userID, "admin")
}
```

### 3. 依赖注入配置

```go
// internal/bootstrap/container.go
func NewContainer(cfg *config.Config) (*Container, error) {
    // ... 初始化 rbacDomainService ...

    // 设置中间件依赖
    middleware.SetTokenValidator(jwtIssuer)
    middleware.SetRoleChecker(middleware.NewRBACRoleChecker(rbacDomainService))

    // ...
}
```

### 4. Admin 中间件逻辑

```go
// internal/adapters/http/middleware/admin.go
func Admin() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从 context 获取用户 ID（由 Auth 中间件设置）
        userID, exists := c.Get("userID")
        if !exists {
            response.Error(c, apperrors.ErrUnauthorized)
            c.Abort()
            return
        }

        // 2. 检查是否是管理员
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
```

## 路由配置

所有管理员路由都使用了 `Auth()` 和 `Admin()` 中间件链：

```go
// internal/adapters/http/router.go
admin := api.Group("/admin")
admin.Use(middleware.Auth(), middleware.Admin())
{
    // 用户管理
    adminUsers := admin.Group("/users")
    {
        adminUsers.GET("", userHandler.ListUsers)
        adminUsers.DELETE("/:id", userHandler.DeleteUser)
        // ...
    }

    // 订单管理
    adminOrders := admin.Group("/orders")
    // ...

    // RBAC 管理
    adminMenus := admin.Group("/menus")
    adminRoles := admin.Group("/roles")
    // ...
}
```

## 验证流程

### 正常流程

1. 用户发送请求到 `/api/admin/users`
2. `Auth()` 中间件验证 JWT Token，提取 userID 并存入 context
3. `Admin()` 中间件调用 `roleChecker.IsAdmin(userID)`
4. `rbacRoleChecker` 调用 `rbacService.CheckUserHasRole(userID, "admin")`
5. RBAC 服务查询数据库，检查用户是否有 `admin` 角色
6. 如果有，继续执行 handler；如果没有，返回 403 Forbidden

### 错误流程

| 场景 | HTTP 状态码 | 错误信息 |
|------|------------|---------|
| 未提供 JWT Token | 401 | Unauthorized |
| JWT Token 无效 | 401 | Unauthorized |
| 用户没有 admin 角色 | 403 | Forbidden |
| RoleChecker 未配置 | 403 | Forbidden |
| 数据库查询失败 | 403 | Forbidden |

## 测试方法

### 前置条件

1. 确保数据库中有 `roles` 表和 `user_roles` 关联表
2. 创建一个角色记录：`code = "admin"`
3. 为测试用户分配 admin 角色

### 测试场景

#### 1. 测试普通用户访问管理员端点（应该被拒绝）

```bash
# 登录普通用户
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password"
  }'

# 使用返回的 token 访问管理员端点
curl -X GET http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer <user_token>"

# 预期结果: 403 Forbidden
```

#### 2. 测试管理员用户访问管理员端点（应该成功）

```bash
# 登录管理员用户
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password"
  }'

# 使用返回的 token 访问管理员端点
curl -X GET http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer <admin_token>"

# 预期结果: 200 OK，返回用户列表
```

#### 3. 测试未认证用户访问管理员端点（应该被拒绝）

```bash
curl -X GET http://localhost:8080/api/admin/users

# 预期结果: 401 Unauthorized
```

### 使用 SQL 为用户分配管理员角色

```sql
-- 1. 查找或创建 admin 角色
INSERT INTO roles (id, name, code, description, is_active, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'Administrator',
    'admin',
    'System administrator with full access',
    true,
    NOW(),
    NOW()
)
ON CONFLICT (code) DO NOTHING;

-- 2. 为用户分配 admin 角色
INSERT INTO user_roles (user_id, role_id)
VALUES (
    '<user_id>',  -- 替换为实际的用户 ID
    '00000000-0000-0000-0000-000000000001'
)
ON CONFLICT DO NOTHING;

-- 3. 验证分配结果
SELECT u.email, r.code, r.name
FROM users u
JOIN user_roles ur ON u.id = ur.user_id
JOIN roles r ON ur.role_id = r.id
WHERE u.id = '<user_id>';
```

## 扩展建议

### 1. 添加更细粒度的权限检查中间件

```go
// middleware/permission.go
func RequirePermission(permissionCode string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetString("userID")
        hasPermission, err := rbacService.CheckPermission(c.Request.Context(), userID, permissionCode)
        if err != nil || !hasPermission {
            response.Error(c, apperrors.ErrForbidden)
            c.Abort()
            return
        }
        c.Next()
    }
}

// 使用示例
adminUsers.DELETE("/:id",
    middleware.RequirePermission("user:delete"),
    userHandler.DeleteUser,
)
```

### 2. 添加审计日志

```go
func (r *rbacRoleChecker) IsAdmin(userID string) (bool, error) {
    isAdmin, err := r.rbacService.CheckUserHasRole(context.Background(), userID, "admin")

    // 记录管理员访问日志
    if isAdmin && err == nil {
        log.Info("admin access",
            zap.String("user_id", userID),
            zap.Time("timestamp", time.Now()),
        )
    }

    return isAdmin, err
}
```

### 3. 添加 Redis 缓存

```go
func (r *rbacRoleChecker) IsAdmin(userID string) (bool, error) {
    // 尝试从缓存获取
    cacheKey := fmt.Sprintf("user:role:admin:%s", userID)
    if cached, err := r.cache.Get(cacheKey); err == nil {
        return cached.(bool), nil
    }

    // 从数据库查询
    isAdmin, err := r.rbacService.CheckUserHasRole(context.Background(), userID, "admin")
    if err != nil {
        return false, err
    }

    // 写入缓存（5分钟过期）
    r.cache.Set(cacheKey, isAdmin, 5*time.Minute)

    return isAdmin, nil
}
```

## 安全考虑

1. **默认拒绝策略**：如果 RoleChecker 未配置或查询失败，默认拒绝访问
2. **中间件顺序**：Admin 中间件必须在 Auth 中间件之后使用
3. **角色编码规范**：使用统一的角色代码（如 "admin"），避免硬编码
4. **错误处理**：所有错误都返回 403 Forbidden，避免信息泄露
5. **审计日志**：建议记录所有管理员操作，便于安全审计

## 总结

通过实现 RoleChecker 适配器，我们成功地：

1. ✅ 保持了 DDD 分层架构的清晰性
2. ✅ 遵循了依赖倒置原则（中间件依赖接口，不依赖具体实现）
3. ✅ 实现了基于 RBAC 的管理员权限控制
4. ✅ 为所有 `/api/admin/*` 端点提供了安全保护
5. ✅ 保持了代码的可测试性和可扩展性

后续可以基于此实现更细粒度的权限控制，如基于权限码的中间件 `RequirePermission()`。
