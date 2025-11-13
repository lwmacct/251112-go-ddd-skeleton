# RBAC 系统集成指南

## 🎉 系统状态

RBAC 系统已完全集成！所有层次的实现都已完成，您只需要执行数据库初始化即可开始使用。

## ✅ 已完成的集成工作

### 1. Domain 层 (领域层) ✅
```
internal/domain/rbac/
├── role.go              # 角色实体 + 业务规则
├── permission.go        # 权限实体
├── menu.go              # 菜单实体（支持两层树形结构）
├── repository.go        # 仓储接口定义
├── service.go           # 领域服务（权限检查、菜单树构建）
└── errors.go            # 领域错误定义
```

**核心领域服务**：
- `CheckPermission(userID, permissionCode)` - 检查用户是否具有指定权限
- `CheckUserHasRole(userID, roleCode)` - 检查用户是否具有指定角色
- `GetUserMenuTree(userID)` - 获取用户的菜单树
- `ValidateMenuHierarchy(parentID)` - 验证菜单层级（两层限制）

### 2. Infrastructure 层 (基础设施层) ✅
```
internal/infrastructure/
├── persistence/
│   ├── model/           # GORM 模型（6张RBAC表）
│   ├── mapper/rbac.go   # Domain ↔ Model 映射
│   └── repository/      # Repository 实现
│       ├── role_repo.go
│       ├── permission_repo.go
│       └── menu_repo.go
└── seed/                # 种子数据系统
    ├── seeder.go
    ├── rbac_seeder.go   # RBAC 数据初始化
    ├── user_seeder.go   # 默认管理员
    └── data/            # YAML 数据文件
        ├── roles.yaml
        ├── permissions.yaml
        ├── menus.yaml
        ├── role_permissions.yaml
        ├── role_menus.yaml
        └── users.yaml
```

**数据库表**（已完成迁移）：
- `roles` - 角色表
- `permissions` - 权限表
- `menus` - 菜单表（支持 parent_id 两层结构）
- `user_roles` - 用户-角色关联（多对多）
- `role_permissions` - 角色-权限关联（多对多）
- `role_menus` - 角色-菜单关联（多对多）
- `seed_history` - Seed 执行历史记录

### 3. Application 层 (应用层) ✅
```
internal/application/
├── menu/
│   ├── commands.go      # CreateMenu, UpdateMenu, DeleteMenu等
│   ├── queries.go       # GetUserMenuTree, GetAllMenuTree等
│   └── dto.go           # MenuTreeResponse（支持递归children）
└── role/
    ├── service.go       # 角色管理服务（CRUD + 关联操作）
    └── dto.go           # RoleDTO, UserRoleDTO等
```

**核心应用服务方法**：
- **菜单服务**：GetUserMenuTree, GetAllMenuTree, CreateMenu, UpdateMenu, DeleteMenu
- **角色服务**：CreateRole, UpdateRole, DeleteRole, AssignRoleToUser, AssignPermissionsToRole

### 4. Adapters 层 (HTTP Handler) ✅
```
internal/adapters/http/
├── handler/rbac/
│   ├── menu_handler.go  # 菜单管理（8个端点）
│   └── role_handler.go  # 角色管理（11个端点）
├── middleware/
│   ├── admin.go         # Admin 权限中间件
│   └── role_checker.go  # RoleChecker 实现
└── router.go            # 路由配置（已完成）
```

**HTTP 端点**（已配置）：

**用户端点**（需要认证）：
- `GET /api/menus/user/tree` - 获取当前用户菜单树 ⭐

**管理员端点**（需要认证 + admin角色）：
- **菜单管理**：
  - `POST /api/admin/menus` - 创建菜单
  - `PUT /api/admin/menus/:id` - 更新菜单
  - `DELETE /api/admin/menus/:id` - 删除菜单
  - `GET /api/admin/menus/tree` - 获取所有菜单树
  - `PUT /api/admin/menus/order` - 更新菜单排序
  - `POST /api/admin/roles/:roleId/menus` - 为角色分配菜单
  - `GET /api/admin/roles/:roleId/menus` - 获取角色菜单树

- **角色管理**：
  - `POST /api/admin/roles` - 创建角色
  - `PUT /api/admin/roles/:id` - 更新角色
  - `DELETE /api/admin/roles/:id` - 删除角色
  - `GET /api/admin/roles/:id` - 获取角色详情
  - `GET /api/admin/roles` - 列出所有角色
  - `POST /api/admin/roles/:roleId/permissions` - 为角色分配权限
  - `GET /api/admin/roles/:roleId/permissions` - 获取角色权限

- **用户-角色管理**：
  - `POST /api/admin/users/:userId/roles/:roleId` - 为用户分配角色
  - `DELETE /api/admin/users/:userId/roles/:roleId` - 移除用户角色
  - `GET /api/admin/users/:userId/roles` - 获取用户角色列表

### 5. 依赖注入 ✅

`internal/bootstrap/container.go` 已完成所有 RBAC 服务的依赖注入：

