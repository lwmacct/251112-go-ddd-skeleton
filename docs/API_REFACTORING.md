# API 路由架构总结

## 当前 API 架构

本项目采用 RESTful API 设计，路由按照功能和权限进行清晰分组。

## 路由结构

```
========== 公开端点 ==========
POST /api/auth/register     # 用户注册
POST /api/auth/login        # 用户登录
POST /api/auth/refresh      # 刷新访问令牌

========== 用户个人中心（需要认证）==========
# 个人信息管理
GET /api/user               # 获取当前用户信息
PUT /api/user               # 更新当前用户信息
PATCH /api/user             # 部分更新当前用户
DELETE /api/user            # 注销账号
PUT /api/user/password      # 修改密码
PUT /api/user/email         # 修改邮箱（TODO）
POST /api/user/avatar       # 上传头像（TODO）

# 会话管理
GET /api/user/sessions      # 查看活跃会话
DELETE /api/user/sessions/:id   # 撤销指定会话

# 令牌管理
GET /api/user/tokens        # 获取个人访问令牌列表
POST /api/user/tokens       # 创建个人访问令牌
DELETE /api/user/tokens/:id # 撤销令牌

# 双因素认证
POST /api/user/2fa/enable   # 启用双因素认证
POST /api/user/2fa/verify   # 验证双因素认证
POST /api/user/2fa/disable  # 禁用双因素认证

# 认证相关
POST /api/auth/logout       # 登出

========== 订单管理（需要认证）==========
POST /api/orders            # 创建订单
GET /api/orders             # 列出当前用户订单
GET /api/orders/:id         # 获取订单详情
POST /api/orders/:id/cancel # 取消订单
POST /api/orders/:id/payment    # 处理支付
GET /api/orders/:id/payment     # 获取支付信息
POST /api/orders/:id/shipment   # 创建发货
GET /api/orders/:id/shipment    # 获取发货信息

========== 用户菜单（需要认证）==========
GET /api/menus/user/tree    # 获取当前用户的菜单树（前端侧边栏）

========== 管理员接口（需要认证 + 管理员权限）==========
# 用户管理
GET /api/admin/users        # 列出所有用户
GET /api/admin/users/:id    # 获取指定用户信息
PUT /api/admin/users/:id    # 更新指定用户信息
DELETE /api/admin/users/:id # 删除用户
POST /api/admin/users/:id/ban   # 封禁用户（TODO）
POST /api/admin/users/:id/unban # 解封用户（TODO）

# 订单管理
GET /api/admin/orders       # 列出所有订单
GET /api/admin/orders/:id   # 获取任意订单详情
PUT /api/admin/orders/:id/status    # 更新订单状态

# 菜单管理
POST /api/admin/menus       # 创建菜单
PUT /api/admin/menus/:id    # 更新菜单
DELETE /api/admin/menus/:id # 删除菜单
GET /api/admin/menus/tree   # 获取所有菜单树
PUT /api/admin/menus/order  # 更新菜单排序

# 角色管理
POST /api/admin/roles       # 创建角色
PUT /api/admin/roles/:id    # 更新角色
DELETE /api/admin/roles/:id # 删除角色
GET /api/admin/roles/:id    # 获取角色详情
GET /api/admin/roles        # 列出所有角色

# 角色-菜单关联
POST /api/admin/roles/:roleId/menus    # 为角色分配菜单
GET /api/admin/roles/:roleId/menus     # 获取角色菜单树

# 角色-权限关联
POST /api/admin/roles/:roleId/permissions  # 为角色分配权限
GET /api/admin/roles/:roleId/permissions   # 获取角色权限列表

# 用户-角色管理
POST /api/admin/users/:userId/roles/:roleId    # 为用户分配角色
DELETE /api/admin/users/:userId/roles/:roleId  # 移除用户角色
GET /api/admin/users/:userId/roles             # 获取用户角色列表
```

## 设计原则

### 1. RESTful 最佳实践

#### 单数 vs 复数
- **`/user`** (单数) - 表示当前用户，单一资源，无需 ID
- **`/users`** (复数) - 表示用户集合（仅管理员可访问）

这遵循了 GitHub、GitLab 等大厂的 API 设计模式。

#### 语义清晰
- `/api/user` - 明确表示"当前用户"，从 JWT Token 获取身份
- `/api/admin/users/:id` - 明确表示"管理员操作指定用户"

### 2. 权限分离

```go
// 用户个人中心 - 从 JWT 获取用户 ID
func (h *UserHandler) GetProfile(c *gin.Context) {
    userID := c.GetString("userID")  // 从 Auth 中间件设置的 context 获取
    // 用户只能访问自己的数据
}

// 管理员接口 - 从 URL 参数获取
func (h *UserHandler) GetUser(c *gin.Context) {
    userID := c.Param("id")  // 从 URL 获取
    // 管理员可以访问任意用户数据
}
```

### 3. 中间件链

```go
// 用户个人中心和订单
authenticated := api.Group("")
authenticated.Use(middleware.Auth())  // 仅需认证

// 管理员接口
admin := api.Group("/admin")
admin.Use(middleware.Auth(), middleware.Admin())  // 认证 + 管理员检查
```

## 安全机制

### 1. 用户数据隔离

