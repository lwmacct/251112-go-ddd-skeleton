# Go DDD Skeleton

基于 Golang 的领域驱动设计（DDD）项目骨架，集成了用户认证、订单管理等完整功能。

## 技术栈

- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: PostgreSQL
- **缓存**: Redis
- **认证**: JWT + 2FA (TOTP)
- **日志**: Zap

## 项目结构

```
251112-go-ddd-skeleton/
├── main.go                # 应用程序主入口（urfave/cli）
├── internal/              # 私有应用代码
│   ├── commands/          # CLI 子命令（业务入口）
│   │   ├── api/          # REST API 服务命令
│   │   ├── worker/       # 后台任务命令
│   │   └── migrate/      # 数据库迁移命令
│   ├── domain/            # 领域层（业务逻辑）
│   ├── application/       # 应用层（用例）
│   ├── infrastructure/    # 基础设施层
│   │   ├── seed/         # 种子数据（RBAC 初始化）
│   │   ├── persistence/  # 持久化
│   │   ├── cache/        # Redis 缓存
│   │   ├── auth/         # 认证（JWT、密码哈希）
│   │   └── ...
│   ├── adapters/          # 适配器层（HTTP）
│   ├── shared/            # 共享工具
│   ├── config/            # 配置管理
│   └── bootstrap/         # 应用启动和依赖注入
├── configs/               # 配置文件
├── docs/                  # 文档
│   ├── RBAC_IMPLEMENTATION.md
│   ├── ADMIN_MIDDLEWARE_SETUP.md
│   └── SEED_USAGE.md     # Seed 使用指南
├── docker/                # Docker 配置
└── README.md
```

## 快速开始

### 前置要求

- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose (可选)

### 使用 Docker Compose

```bash
# 启动所有服务
docker-compose -f docker/docker-compose.yaml up -d

# 运行数据库迁移
docker exec -it go-ddd-api ./main migrate up

# 初始化种子数据
docker exec -it go-ddd-api ./main migrate seed

# 查看日志
docker-compose -f docker/docker-compose.yaml logs -f api
```

### 本地开发

1. **安装依赖**

```bash
go mod download
```

2. **配置环境变量**

```bash
cp .env.example .env
# 编辑 .env 文件，配置数据库和Redis连接
```

3. **运行数据库迁移和初始化数据**

```bash
# 执行迁移
go run main.go migrate up

# 初始化种子数据（包括默认角色、权限、菜单和管理员账户）
go run main.go migrate seed

# 查看迁移状态
go run main.go migrate status
```

**默认管理员账户**：
- Email: `admin@example.com`
- Password: `Admin@123456`

⚠️ **重要**: 生产环境部署后请立即修改默认密码！

4. **启动 API 服务**

```bash
# 使用默认配置（0.0.0.0:8080）
go run main.go api

# 或指定端口和主机
go run main.go api --host 127.0.0.1 --port 8080
```

5. **启动 Worker（可选）**

```bash
# 使用默认间隔（1分钟）
go run main.go worker

# 或自定义执行间隔
go run main.go worker --interval 30s
```

6. **编译独立二进制文件（可选）**

```bash
# 编译主程序（包含所有命令）
go build -o main .

# 运行各种命令
./main api --host 127.0.0.1 --port 8080
./main migrate up
./main migrate seed
./main worker --interval 30s

# 或者编译为独立的二进制文件
go build -o bin/api cmd/api/main.go
go build -o bin/worker cmd/worker/main.go
go build -o bin/migrate cmd/migrate/main.go

./bin/api --host 127.0.0.1 --port 8080
./bin/migrate up
./bin/migrate seed
./bin/worker --interval 30s
```

## API 端点

### 认证

- `POST /api/v1/auth/register` - 注册用户
- `POST /api/v1/auth/login` - 登录
- `POST /api/v1/auth/refresh` - 刷新令牌
- `POST /api/v1/auth/logout` - 登出

### 用户个人中心

- `GET /api/v1/user` - 获取当前用户信息
- `PUT /api/v1/user` - 更新当前用户信息
- `PATCH /api/v1/user` - 部分更新当前用户
- `DELETE /api/v1/user` - 注销账号
- `PUT /api/v1/user/password` - 修改密码
- `PUT /api/v1/user/email` - 修改邮箱
- `POST /api/v1/user/avatar` - 上传头像

### 安全与会话管理

- `GET /api/v1/user/sessions` - 查看活跃会话
- `DELETE /api/v1/user/sessions/:id` - 删除指定会话
- `GET /api/v1/user/tokens` - 查看个人访问令牌
- `POST /api/v1/user/tokens` - 创建个人访问令牌
- `DELETE /api/v1/user/tokens/:id` - 撤销令牌
- `POST /api/v1/user/2fa/enable` - 启用双因素认证
- `POST /api/v1/user/2fa/disable` - 禁用双因素认证

