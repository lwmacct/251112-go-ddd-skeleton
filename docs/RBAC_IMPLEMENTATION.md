# RBAC 菜单权限系统实现指南

## 已完成的核心功能

### 1. Domain层（领域层）✅
- `internal/domain/rbac/role.go` - 角色实体
- `internal/domain/rbac/permission.go` - 权限实体
- `internal/domain/rbac/menu.go` - 菜单实体（支持两层结构）
- `internal/domain/rbac/repository.go` - 仓储接口定义
- `internal/domain/rbac/service.go` - 领域服务（权限检查、菜单树构建）
- `internal/domain/rbac/errors.go` - 领域错误定义

### 2. Infrastructure层（基础设施层）✅
- `internal/infrastructure/persistence/model/rbac.go` - GORM模型（6张表）
- `internal/infrastructure/persistence/mapper/rbac.go` - 实体-模型映射器
- `internal/infrastructure/persistence/repository/role_repo.go` - 角色仓储实现
- `internal/infrastructure/persistence/repository/permission_repo.go` - 权限仓储实现
- `internal/infrastructure/persistence/repository/menu_repo.go` - 菜单仓储实现（含用户菜单查询）
- `migrations/000005_create_rbac_tables.up.sql` - 数据库迁移文件

### 3. Application层（应用层）✅
- `internal/application/menu/` - 菜单应用服务（完整 CQRS）
  - `commands.go` - CreateMenu, UpdateMenu, DeleteMenu, AssignMenusToRole等
  - `queries.go` - GetUserMenuTree, GetAllMenuTree, GetRoleMenuTree
  - `dto.go` - 菜单树DTO（支持递归children）
- `internal/application/role/` - 角色应用服务
  - `service.go` - 角色管理服务（CRUD + 关联操作）
  - `dto.go` - 角色相关DTO

### 4. Adapters层（HTTP Handler）✅
- `internal/adapters/http/handler/rbac/menu_handler.go` - 菜单管理Handler
- `internal/adapters/http/handler/rbac/role_handler.go` - 角色管理Handler
- `internal/adapters/http/middleware/admin.go` - 管理员权限中间件
- `internal/adapters/http/middleware/role_checker.go` - RoleChecker 实现
- `internal/adapters/http/router.go` - RBAC 路由配置

### 5. Seed 数据初始化 ✅
- `internal/infrastructure/seed/` - 种子数据系统
  - `seeder.go` - Seeder 接口和 Manager
  - `rbac_seeder.go` - RBAC seed 实现
  - `user_seeder.go` - 默认管理员seed
  - `data/*.yaml` - 默认数据文件（4 角色、17 权限、7 菜单）

---

## 数据库表结构

```sql
-- 6张核心表
1. roles                - 角色表
2. permissions          - 权限表
3. menus                - 菜单表（支持parent_id两层结构）
4. user_roles           - 用户-角色关联（多对多）
5. role_permissions     - 角色-权限关联（多对多）
6. role_menus           - 角色-菜单关联（多对多）
```

---

## 核心使用示例

### 1. 获取用户的菜单树（前端调用）

```go
// 在HTTP Handler中
func (h *MenuHandler) GetUserMenuTree(c *gin.Context) {
    // 从JWT token获取用户ID
    userID := c.GetString("user_id")

    // 调用应用服务
    menuTree, err := h.menuService.GetUserMenuTree(c.Request.Context(), userID)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, menuTree)
}
```

**响应示例：**
```json
{
  "menus": [
    {
      "id": "1",
      "name": "系统管理",
      "path": "/system",
      "icon": "setting",
      "type": "dir",
      "sortOrder": 0,
      "isVisible": true,
      "children": [
        {
          "id": "11",
          "name": "用户管理",
          "path": "/system/users",
          "icon": "user",
          "type": "menu",
          "sortOrder": 0,
          "component": "system/users/index",
          "permission": "user:read"
        },
        {
          "id": "12",
          "name": "角色管理",
          "path": "/system/roles",
          "icon": "team",
          "type": "menu",
          "sortOrder": 1,
          "component": "system/roles/index",
          "permission": "role:read"
        }
      ]
    },
    {
      "id": "2",
      "name": "订单管理",
      "path": "/orders",
      "icon": "order",
      "type": "menu",
      "sortOrder": 1,
      "component": "orders/index",
      "permission": "order:read"
    }
  ]
}
```

---

## 系统初始化和使用

### 1. 数据库迁移和 Seed

