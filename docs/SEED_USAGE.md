# 数据库 Seed 使用指南

## 概述

本项目提供了完整的数据库 seed 功能，用于初始化 RBAC 权限系统的基础数据，包括角色、权限、菜单和默认管理员账户。

## 功能特性

- ✅ **幂等性**：支持多次安全执行，不会重复插入数据
- ✅ **事务保证**：所有操作在事务中执行，失败自动回滚
- ✅ **执行历史**：记录到 `seed_history` 表，可追溯历史
- ✅ **依赖管理**：自动处理数据依赖关系（角色→权限→菜单→用户）
- ✅ **详细日志**：友好的执行过程输出

## 快速开始

### 1. 执行数据库迁移

在运行 seed 之前，需要先创建数据库表结构：

```bash
./main migrate up
```

### 2. 执行 seed 数据初始化

```bash
./main migrate seed
```

预期输出：

```
Connecting to database: postgres@localhost:5432/go_ddd_db
🌱 Running database seeders...
▶️  Running seed: RBAC_SEED
  📦 Loading YAML data...
  ✓ Creating roles...
  ✓ Creating permissions...
  ✓ Creating menus...
  ✓ Assigning permissions to roles...
  ✓ Assigning menus to roles...
  ✅ Created 4 roles, 17 permissions, 7 menus
✅ Seed completed: RBAC_SEED

▶️  Running seed: DEFAULT_ADMIN_USER
  📦 Loading user data...
  ✓ Created user: admin@example.com
  ✓ Assigned role: admin
  ✅ Created 1 default user(s)
✅ Seed completed: DEFAULT_ADMIN_USER

✅ All seeds completed successfully!

╔══════════════════════════════════════════════════════╗
║                   🎉 Seed 完成！                     ║
╚══════════════════════════════════════════════════════╝

📋 默认管理员账户:
   Email:    admin@example.com
   Password: Admin@123456

⚠️  请在首次登录后立即修改默认密码！
```

### 3. 验证 seed 结果

#### 方式一：使用 SQL 查询

```sql
-- 查看角色
SELECT * FROM roles;

-- 查看权限
SELECT * FROM permissions;

-- 查看菜单树
SELECT
    m1.name AS parent_menu,
    m2.name AS child_menu
FROM menus m1
LEFT JOIN menus m2 ON m2.parent_id = m1.id
WHERE m1.parent_id IS NULL
ORDER BY m1.sort_order, m2.sort_order;

-- 查看默认用户及其角色
SELECT
    u.email,
    u.username,
    r.name AS role
FROM users u
JOIN user_roles ur ON ur.user_id = u.id
JOIN roles r ON r.id = ur.role_id
WHERE u.email = 'admin@example.com';

-- 查看 seed 执行历史
SELECT * FROM seed_history ORDER BY executed_at DESC;
```

#### 方式二：通过 API 测试

```bash
# 1. 登录获取 token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "Admin@123456"
  }'

# 2. 获取用户菜单树
curl -X GET http://localhost:8080/api/menus/user/tree \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 3. 访问管理员端点测试权限
curl -X GET http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 初始化的数据

### 1. 角色（4 个）

| Code | Name | Description |
|------|------|-------------|
| `admin` | 超级管理员 | 拥有系统所有权限 |
| `user` | 普通用户 | 基础业务权限 |
| `editor` | 编辑员 | 内容编辑权限 |
| `viewer` | 访客 | 只读权限 |

### 2. 权限（17 个）

#### 用户管理权限
- `user:create` - 创建用户
- `user:read` - 查看用户
- `user:update` - 更新用户
- `user:delete` - 删除用户

#### 角色管理权限
- `role:create` - 创建角色
- `role:read` - 查看角色
- `role:update` - 更新角色
- `role:delete` - 删除角色

#### 菜单管理权限
- `menu:create` - 创建菜单
- `menu:read` - 查看菜单
- `menu:update` - 更新菜单
- `menu:delete` - 删除菜单

#### 权限管理权限
- `permission:read` - 查看权限

#### 订单管理权限
- `order:create` - 创建订单
- `order:read` - 查看订单
- `order:update` - 更新订单
- `order:delete` - 删除订单

### 3. 菜单树（7 个菜单，两层结构）

```
系统管理/ (目录)
├── 用户管理 (菜单) - 需要 user:read 权限
├── 角色管理 (菜单) - 需要 role:read 权限
├── 菜单管理 (菜单) - 需要 menu:read 权限
└── 权限管理 (菜单) - 需要 permission:read 权限

订单管理 (菜单) - 需要 order:read 权限

个人中心 (菜单) - 无权限要求
```

### 4. 角色权限映射

| 角色 | 权限 |
|-----|------|
| `admin` | 所有权限（17 个） |
| `editor` | 部分读写权限（6 个） |
| `user` | 基础业务权限（2 个） |
| `viewer` | 只读权限（5 个） |

### 5. 角色菜单映射

| 角色 | 可见菜单 |
|-----|---------|
| `admin` | 所有菜单（7 个） |
| `editor` | 部分管理菜单（4 个） |
| `user` | 业务菜单（2 个） |
| `viewer` | 只读菜单（4 个） |

### 6. 默认用户

- **Email**: `admin@example.com`
- **Password**: `Admin@123456`
- **角色**: `admin`
- **状态**: 激活

## 高级用法

### 重复执行（幂等性）

如果再次执行 seed 命令，系统会自动跳过已执行的 seed：

```bash
./main migrate seed
```

输出：

```
🌱 Running database seeders...
⏭️  Skipping seed: RBAC_SEED (already executed)
⏭️  Skipping seed: DEFAULT_ADMIN_USER (already executed)
✅ All seeds completed successfully!
```

### 查看 seed 历史

可以通过查询数据库查看执行历史：

```sql
SELECT
    name,
    status,
    executed_at,
    CASE WHEN error = '' THEN NULL ELSE error END AS error