### 管理员接口

- `GET /api/v1/admin/users` - 列出所有用户
- `GET /api/v1/admin/users/:id` - 获取指定用户信息
- `PUT /api/v1/admin/users/:id` - 更新指定用户信息
- `DELETE /api/v1/admin/users/:id` - 删除用户
- `POST /api/v1/admin/users/:id/ban` - 封禁用户
- `POST /api/v1/admin/users/:id/unban` - 解封用户

### 订单管理

- `POST /api/v1/orders` - 创建订单
- `GET /api/v1/orders` - 列出当前用户订单
- `GET /api/v1/orders/:id` - 获取订单详情
- `POST /api/v1/orders/:id/cancel` - 取消订单
- `POST /api/v1/orders/:id/payment` - 处理支付
- `GET /api/v1/orders/:id/payment` - 获取支付信息
- `POST /api/v1/orders/:id/shipment` - 创建发货
- `GET /api/v1/orders/:id/shipment` - 获取发货信息

### 管理员订单接口

- `GET /api/v1/admin/orders` - 列出所有订单
- `GET /api/v1/admin/orders/:id` - 获取任意订单详情
- `PUT /api/v1/admin/orders/:id/status` - 更新订单状态

### RBAC 权限管理接口

#### 菜单管理
- `GET /api/menus/user/tree` - 获取当前用户菜单树（前端侧边栏）
- `POST /api/admin/menus` - 创建菜单
- `PUT /api/admin/menus/:id` - 更新菜单
- `DELETE /api/admin/menus/:id` - 删除菜单
- `GET /api/admin/menus/tree` - 获取所有菜单树
- `PUT /api/admin/menus/order` - 更新菜单排序

#### 角色管理
- `POST /api/admin/roles` - 创建角色
- `PUT /api/admin/roles/:id` - 更新角色
- `DELETE /api/admin/roles/:id` - 删除角色
- `GET /api/admin/roles/:id` - 获取角色详情
- `GET /api/admin/roles` - 列出所有角色

#### 角色权限关联
- `POST /api/admin/roles/:roleId/permissions` - 为角色分配权限
- `GET /api/admin/roles/:roleId/permissions` - 获取角色权限列表
- `POST /api/admin/roles/:roleId/menus` - 为角色分配菜单
- `GET /api/admin/roles/:roleId/menus` - 获取角色菜单树

#### 用户角色管理
- `POST /api/admin/users/:userId/roles/:roleId` - 为用户分配角色
- `DELETE /api/admin/users/:userId/roles/:roleId` - 移除用户角色
- `GET /api/admin/users/:userId/roles` - 获取用户角色列表

## 核心设计

### DDD 分层架构

- **Domain**: 纯业务逻辑，无外部依赖
- **Application**: 用例编排，实现 CQRS 模式
- **Infrastructure**: 技术实现（数据库、缓存等）
- **Adapters**: 外部接口（HTTP、gRPC 等）

### CQRS 模式

- Commands（命令）: 写操作
- Queries（查询）: 读操作
- 逻辑分离但不过度细分文件

### 依赖倒置

- Domain 层定义接口
- Infrastructure 层实现接口
- Application 层定义端口
- Infrastructure 层实现端口

## 数据库设计

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

## 数据库初始化

### 迁移和 Seed

```bash
# 1. 创建数据库表结构
./main migrate up

# 2. 初始化基础数据（角色、权限、菜单、管理员账户）
./main migrate seed

# 3. 查看执行状态
./main migrate status
```

### 默认数据

执行 `seed` 命令后将自动创建：
- **4 个角色**: admin, user, editor, viewer
- **17 个权限**: user:*, role:*, menu:*, order:*, permission:read
- **7 个菜单**: 系统管理（含 4 个子菜单）、订单管理、个人中心
- **1 个管理员**: admin@example.com / Admin@123456

详细说明请参考：`docs/SEED_USAGE.md`

## 开发指南

### 添加新功能

1. 在 `internal/domain/` 创建领域实体和接口
2. 在 `internal/application/` 实现用例（Commands + Queries）
3. 在 `internal/infrastructure/` 实现技术细节
4. 在 `internal/adapters/http/handler/` 添加 HTTP 处理器
5. 在 `internal/adapters/http/router.go` 注册路由
6. 在 `internal/bootstrap/container.go` 更新依赖注入

### 添加新的 CLI 命令

1. 在 `internal/commands/` 创建新的命令包
2. 定义 `Command` 变量（类型为 `*cli.Command`）
3. 在 `main.go` 的 `buildCommands()` 函数中引入新命令

### 运行测试

```bash
go test ./...
```

### 代码检查

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run
```

## License

MIT License
