# Golang DDD 项目架构

## 技术栈

- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: PostgreSQL
- **缓存**: Redis
- **依赖注入**: Wire (可选) 或手动装配

## 核心设计原则

### 1. DDD 分层架构

```
adapters → application → domain
              ↑
infrastructure (实现 application 的 ports)
```

### 2. CQRS 模式

- **命令（Command）**: 写操作，在 `commands.go` 中
- **查询（Query）**: 读操作，在 `queries.go` 中
- 逻辑分离但不过度细分文件

### 3. 依赖倒置

- Domain 层定义接口（Repository）
- Infrastructure 层实现接口
- Application 层定义端口（Ports）
- Infrastructure 层实现端口

---

## 目录结构

```
251112-go-ddd-skeleton/
│
├── main.go                                 # 应用程序统一 CLI 入口（可选，包含所有命令）
│
├── cmd/                                    # 独立的命令入口（各自 main 包）
│   ├── api/
│   │   └── main.go                        # REST API 服务入口
│   ├── worker/
│   │   └── main.go                        # 后台任务处理器入口
│   └── migrate/
│       └── main.go                        # 数据库迁移工具入口
│
├── internal/                               # 私有应用代码
│   │
│   ├── commands/                           # CLI 子命令实现（业务逻辑）
│   │   ├── api/
│   │   │   └── api.go                     # API 服务命令逻辑
│   │   ├── worker/
│   │   │   └── worker.go                  # Worker 命令逻辑
│   │   └── migrate/
│   │       └── migrate.go                 # 迁移命令逻辑
│   │
│   ├── domain/                             # 领域层（纯业务逻辑，无外部依赖）
│   │   │
│   │   ├── user/                           # 用户聚合
│   │   │   ├── user.go                    # User 实体 + 值对象（Email, Password, Username）
│   │   │   ├── repository.go              # 仓储接口定义
│   │   │   ├── service.go                 # 领域服务（可选）
│   │   │   └── errors.go                  # 领域错误
│   │   │
│   │   ├── auth/                           # 认证聚合
│   │   │   ├── auth.go                    # TwoFactor + PAT + Session 实体
│   │   │   ├── repository.go              # 认证仓储接口
│   │   │   ├── service.go                 # 认证领域服务
│   │   │   └── errors.go
│   │   │
│   │   ├── order/                          # 订单聚合
│   │   │   ├── order.go                   # Order + OrderItem + 值对象
│   │   │   ├── payment.go                 # Payment 实体
│   │   │   ├── shipment.go                # Shipment 实体
│   │   │   ├── repository.go              # 订单仓储接口
│   │   │   ├── service.go                 # 订单领域服务
│   │   │   └── errors.go
│   │   │
│   │   └── rbac/                           # RBAC 权限控制聚合
│   │       ├── role.go                    # Role 实体
│   │       ├── permission.go              # Permission 实体
│   │       ├── menu.go                    # Menu 实体（支持两层树形结构）
│   │       ├── repository.go              # RBAC 仓储接口
│   │       ├── service.go                 # RBAC 领域服务（权限检查、菜单树构建）
│   │       └── errors.go                  # 领域错误
│   │
│   ├── application/                        # 应用层（用例编排，CQRS 体现）
│   │   │
│   │   ├── user/                           # 用户应用服务
│   │   │   ├── service.go                 # 应用服务（协调命令和查询）
│   │   │   ├── commands.go                # 写操作：CreateUser, UpdateUser, DeleteUser, ChangePassword
│   │   │   ├── queries.go                 # 读操作：GetUser, ListUsers, FindByEmail
│   │   │   └── dto.go                     # 数据传输对象
│   │   │
│   │   ├── auth/                           # 认证应用服务
│   │   │   ├── service.go
│   │   │   ├── commands.go                # Login, Logout, Enable2FA, Verify2FA, CreatePAT, RevokePAT
│   │   │   ├── queries.go                 # GetSessions, ListPATs
│   │   │   └── dto.go
│   │   │
│   │   ├── order/                          # 订单应用服务
│   │   │   ├── service.go
│   │   │   ├── commands.go                # CreateOrder, CancelOrder, ProcessPayment, RefundPayment, CreateShipment
│   │   │   ├── queries.go                 # GetOrder, ListOrders, GetPayment, GetShipment
│   │   │   └── dto.go
│   │   │
│   │   ├── role/                           # 角色应用服务
│   │   │   ├── service.go                 # 应用服务
│   │   │   ├── dto.go                     # 数据传输对象
│   │   │   └── (待实现 commands/queries)  # 角色管理相关操作
│   │   │
│   │   └── menu/                           # 菜单应用服务
│   │       ├── commands.go                # CreateMenu, UpdateMenu, DeleteMenu
│   │       ├── queries.go                 # GetUserMenuTree, GetAllMenuTree, GetRoleMenuTree
│   │       └── dto.go                     # 菜单树 DTO（支持递归 children）
│   │
│   ├── infrastructure/                     # 基础设施层
│   │   │
│   │   ├── seed/                           # 数据库种子数据（Seed）
│   │   │   ├── seeder.go                  # Seeder 接口和 Manager
│   │   │   ├── rbac_seeder.go             # RBAC seed 实现
│   │   │   ├── user_seeder.go             # 用户 seed 实现
│   │   │   ├── utils.go                   # 工具函数
│   │   │   └── data/                      # seed 数据文件（YAML）
│   │   │       ├── roles.yaml            # 默认角色
│   │   │       ├── permissions.yaml      # 默认权限
│   │   │       ├── menus.yaml            # 默认菜单
│   │   │       ├── role_permissions.yaml # 角色-权限关联
│   │   │       ├── role_menus.yaml       # 角色-菜单关联
│   │   │       └── users.yaml            # 默认用户
│   │   │
│   │   ├── persistence/                    # 持久化
│   │   │   ├── postgres.go                # PostgreSQL 连接与配置
│   │   │   ├── transaction.go             # 事务管理（UnitOfWork）
│   │   │   │
│   │   │   ├── model/                     # GORM 模型（集中管理，按领域分文件）
│   │   │   │   ├── user.go               # User 模型
│   │   │   │   ├── two_factor.go         # TwoFactor 模型
│   │   │   │   ├── personal_access_token.go # PersonalAccessToken 模型
│   │   │   │   ├── session.go            # Session 模型
│   │   │   │   ├── order.go              # Order 模型
│   │   │   │   ├── order_item.go         # OrderItem 模型
│   │   │   │   ├── payment.go            # Payment 模型
│   │   │   │   ├── shipment.go           # Shipment 模型
│   │   │   │   ├── invoice.go            # Invoice 模型
│   │   │   │   ├── role.go               # Role 模型
│   │   │   │   ├── permission.go         # Permission 模型
│   │   │   │   ├── menu.go               # Menu 模型
│   │   │   │   ├── user_role.go          # UserRole 关联表
│   │   │   │   ├── role_permission.go    # RolePermission 关联表
│   │   │   │   ├── role_menu.go          # RoleMenu 关联表
│   │   │   │   └── schema.go             # 统一模型注册（AutoMigrate）
│   │   │   │
│   │   │   ├── mapper/                    # Domain Entity <-> GORM Model 转换
│   │   │   │   ├── user.go
│   │   │   │   ├── auth.go
│   │   │   │   ├── order.go
│   │   │   │   └── rbac.go               # RBAC 实体映射
│   │   │   │
│   │   │   └── repository/                # 仓储实现
│   │   │       ├── user_repo.go          # 实现 domain/user/repository.go
│   │   │       ├── auth_repo.go          # 实现 domain/auth/repository.go
│   │   │       ├── order_repo.go         # 实现 domain/order/repository.go
│   │   │       ├── role_repo.go          # 实现 domain/rbac 角色仓储
│   │   │       ├── permission_repo.go    # 实现 domain/rbac 权限仓储
│   │   │       └── menu_repo.go          # 实现 domain/rbac 菜单仓储
│   │   │
│   │   ├── cache/                          # Redis 缓存
│   │   │   ├── redis.go                   # Redis 客户端配置
│   │   │   ├── session.go                 # Session 存储
│   │   │   └── distributed_lock.go        # 分布式锁
│   │   │
│   │   ├── auth/                           # 认证基础设施（实现 application 端口）
│   │   │   ├── jwt.go                     # JWT Token 实现
│   │   │   ├── password.go                # 密码哈希
│   │   │   └── totp.go                    # TOTP 2FA
│   │   │
│   │   ├── email/                          # 邮件服务
│   │   │   └── smtp.go                    # SMTP 实现
│   │   │
│   │   ├── payment/                        # 支付网关
│   │   │   └── stripe.go                  # Stripe 实现
│   │   │
│   │   └── logger/                         # 日志
│   │       └── logger.go                  # Zap 或其他日志实现
│   │
│   ├── adapters/                           # 适配器层（接口层）
│   │   │
│   │   └── http/                           # HTTP 适配器（Gin）
│   │       ├── server.go                  # HTTP 服务器
│   │       ├── router.go                  # 路由配置
│   │       │
│   │       ├── middleware/                # 中间件
│   │       │   ├── auth.go               # JWT 认证中间件
│   │       │   ├── admin.go              # 管理员权限中间件
│   │       │   ├── cors.go               # CORS
│   │       │   ├── rate_limit.go         # 限流（待实现）
│   │       │   └── logger.go             # 日志中间件
│   │       │
│   │       ├── handler/                   # HTTP 处理器（按模块分子目录）
│   │       │   ├── user/                 # 用户模块
│   │       │   │   └── handler.go        # 用户相关所有端点
│   │       │   ├── auth/                 # 认证模块
│   │       │   │   └── handler.go        # 认证相关所有端点
│   │       │   ├── order/                # 订单模块
│   │       │   │   └── handler.go        # 订单相关所有端点
│   │       │   └── rbac/                 # RBAC 权限控制模块
│   │       │       ├── role_handler.go   # 角色管理端点
│   │       │       └── menu_handler.go   # 菜单管理端点
│   │       │
│   │       └── response/                  # 统一响应
│   │           └── response.go           # 响应格式 + 错误处理
│   │
│   ├── shared/                             # 内部共享代码
│   │   ├── errors/                         # 错误处理
│   │   │   └── errors.go                  # 错误类型 + 错误码
│   │   ├── validator/                      # 验证器
│   │   │   └── validator.go               # 基于 go-playground/validator
│   │   └── pagination/                     # 分页
│   │       └── pagination.go
│   │
│   ├── config/                             # 配置管理
│   │   └── config.go                      # 配置结构 + 加载器（Viper）
│   │
│   └── bootstrap/                          # 应用启动和依赖注入
│       └── container.go                   # 依赖注入容器（手动装配）
│
├── api/                                    # API 定义（独立于实现）
│   └── openapi/                           # OpenAPI/Swagger 规范
│       └── openapi.yaml                   # OpenAPI 3.0 规范
│
├── configs/                                # 配置文件
│   ├── config.yaml                        # 默认配置
│   ├── config.dev.yaml                    # 开发环境
│   └── config.prod.yaml                   # 生产环境
│
├── migrations/                             # 数据库迁移（集中管理）
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_auth_tables.up.sql
│   ├── 000002_create_auth_tables.down.sql
│   ├── 000003_create_orders_table.up.sql
│   ├── 000003_create_orders_table.down.sql
│   ├── 000004_create_order_details_tables.up.sql
│   ├── 000004_create_order_details_tables.down.sql
│   ├── 000005_create_rbac_tables.up.sql
│   └── 000005_create_rbac_tables.down.sql
│
├── scripts/                                # 脚本工具
│   ├── setup.sh                           # 项目初始化
│   └── dev.sh                             # 本地开发启动
│
├── test/                                   # 测试
│   ├── integration/                       # 集成测试
│   │   ├── user_test.go
│   │   ├── auth_test.go
│   │   └── order_test.go
│   └── fixtures/                          # 测试数据
│       └── fixtures.go
│
├── docs/                                   # 文档
│   ├── RBAC_IMPLEMENTATION.md             # RBAC 实现指南
│   ├── RBAC_INTEGRATION.md                # RBAC 集成说明
│   ├── ADMIN_MIDDLEWARE_SETUP.md          # Admin 中间件配置指南
│   ├── SEED_USAGE.md                      # 数据库 Seed 使用指南
│   └── architecture/                      # 架构文档
│       ├── README.md                     # 架构概览
│       └── ddd.md                        # DDD 设计说明
│
├── docker/                                 # Docker 相关
│   ├── Dockerfile                         # 生产镜像
│   └── docker-compose.yaml                # 本地开发环境
│
├── go.mod                                  # Go 模块
├── go.sum
├── Taskfile.yml                            # Task 命令
├── .env.example                            # 环境变量示例
├── .gitignore
├── .golangci.yml                           # golangci-lint 配置
└── README.md
```