```go
// ✅ 正确：用户只能操作自己的数据
userID := c.GetString("userID")  // 从认证 token 获取，不可伪造
h.userService.GetUser(ctx, userID)

// ❌ 错误：允许用户通过 URL 操作他人数据
userID := c.Param("id")  // 来自 URL，可被篡改
```

### 2. 管理员权限检查

Admin 中间件会进行两层检查：
1. **认证检查**（Auth 中间件）- 验证 JWT Token 有效性
2. **角色检查**（Admin 中间件）- 通过 RBAC 系统验证用户是否具有 `admin` 角色

```go
// 实现位置：internal/adapters/http/middleware/admin.go
func Admin() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, exists := c.Get("userID")
        if !exists {
            response.Error(c, apperrors.ErrUnauthorized)
            c.Abort()
            return
        }

        // 使用 RoleChecker 检查管理员权限
        isAdmin, err := roleChecker.IsAdmin(userID.(string))
        if err != nil || !isAdmin {
            response.Error(c, apperrors.ErrForbidden)
            c.Abort()
            return
        }

        c.Next()
    }
}
```

### 3. RBAC 权限系统

项目实现了完整的 RBAC（基于角色的访问控制）系统：

- **默认角色**：admin, user, editor, viewer
- **权限管理**：细粒度的操作权限（如 `user:create`, `order:read`）
- **菜单控制**：根据角色动态生成前端菜单树

详细说明请参考：
- `docs/RBAC_IMPLEMENTATION.md` - RBAC 实现指南
- `docs/ADMIN_MIDDLEWARE_SETUP.md` - Admin 中间件配置

## 架构优势

### ✅ 安全性
- 用户无法访问他人数据
- 明确的权限边界
- 基于角色的动态权限控制

### ✅ 语义清晰
- `/user` 明确表示"当前用户"
- `/admin` 明确表示"管理功能"
- 路由结构一目了然

### ✅ 符合标准
- 遵循 RESTful 最佳实践
- 与主流 API（GitHub、GitLab）设计一致

### ✅ 易于扩展
- 清晰的 DDD 分层架构
- 便于添加新的资源和端点

### ✅ 前端友好
- 个人中心 API 不需要知道用户 ID
- 简化前端状态管理
- 动态菜单树支持（根据权限自动过滤）

## DDD 架构

```
Adapters (HTTP)
    ↓ 调用
Application (用例)
    ↓ 调用
Domain (业务逻辑)
    ↑ 定义接口
Infrastructure (技术实现)
```

### 各层职责

- **Adapters 层**：HTTP 请求处理、参数验证、响应格式化
- **Application 层**：用例编排、CQRS 实现、DTO 转换
- **Domain 层**：纯业务逻辑、实体、值对象、领域服务
- **Infrastructure 层**：数据库、缓存、外部服务实现

## 初始化和测试

### 数据库初始化

```bash
# 1. 创建数据库表结构
./main migrate up

# 2. 初始化基础数据（角色、权限、菜单、管理员账户）
./main migrate seed
```

初始化后会自动创建：
- **默认管理员**：admin@example.com / Admin@123456
- **4 个角色**：admin, user, editor, viewer
- **17 个权限**：user:*, role:*, menu:*, order:*
- **7 个菜单**：系统管理（含 4 个子菜单）、订单管理、个人中心

详细说明请参考：`docs/SEED_USAGE.md`

### API 测试示例

#### 1. 用户注册和登录

```bash
# 注册新用户
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "username": "testuser"
  }'

# 登录获取 token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

#### 2. 用户个人中心

```bash
# 获取当前用户信息
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/user

# 更新当前用户信息
curl -X PUT \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"username": "newname"}' \
  http://localhost:8080/api/user

# 获取用户菜单树（前端侧边栏）
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/menus/user/tree
```

#### 3. 管理员接口

```bash
# 登录管理员账户
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "Admin@123456"
  }'

# 列出所有用户（需要管理员令牌）
curl -H "Authorization: Bearer $ADMIN_TOKEN" \
  http://localhost:8080/api/admin/users

# 获取指定用户
curl -H "Authorization: Bearer $ADMIN_TOKEN" \
  http://localhost:8080/api/admin/users/USER_ID

# 获取所有菜单树
curl -H "Authorization: Bearer $ADMIN_TOKEN" \
  http://localhost:8080/api/admin/menus/tree

# 创建新角色
curl -X POST \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "编辑员",
    "code": "editor",
    "description": "内容编辑权限"
  }' \
  http://localhost:8080/api/admin/roles
```

#### 4. 权限测试

```bash
# 普通用户访问管理员接口（应该返回 403 Forbidden）
curl -H "Authorization: Bearer $USER_TOKEN" \
  http://localhost:8080/api/admin/users
```

## 参考资料

- [项目 CLAUDE.md](../CLAUDE.md) - 完整架构文档
- [GitHub REST API](https://docs.github.com/en/rest) - API 设计参考
- [Microsoft API Design Best Practices](https://learn.microsoft.com/en-us/azure/architecture/best-practices/api-design)
- [Domain-Driven Design](https://www.domainlanguage.com/ddd/)

---

**最后更新**: 2024
**版本**: 2.0 - 包含完整 RBAC 系统和 Seed 初始化
