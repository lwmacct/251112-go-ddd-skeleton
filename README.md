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
│   ├── adapters/          # 适配器层（HTTP）
│   ├── shared/            # 共享工具
│   ├── config/            # 配置管理
│   └── bootstrap/         # 应用启动和依赖注入
├── configs/               # 配置文件
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
docker exec -it go-ddd-api /app/migrate up

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

3. **运行数据库迁移**

```bash
# 执行迁移
go run cmd/migrate/main.go up

# 查看迁移状态
go run cmd/migrate/main.go status
```

4. **启动 API 服务**

```bash
# 使用默认配置（0.0.0.0:8080）
go run cmd/api/main.go

# 或指定端口和主机
go run cmd/api/main.go --host 127.0.0.1 --port 8080
```

5. **启动 Worker（可选）**

```bash
# 使用默认间隔（1分钟）
go run cmd/worker/main.go

# 或自定义执行间隔
go run cmd/worker/main.go --interval 30s
```

6. **编译独立二进制文件（可选）**

```bash
# 编译 API 服务
go build -o bin/api cmd/api/main.go

# 编译 Worker
go build -o bin/worker cmd/worker/main.go

# 编译迁移工具
go build -o bin/migrate cmd/migrate/main.go

# 运行编译后的二进制文件
./bin/api --host 127.0.0.1 --port 8080
./bin/migrate up
./bin/worker --interval 30s
```

## API 端点

### 认证

- `POST /api/v1/auth/register` - 注册用户
- `POST /api/v1/auth/login` - 登录
- `POST /api/v1/auth/refresh` - 刷新令牌
- `POST /api/v1/auth/logout` - 登出

### 用户

- `GET /api/v1/users/:id` - 获取用户信息
- `PUT /api/v1/users/:id` - 更新用户信息
- `DELETE /api/v1/users/:id` - 删除用户
- `GET /api/v1/users` - 列出用户

### 订单

- `POST /api/v1/orders` - 创建订单
- `GET /api/v1/orders/:id` - 获取订单详情
- `GET /api/v1/orders` - 列出订单
- `POST /api/v1/orders/:id/cancel` - 取消订单
- `POST /api/v1/orders/:id/payment` - 处理支付
- `POST /api/v1/orders/:id/shipment` - 创建发货

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
golangci-lint run
```

## License

MIT License