FROM seed_history
ORDER BY executed_at DESC;
```

### 强制重新执行（开发中）

如果需要强制重新执行（例如在开发环境中测试），可以先清空 seed_history 表：

```sql
-- 清空 seed 历史（会导致重新执行）
TRUNCATE TABLE seed_history;

-- 或者删除特定 seed 的记录
DELETE FROM seed_history WHERE name = 'RBAC_SEED';
```

然后重新运行 seed：

```bash
./main migrate seed
```

## 自定义 seed 数据

### 修改默认数据

所有 seed 数据存储在 YAML 文件中，可以根据需要修改：

```
internal/infrastructure/seed/data/
├── roles.yaml              # 修改角色
├── permissions.yaml        # 修改权限
├── menus.yaml             # 修改菜单
├── role_permissions.yaml  # 修改角色-权限关联
├── role_menus.yaml        # 修改角色-菜单关联
└── users.yaml             # 修改默认用户
```

修改后重新编译并执行：

```bash
go build -o main .
./main migrate seed
```

### 添加新的 Seeder

1. 创建新的 seeder 文件（例如 `internal/infrastructure/seed/custom_seeder.go`）：

```go
package seed

import (
    "context"
    "gorm.io/gorm"
)

type CustomSeeder struct{}

func NewCustomSeeder() *CustomSeeder {
    return &CustomSeeder{}
}

func (s *CustomSeeder) Name() string {
    return "CUSTOM_SEED"
}

func (s *CustomSeeder) ShouldRun(ctx context.Context, db *gorm.DB) (bool, error) {
    // 实现检查逻辑
    return true, nil
}

func (s *CustomSeeder) Run(ctx context.Context, db *gorm.DB) error {
    // 实现 seed 逻辑
    return nil
}
```

2. 在 `migrate.go` 中注册：

```go
manager.Register(seed.NewRBACSeeder())
manager.Register(seed.NewUserSeeder(passwordHasher))
manager.Register(seed.NewCustomSeeder()) // 添加新的 seeder
```

## 常见问题

### Q1: seed 执行失败怎么办？

**A**: seed 使用事务保证，失败会自动回滚。检查错误信息，修复问题后重新执行即可。

### Q2: 如何重置所有 seed 数据？

**A**:

```bash
# 1. 删除所有数据（谨慎操作！）
./main migrate down  # 如果有 down 命令

# 2. 重新迁移
./main migrate up

# 3. 重新 seed
./main migrate seed
```

或者直接清空相关表：

```sql
TRUNCATE TABLE user_roles, role_menus, role_permissions CASCADE;
TRUNCATE TABLE users, roles, permissions, menus CASCADE;
TRUNCATE TABLE seed_history;
```

### Q3: 默认密码是什么？

**A**: `Admin@123456`。强烈建议首次登录后立即修改！

### Q4: 如何添加新的默认角色？

**A**: 编辑 `internal/infrastructure/seed/data/roles.yaml`，添加新角色，然后：

```bash
# 清除 RBAC_SEED 历史
DELETE FROM seed_history WHERE name = 'RBAC_SEED';

# 重新执行
./main migrate seed
```

### Q5: seed 会影响现有数据吗？

**A**: 不会。seed 会检查数据是否已存在（通过 code 字段），已存在的数据会被跳过。

## 安全建议

1. **生产环境部署**：
   - 首次部署后立即修改默认管理员密码
   - 考虑禁用或删除不需要的默认角色
   - 定期审查权限配置

2. **密码策略**：
   - 默认密码使用 bcrypt 哈希（cost=10）
   - 建议在用户服务中添加密码复杂度验证

3. **权限最小化**：
   - 为普通用户分配最少必要权限
   - 定期审查用户权限

4. **审计日志**：
   - `seed_history` 表记录所有执行历史
   - 建议添加应用层面的操作审计日志

## 故障排查

### 问题：seed 执行时提示 "table does not exist"

**解决**：先执行数据库迁移

```bash
./main migrate up
./main migrate seed
```

### 问题：seed 执行失败，提示 "permission denied"

**解决**：检查数据库用户权限

```sql
GRANT ALL PRIVILEGES ON DATABASE your_db TO your_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO your_user;
```

### 问题：重复执行 seed 报错 "duplicate key"

**解决**：这不应该发生（seed 有去重检查）。如果出现，检查数据库唯一索引：

```sql
-- 检查唯一索引
\d+ roles
\d+ permissions
\d+ menus
```

## 性能考虑

- **批量插入**：当前实现逐条插入，对于大量数据可以优化为批量插入
- **索引利用**：利用 code 字段的唯一索引加速查询
- **事务大小**：每个 seeder 在独立事务中执行，避免长事务

## 扩展阅读

- [RBAC 实现指南](./RBAC_IMPLEMENTATION.md)
- [RBAC 集成说明](./RBAC_INTEGRATION.md)
- [Admin 中间件配置](./ADMIN_MIDDLEWARE_SETUP.md)

## 支持

如有问题，请参考：
- 项目 README.md
- Claude Code 架构文档
- GitHub Issues

---

**最后更新**: 2024
**版本**: 1.0