```go
// RBAC仓储
roleRepo := repository.NewRoleRepo(db)
permissionRepo := repository.NewPermissionRepo(db)
menuRepo := repository.NewMenuRepo(db)

// RBAC领域服务
rbacDomainService := rbac.NewService(roleRepo, permissionRepo, menuRepo)

// RBAC应用服务
menuService := appmenu.NewService(rbacDomainService, menuRepo)
roleService := approle.NewService(roleRepo, permissionRepo, rbacDomainService)

// HTTP Handler
menuHandler := rbachandler.NewMenuHandler(menuService)
roleHandler := rbachandler.NewRoleHandler(roleService)

// 配置 Admin 中间件
middleware.SetRoleChecker(middleware.NewRBACRoleChecker(rbacDomainService))

// 路由
router := http.SetupRouter(userHandler, authHandler, orderHandler, menuHandler, roleHandler)
```

### 6. Admin 权限中间件 ✅

完全实现并配置：
- `internal/adapters/http/middleware/admin.go` - Admin 中间件
- `internal/adapters/http/middleware/role_checker.go` - RoleChecker 实现

**工作原理**：
1. Auth 中间件验证 JWT Token，提取 userID
2. Admin 中间件调用 RoleChecker.IsAdmin(userID)
3. RoleChecker 通过 RBAC 领域服务查询数据库
4. 验证用户是否具有 `admin` 角色

详细配置说明：[`docs/ADMIN_MIDDLEWARE_SETUP.md`](./ADMIN_MIDDLEWARE_SETUP.md)

---

## 🚀 快速开始

### 步骤 1：数据库初始化

```bash
# 1. 创建数据库表结构
./main migrate up

# 2. 初始化 RBAC 基础数据
./main migrate seed
```

### 步骤 2：查看初始化的数据

**默认角色**（4个）：
- `admin` - 超级管理员（所有权限）
- `user` - 普通用户（基础业务权限）
- `editor` - 编辑员（内容编辑权限）
- `viewer` - 访客（只读权限）

**默认权限**（17个）：
- 用户管理：`user:create`, `user:read`, `user:update`, `user:delete`
- 角色管理：`role:create`, `role:read`, `role:update`, `role:delete`
- 菜单管理：`menu:create`, `menu:read`, `menu:update`, `menu:delete`
- 权限管理：`permission:read`
- 订单管理：`order:create`, `order:read`, `order:update`, `order:delete`

**默认菜单**（7个，两层树形结构）：
```
系统管理/ (目录)
├── 用户管理 (菜单)
├── 角色管理 (菜单)
├── 菜单管理 (菜单)
└── 权限管理 (菜单)

订单管理 (菜单)
个人中心 (菜单)
```

**默认管理员账户**：
- Email: `admin@example.com`
- Password: `Admin@123456`
- Role: `admin`

⚠️ **重要**：生产环境部署后请立即修改默认密码！

详细 Seed 使用指南：[`docs/SEED_USAGE.md`](./SEED_USAGE.md)

### 步骤 3：测试 API

#### 3.1 登录管理员账户

```bash
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

#### 3.2 获取用户菜单树（核心功能）

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
        }
      ]
    }
  ]
}
```

#### 3.3 测试管理员接口

```bash
# 获取所有菜单树
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

#### 3.4 验证权限检查

```bash
# 注册普通用户
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "username": "testuser"
  }'

# 登录普通用户
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

# 尝试访问管理员接口（应该被拒绝）
curl -H "Authorization: Bearer USER_TOKEN" \
  http://localhost:8080/api/admin/users

# 预期响应：403 Forbidden
{
  "error": "Forbidden"
}
```

---

## 📊 SQL 验证查询

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
    r.name AS role_name,
    r.code AS role_code
FROM users u
JOIN user_roles ur ON ur.user_id = u.id
JOIN roles r ON r.id = ur.role_id
WHERE u.email = 'admin@example.com';

-- 查看角色的权限
SELECT
    r.name AS role_name,
    p.code AS permission_code,
    p.name AS permission_name
FROM roles r
JOIN role_permissions rp ON rp.role_id = r.id
JOIN permissions p ON p.id = rp.permission_id
WHERE r.code = 'admin'
ORDER BY p.code;

-- 查看角色的菜单
SELECT
    r.name AS role_name,
    m.name AS menu_name,
    m.path AS menu_path
FROM roles r
JOIN role_menus rm ON rm.role_id = r.id
JOIN menus m ON m.id = rm.menu_id
WHERE r.code = 'admin'
ORDER BY m.sort_order;

-- 查看 seed 执行历史
SELECT * FROM seed_history ORDER BY executed_at DESC;
```

---

## 🚀 前端集成示例

### React/TypeScript 示例

```typescript
// types.ts
interface MenuTreeItem {
  id: string;
  name: string;
  path: string;
  icon: string;
  type: 'dir' | 'menu' | 'link';
  sortOrder: number;
  isVisible: boolean;
  component?: string;
  permission?: string;
  children?: MenuTreeItem[];
}

interface MenuTreeResponse {
  menus: MenuTreeItem[];
}

// api.ts
async function getUserMenus(): Promise<MenuTreeResponse> {
  const response = await fetch('/api/menus/user/tree', {
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('token')}`
    }
  });

  if (!response.ok) {
    throw new Error('Failed to fetch menus');
  }

  return response.json();
}