---

## 数据库表设计

### 用户相关表

- `users` - 用户基本信息
- `two_factor_auth` - 双因素认证
- `personal_access_tokens` - 个人访问令牌
- `sessions` - 用户会话

### 订单相关表

- `orders` - 订单主表
- `order_items` - 订单明细
- `payments` - 支付记录
- `shipments` - 发货记录
- `invoices` - 发票记录

### RBAC 权限控制相关表

- `roles` - 角色表
- `permissions` - 权限表
- `menus` - 菜单表（支持两层树形结构）
- `user_roles` - 用户-角色关联表（多对多）
- `role_permissions` - 角色-权限关联表（多对多）
- `role_menus` - 角色-菜单关联表（多对多）

### 系统相关表

- `seed_history` - 种子数据执行历史记录表

---

## 关键要点

### ✅ 做到了

1. **DDD 分层清晰**：Domain → Application → Infrastructure → Adapters
2. **CQRS 分离**：Commands 和 Queries 独立文件但不过度细分
3. **依赖倒置**：Domain 定义接口，Infrastructure 实现
4. **GORM 模型集中**：避免循环依赖，模型按实体细粒度拆分
5. **端口-适配器**：Application 定义端口，Infrastructure 实现
6. **符合 Go 习惯**：包名简洁，文件适度聚合
7. **独立入口**：`cmd/` 下各服务独立入口，减少依赖耦合
8. **RBAC 权限控制**：完整的基于角色的访问控制系统，支持菜单树和权限管理
9. **Handler 模块化**：HTTP 处理器按业务模块划分子目录，结构更清晰
10. **Seed 数据初始化**：支持幂等性的数据库种子数据初始化，包含默认角色、权限、菜单和管理员账户