```bash
# 1. 创建数据库表结构
./main migrate up

# 2. 初始化 RBAC 基础数据（角色、权限、菜单、管理员）
./main migrate seed
```

执行 seed 后将自动创建：
- **4 个默认角色**：
  - `admin` - 超级管理员（所有权限）
  - `user` - 普通用户（基础业务权限）
  - `editor` - 编辑员（内容编辑权限）
  - `viewer` - 访客（只读权限）

- **17 个默认权限**：
  - 用户管理：`user:create`, `user:read`, `user:update`, `user:delete`
  - 角色管理：`role:create`, `role:read`, `role:update`, `role:delete`
  - 菜单管理：`menu:create`, `menu:read`, `menu:update`, `menu:delete`
  - 权限管理：`permission:read`
  - 订单管理：`order:create`, `order:read`, `order:update`, `order:delete`

- **7 个默认菜单**（两层树形结构）：
  ```
  系统管理/ (目录)
  ├── 用户管理 (菜单)
  ├── 角色管理 (菜单)
  ├── 菜单管理 (菜单)
  └── 权限管理 (菜单)

  订单管理 (菜单)
  个人中心 (菜单)
  ```

- **1 个默认管理员**：
  - Email: `admin@example.com`
  - Password: `Admin@123456`
  - Role: `admin`

⚠️ **重要**：生产环境部署后请立即修改默认密码！

详细 Seed 使用指南请参考：`docs/SEED_USAGE.md`

### 2. API 端点

#### 用户端点（需要认证）
```
GET /api/menus/user/tree    # 获取当前用户的菜单树（前端侧边栏）
```

#### 管理员端点（需要认证 + 管理员权限）

**菜单管理**：
```
POST   /api/admin/menus        # 创建菜单
PUT    /api/admin/menus/:id    # 更新菜单
DELETE /api/admin/menus/:id    # 删除菜单
GET    /api/admin/menus/tree   # 获取所有菜单树
PUT    /api/admin/menus/order  # 更新菜单排序
POST   /api/admin/roles/:roleId/menus    # 为角色分配菜单
GET    /api/admin/roles/:roleId/menus    # 获取角色菜单树
```

**角色管理**：
```
POST   /api/admin/roles        # 创建角色
PUT    /api/admin/roles/:id    # 更新角色
DELETE /api/admin/roles/:id    # 删除角色
GET    /api/admin/roles/:id    # 获取角色详情
GET    /api/admin/roles        # 列出所有角色
POST   /api/admin/roles/:roleId/permissions      # 为角色分配权限
GET    /api/admin/roles/:roleId/permissions      # 获取角色权限
```

**用户-角色管理**：
```
POST   /api/admin/users/:userId/roles/:roleId    # 为用户分配角色
DELETE /api/admin/users/:userId/roles/:roleId    # 移除用户角色
GET    /api/admin/users/:userId/roles            # 获取用户角色
```

### 3. 管理员权限检查

Admin 中间件已完全配置：

```go
// 实现位置：internal/adapters/http/middleware/admin.go
// RoleChecker 位置：internal/adapters/http/middleware/role_checker.go

// Admin 中间件会检查：
// 1. 用户是否已认证（Auth 中间件）
// 2. 用户是否具有 admin 角色（通过 RBAC 系统）

admin := api.Group("/admin")
admin.Use(middleware.Auth(), middleware.Admin())
```

RoleChecker 在 `internal/bootstrap/container.go` 中配置：
```go
middleware.SetRoleChecker(middleware.NewRBACRoleChecker(rbacDomainService))
```

详细配置说明请参考：`docs/ADMIN_MIDDLEWARE_SETUP.md`

### 4. 测试流程

#### 步骤1：使用管理员账户登录

```bash
# 登录管理员账户
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "Admin@123456"
  }'

# 响应示例
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_at": "2024-01-01T12:00:00Z"
}
```

#### 步骤2：获取用户菜单树

```bash
curl -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  http://localhost:8080/api/menus/user/tree
```

**响应示例**：
```json
{
  "menus": [
    {
      "id": "01H...",
      "name": "系统管理",
      "path": "/system",
      "icon": "setting",
      "type": "dir",
      "sortOrder": 0,
      "isVisible": true,
      "children": [
        {
          "id": "01H...",
          "name": "用户管理",
          "path": "/system/users",
          "icon": "user",
          "type": "menu",
          "sortOrder": 0,
          "component": "system/users/index",
          "permission": "user:read"
        },
        {
          "id": "01H...",
          "name": "角色管理",
          "path": "/system/roles",
          "icon": "team",
          "type": "menu",
          "sortOrder": 1,
          "component": "system/roles/index",
          "permission": "role:read"
        }
      ]
    }
  ]
}
```