// MenuItem.tsx - 递归菜单组件
import { Link } from 'react-router-dom';

interface MenuItemProps {
  menu: MenuTreeItem;
}

function MenuItem({ menu }: MenuItemProps) {
  // 目录类型（dir）不可点击，仅显示子菜单
  if (menu.type === 'dir' && menu.children) {
    return (
      <div className="menu-group">
        <div className="menu-group-title">
          <Icon type={menu.icon} />
          <span>{menu.name}</span>
        </div>
        <div className="submenu">
          {menu.children.map(child => (
            <MenuItem key={child.id} menu={child} />
          ))}
        </div>
      </div>
    );
  }

  // 菜单类型（menu）可点击
  return (
    <Link to={menu.path} className="menu-item">
      <Icon type={menu.icon} />
      <span>{menu.name}</span>
    </Link>
  );
}

// Sidebar.tsx - 侧边栏组件
import { useEffect, useState } from 'react';

function Sidebar() {
  const [menus, setMenus] = useState<MenuTreeItem[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getUserMenus()
      .then(data => {
        setMenus(data.menus);
        setLoading(false);
      })
      .catch(err => {
        console.error('Failed to load menus:', err);
        setLoading(false);
      });
  }, []);

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <aside className="sidebar">
      <nav>
        {menus.map(menu => (
          <MenuItem key={menu.id} menu={menu} />
        ))}
      </nav>
    </aside>
  );
}

export default Sidebar;
```

### Vue 3 示例

```vue
<template>
  <aside class="sidebar">
    <nav>
      <MenuItem v-for="menu in menus" :key="menu.id" :menu="menu" />
    </nav>
  </aside>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import MenuItem from './MenuItem.vue';

interface MenuTreeItem {
  id: string;
  name: string;
  path: string;
  icon: string;
  type: 'dir' | 'menu' | 'link';
  children?: MenuTreeItem[];
}

const menus = ref<MenuTreeItem[]>([]);

async function getUserMenus() {
  const response = await fetch('/api/menus/user/tree', {
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('token')}`
    }
  });
  return response.json();
}

onMounted(async () => {
  const data = await getUserMenus();
  menus.value = data.menus;
});
</script>
```

---

## 🔧 扩展建议

### 1. 细粒度权限中间件

可以添加基于权限的中间件（当前仅有基于角色的 Admin 中间件）：

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
cacheKey := fmt.Sprintf("user:menus:%s", userID)
if cached, err := cache.Get(ctx, cacheKey); err == nil {
    return cached, nil
}

menuTree, err := domainService.GetUserMenuTree(ctx, userID)
if err != nil {
    return nil, err
}

cache.Set(ctx, cacheKey, menuTree, 5*time.Minute)
return menuTree, nil
```

### 3. 审计日志

```go
// 记录敏感操作
auditLog.Record(ctx, AuditLog{
    UserID:   adminID,
    Action:   "AssignRole",
    Resource: fmt.Sprintf("user:%s,role:%s", userID, roleID),
    IP:       c.ClientIP(),
})
```

### 4. 自定义权限数据

修改 YAML 文件后重新执行 seed：

```bash
# 编辑 seed 数据
vi internal/infrastructure/seed/data/roles.yaml
vi internal/infrastructure/seed/data/menus.yaml

# 清空 seed 历史
DELETE FROM seed_history WHERE name = 'RBAC_SEED';

# 重新执行 seed
./main migrate seed
```

---

## 📖 相关文档

- [`docs/RBAC_IMPLEMENTATION.md`](./RBAC_IMPLEMENTATION.md) - RBAC 实现指南（技术细节）
- [`docs/SEED_USAGE.md`](./SEED_USAGE.md) - Seed 数据使用指南
- [`docs/ADMIN_MIDDLEWARE_SETUP.md`](./ADMIN_MIDDLEWARE_SETUP.md) - Admin 中间件配置
- [`docs/API_REFACTORING.md`](./API_REFACTORING.md) - API 路由架构总结
- [`CLAUDE.md`](../CLAUDE.md) - 完整项目架构文档
- [`README.md`](../README.md) - 项目快速开始

---

## 🎉 总结

RBAC 系统已完全集成并可直接使用！

### ✅ 核心功能
- 基于角色的访问控制（4个默认角色）
- 细粒度权限管理（17个默认权限）
- 动态菜单树（7个默认菜单）
- Admin 权限中间件（完全实现）
- Seed 数据初始化（支持幂等性）

### 🚀 立即开始
```bash
./main migrate up      # 创建表结构
./main migrate seed    # 初始化数据
./main api             # 启动服务
```

### 📱 前端集成
调用 `GET /api/menus/user/tree` 获取菜单树，直接渲染侧边栏！

---

**最后更新**：2024
**版本**：2.0 - 完全集成，开箱即用