### ⚠️ 注意事项

1. **Domain 层**不应该依赖任何框架（无 GORM tags）
2. **Application 层**不应该 import Infrastructure 的具体实现
3. **事务管理**应该在 Application 层声明边界
4. **错误处理**统一使用领域错误 + HTTP 状态码映射
5. **测试**每一层都应该有对应的测试
6. **Handler 组织**：按业务模块划分子目录，每个模块一个 handler.go 文件
7. **RBAC 集成**：菜单、角色、权限三者独立管理，通过关联表建立多对多关系

---

## RBAC 权限控制系统

项目已实现完整的 RBAC（基于角色的访问控制）系统，包括以下核心功能：

### 核心实体

1. **Role（角色）**：定义用户组权限集合
2. **Permission（权限）**：细粒度的操作权限（如 `user:create`, `order:read`）
3. **Menu（菜单）**：前端菜单树结构，支持两层嵌套

### 主要特性

- **多对多关系**：用户-角色、角色-权限、角色-菜单
- **菜单树查询**：`GetUserMenuTree()` API 返回用户可见的完整菜单树
- **权限检查**：`CheckPermission(userID, permissionCode)` 验证用户权限
- **两层菜单结构**：支持目录（dir）和菜单（menu）两层结构，优化性能

### 使用示例