#### 步骤3：测试管理员接口

```bash
# 获取所有菜单树（管理员）
curl -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  http://localhost:8080/api/admin/menus/tree

# 创建新角色
curl -X POST \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "产品经理",
    "code": "product_manager",
    "description": "产品管理相关权限"
  }' \
  http://localhost:8080/api/admin/roles

# 为用户分配角色
curl -X POST \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  http://localhost:8080/api/admin/users/USER_ID/roles/ROLE_ID
```

#### 步骤4：验证权限检查

```bash
# 普通用户访问管理员接口（应该返回 403 Forbidden）
# 1. 先注册一个普通用户
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "username": "testuser"
  }'

# 2. 登录普通用户
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

# 3. 尝试访问管理员接口（应该被拒绝）
curl -H "Authorization: Bearer USER_TOKEN" \
  http://localhost:8080/api/admin/users

# 预期响应：403 Forbidden
{
  "error": "Forbidden"
}
```

### 5. SQL 查询验证

```sql
-- 查看所有角色
SELECT * FROM roles;

-- 查看所有权限
SELECT * FROM permissions;

-- 查看菜单树结构
SELECT
    m1.name AS parent_menu,
    m2.name AS child_menu,
    m2.component,
    m2.permission
FROM menus m1
LEFT JOIN menus m2 ON m2.parent_id = m1.id
WHERE m1.parent_id IS NULL
ORDER BY m1.sort_order, m2.sort_order;

-- 查看用户的角色
SELECT
    u.email,
    u.username,
    r.name AS role_name,
    r.code AS role_code
FROM users u
JOIN user_roles ur ON ur.user_id = u.id
JOIN roles r ON r.id = ur.role_id
WHERE u.email = 'admin@example.com';

-- 查看角色的权限
SELECT
    r.name AS role_name,
    p.name AS permission_name,
    p.code AS permission_code
FROM roles r
JOIN role_permissions rp ON rp.role_id = r.id
JOIN permissions p ON p.id = rp.permission_id
WHERE r.code = 'admin';

-- 查看 seed 执行历史
SELECT * FROM seed_history ORDER BY executed_at DESC;
```

---

## 前端集成示例（React/Vue）

### 获取菜单并渲染

```typescript
// API调用
async function getUserMenus() {
  const response = await fetch('/api/v1/rbac/menus/user/tree', {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  const data = await response.json();
  return data.menus;
}

// 菜单渲染（递归组件）
function MenuItem({ menu }) {
  return (
    <div>
      <a href={menu.path}>
        <Icon type={menu.icon} />
        {menu.name}
      </a>
      {menu.children && menu.children.map(child => (
        <MenuItem key={child.id} menu={child} />
      ))}
    </div>
  );
}
```

---

## 测试流程

1. **运行数据库迁移**
```bash
make migrate-up  # 创建RBAC表
```

2. **插入种子数据**
```bash
make migrate-seed  # 插入默认角色和菜单
```

3. **分配角色给测试用户**
```sql
INSERT INTO user_roles (user_id, role_id)
VALUES ('your-user-id', '01H0EXAMPLE0000ADMIN');
```

4. **测试API**
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/v1/rbac/menus/user/tree
```

---

## 架构特点

1. **DDD分层清晰**：Domain → Application → Infrastructure → Adapters
2. **CQRS模式**：Commands 和 Queries 独立文件实现
3. **依赖倒置**：Domain定义接口，Infrastructure实现
4. **两层菜单树**：优化查询性能，满足大多数场景
5. **权限与菜单解耦**：菜单可选关联权限码
6. **扩展性强**：可轻松添加更多权限和菜单
7. **Seed 系统**：支持幂等性的数据库初始化
8. **完整 RBAC**：角色、权限、菜单完整实现

---

## 扩展建议

### 1. 细粒度权限控制

可以添加基于权限的中间件：

```go
// middleware/permission.go
func RequirePermission(rbacService *rbac.Service, permissionCode string) gin.HandlerFunc {
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
    middleware.RequirePermission(rbacService, "user:delete"),
    userHandler.DeleteUser,
)
```

### 2. Redis 缓存优化

```go
// 缓存用户菜单树（TTL: 5分钟）
func (s *Service) GetUserMenuTree(ctx context.Context, userID string) (*MenuTreeResponse, error) {
    cacheKey := fmt.Sprintf("user:menus:%s", userID)
    if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
        return cached, nil
    }

    menuTree, err := s.domainService.GetUserMenuTree(ctx, userID)
    if err != nil {
        return nil, err
    }

    result := s.buildMenuTreeResponse(menuTree)
    s.cache.Set(ctx, cacheKey, result, 5*time.Minute)

    return result, nil
}
```

### 3. 审计日志

```go
// 记录敏感操作
type AuditLog struct {
    UserID    string
    Action    string
    Resource  string
    Detail    string
    IP        string
    CreatedAt time.Time
}

