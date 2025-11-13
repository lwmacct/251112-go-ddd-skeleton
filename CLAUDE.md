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
│   │   └── order/                          # 订单聚合
│   │       ├── order.go                   # Order + OrderItem + 值对象
│   │       ├── payment.go                 # Payment 实体
│   │       ├── shipment.go                # Shipment 实体
│   │       ├── repository.go              # 订单仓储接口
│   │       ├── service.go                 # 订单领域服务
│   │       └── errors.go
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
│   │   └── order/                          # 订单应用服务
│   │       ├── service.go
│   │       ├── commands.go                # CreateOrder, CancelOrder, ProcessPayment, RefundPayment, CreateShipment
│   │       ├── queries.go                 # GetOrder, ListOrders, GetPayment, GetShipment
│   │       └── dto.go
│   │
│   ├── infrastructure/                     # 基础设施层
│   │   │
│   │   ├── persistence/                    # 持久化
│   │   │   ├── postgres.go                # PostgreSQL 连接与配置
│   │   │   ├── transaction.go             # 事务管理（UnitOfWork）
│   │   │   │
│   │   │   ├── model/                     # GORM 模型（集中管理，按领域分文件）
│   │   │   │   ├── user.go               # User + TwoFactor + PAT + Session 模型
│   │   │   │   ├── order.go              # Order + OrderItem + Payment + Shipment + Invoice 模型
│   │   │   │   └── schema.go             # 统一模型注册（AutoMigrate）
│   │   │   │
│   │   │   ├── mapper/                    # Domain Entity <-> GORM Model 转换
│   │   │   │   ├── user.go
│   │   │   │   ├── auth.go
│   │   │   │   └── order.go
│   │   │   │
│   │   │   └── repository/                # 仓储实现
│   │   │       ├── user_repo.go          # 实现 domain/user/repository.go
│   │   │       ├── auth_repo.go          # 实现 domain/auth/repository.go
│   │   │       └── order_repo.go         # 实现 domain/order/repository.go
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
│   │       ├── router.go                  # 路由配置（支持版本化 /v1/）
│   │       │
│   │       ├── middleware/                # 中间件
│   │       │   ├── auth.go               # JWT 认证中间件
│   │       │   ├── cors.go               # CORS
│   │       │   ├── rate_limit.go         # 限流
│   │       │   └── logger.go             # 日志中间件
│   │       │
│   │       ├── handler/                   # HTTP 处理器
│   │       │   ├── user.go               # 用户相关所有端点
│   │       │   ├── auth.go               # 认证相关所有端点
│   │       │   └── order.go              # 订单相关所有端点
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
│   └── 000004_create_order_details_tables.down.sql
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

---

## 关键要点

### ✅ 做到了

1. **DDD 分层清晰**：Domain → Application → Infrastructure → Adapters
2. **CQRS 分离**：Commands 和 Queries 独立文件但不过度细分
3. **依赖倒置**：Domain 定义接口，Infrastructure 实现
4. **GORM 模型集中**：避免循环依赖
5. **端口-适配器**：Application 定义端口，Infrastructure 实现
6. **符合 Go 习惯**：包名简洁，文件适度聚合
7. **独立入口**：`cmd/` 下各服务独立入口，减少依赖耦合

### ⚠️ 注意事项

1. **Domain 层**不应该依赖任何框架（无 GORM tags）
2. **Application 层**不应该 import Infrastructure 的具体实现
3. **事务管理**应该在 Application 层声明边界
4. **错误处理**统一使用领域错误 + HTTP 状态码映射
5. **测试**每一层都应该有对应的测试

---

## 扩展阅读

- [Domain-Driven Design（Eric Evans）](https://www.domainlanguage.com/ddd/)
- [Implementing Domain-Driven Design（Vaughn Vernon）](https://vaughnvernon.com/)
- [Go DDD Example](https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example)
- [Clean Architecture in Go](https://github.com/bxcodec/go-clean-arch)