获取用户菜单树（前端渲染侧边栏）：

```go
// 在 HTTP Handler 中
userID := c.GetString("user_id") // 从 JWT 获取
menuTree, err := menuService.GetUserMenuTree(ctx, userID)
// 返回 JSON 格式的树形结构，包含 children 递归嵌套
```

详细实现指南请参考：`docs/RBAC_IMPLEMENTATION.md`

---

## 数据库 Seed 初始化

项目提供完整的数据库 seed 功能，用于初始化 RBAC 系统的基础数据。

### 功能特性

- ✅ **幂等性**：支持多次安全执行，不会重复插入数据
- ✅ **事务保证**：所有操作在事务中执行，失败自动回滚
- ✅ **执行历史**：记录到 `seed_history` 表，可追溯历史
- ✅ **依赖管理**：自动处理数据依赖关系

### 初始化数据

执行 `./main migrate seed` 将初始化以下数据：

#### 默认角色（4 个）
- `admin` - 超级管理员（所有权限）
- `user` - 普通用户（基础业务权限）
- `editor` - 编辑员（内容编辑权限）
- `viewer` - 访客（只读权限）

#### 默认权限（17 个）
- 用户管理：`user:create`, `user:read`, `user:update`, `user:delete`
- 角色管理：`role:create`, `role:read`, `role:update`, `role:delete`
- 菜单管理：`menu:create`, `menu:read`, `menu:update`, `menu:delete`
- 权限管理：`permission:read`
- 订单管理：`order:create`, `order:read`, `order:update`, `order:delete`

#### 默认菜单（7 个，两层树形结构）
```
系统管理/ (目录)
├── 用户管理 (菜单)
├── 角色管理 (菜单)
├── 菜单管理 (菜单)
└── 权限管理 (菜单)

订单管理 (菜单)
个人中心 (菜单)
```

#### 默认管理员账户
- **Email**: `admin@example.com`
- **Password**: `Admin@123456`
- **角色**: `admin`

⚠️ **生产环境部署后请立即修改默认密码！**

### 使用方法

```bash
# 1. 执行数据库迁移
./main migrate up

# 2. 初始化 seed 数据
./main migrate seed

# 3. 验证数据（查询数据库）
SELECT * FROM roles;
SELECT * FROM seed_history;
```

详细使用指南请参考：`docs/SEED_USAGE.md`

---

## 扩展阅读

- [Domain-Driven Design（Eric Evans）](https://www.domainlanguage.com/ddd/)
- [Implementing Domain-Driven Design（Vaughn Vernon）](https://vaughnvernon.com/)
- [Go DDD Example](https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example)
- [Clean Architecture in Go](https://github.com/bxcodec/go-clean-arch)