// 在关键操作后记录
auditLog.Record(ctx, AuditLog{
    UserID:   adminID,
    Action:   "AssignRole",
    Resource: fmt.Sprintf("user:%s,role:%s", userID, roleID),
    IP:       c.ClientIP(),
})
```

### 4. 三层或多层菜单支持

当前实现支持两层菜单，如需要更多层级，可以修改领域服务的 `ValidateMenuHierarchy` 方法。

### 5. 动态权限

可以在菜单上添加动态权限表达式：

```json
{
  "name": "编辑文章",
  "permission": "article:edit",
  "dynamicPermission": "article.author_id == user.id OR user.role == 'admin'"
}
```

---

## 相关文档

- [`docs/SEED_USAGE.md`](./SEED_USAGE.md) - Seed 数据使用指南
- [`docs/ADMIN_MIDDLEWARE_SETUP.md`](./ADMIN_MIDDLEWARE_SETUP.md) - Admin 中间件配置
- [`docs/RBAC_INTEGRATION.md`](./RBAC_INTEGRATION.md) - RBAC 集成指南
- [`docs/API_REFACTORING.md`](./API_REFACTORING.md) - API 路由架构总结
- [`CLAUDE.md`](../CLAUDE.md) - 完整项目架构文档

---

## 关键文件位置

```
internal/domain/rbac/                  # 领域层（纯业务逻辑）
├── role.go
├── permission.go
├── menu.go
├── repository.go
├── service.go                         # 核心：CheckPermission, GetUserMenuTree
└── errors.go

internal/infrastructure/persistence/   # 基础设施层
├── model/rbac.go                      # GORM模型
├── mapper/rbac.go                     # 映射器
└── repository/
    ├── role_repo.go
    ├── permission_repo.go
    └── menu_repo.go                   # 核心：GetUserMenus

internal/application/menu/             # 应用层
├── queries.go                         # 核心：GetUserMenuTree API
└── dto.go                             # MenuTreeResponse

migrations/
└── 000005_create_rbac_tables.up.sql  # 数据库迁移
```

---

## 🎉 系统完成状态

### ✅ 已完成的功能

你现在拥有了一个完整的、遵循DDD原则的RBAC菜单权限系统：

**核心功能**：
- ✅ 3个领域实体（Role, Permission, Menu）
- ✅ 完整的仓储模式实现
- ✅ 领域服务（权限检查、菜单树构建）
- ✅ 7张数据库表（6张RBAC + 1张seed_history）
- ✅ 用户菜单树查询API（核心功能）
- ✅ 支持两层菜单结构
- ✅ 多对多关系管理（用户-角色-权限-菜单）

**应用层**：
- ✅ 菜单应用服务（Commands + Queries + DTO）
- ✅ 角色应用服务（完整 CRUD + 关联操作）

**HTTP层**：
- ✅ 菜单管理Handler（8个端点）
- ✅ 角色管理Handler（11个端点）
- ✅ Admin 中间件 + RoleChecker 实现
- ✅ 完整路由配置

**Seed 系统**：
- ✅ Seeder 接口和 Manager
- ✅ RBAC Seeder（4角色 + 17权限 + 7菜单）
- ✅ User Seeder（默认管理员）
- ✅ 支持幂等性和事务保证
- ✅ 执行历史记录

**使用方法**：
```bash
# 1. 创建表结构
./main migrate up

# 2. 初始化数据
./main migrate seed

# 3. 使用默认管理员登录
Email: admin@example.com
Password: Admin@123456

# 4. 调用 API 获取菜单树
GET /api/menus/user/tree
```

核心API `GetUserMenuTree(userID)` 可以直接返回用户根据角色能看到的菜单树，前端可以直接渲染！

---

**最后更新**：2024
**版本**：2.0 - 包含完整 RBAC 系统、Seed 初始化和 Admin 中间件
