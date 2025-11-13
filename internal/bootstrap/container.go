package bootstrap

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/adapters/http"
	authhandler "github.com/lwmacct/251112-go-ddd-skeleton/internal/adapters/http/handler/auth"
	orderhandler "github.com/lwmacct/251112-go-ddd-skeleton/internal/adapters/http/handler/order"
	rbachandler "github.com/lwmacct/251112-go-ddd-skeleton/internal/adapters/http/handler/rbac"
	userhandler "github.com/lwmacct/251112-go-ddd-skeleton/internal/adapters/http/handler/user"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/adapters/http/middleware"
	appauth "github.com/lwmacct/251112-go-ddd-skeleton/internal/application/auth"
	appmenu "github.com/lwmacct/251112-go-ddd-skeleton/internal/application/menu"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/application/order"
	approle "github.com/lwmacct/251112-go-ddd-skeleton/internal/application/role"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/application/user"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/config"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/auth"
	domainorder "github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/order"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/rbac"
	domainuser "github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/user"
	infraauth "github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/auth"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/cache"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/logger"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/payment"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/persistence"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/persistence/repository"
)

// Container 依赖注入容器
type Container struct {
	Config *config.Config
	Router *gin.Engine
	Logger *logger.ZapLogger
}

// NewContainer 创建依赖注入容器
func NewContainer(cfg *config.Config) (*Container, error) {
	// 1. 初始化基础设施

	// Logger
	log, err := logger.NewZapLogger(cfg.App.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}
	middleware.SetLogger(log.GetZapLogger())

	// Database
	db, err := persistence.NewPostgres(persistence.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Redis
	_, err = cache.NewRedis(cache.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	// 2. 初始化仓储
	userRepo := repository.NewUserRepository(db)
	tfRepo := repository.NewTwoFactorRepository(db)
	patRepo := repository.NewPATRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	shipmentRepo := repository.NewShipmentRepository(db)
	// RBAC仓储
	roleRepo := repository.NewRoleRepo(db)
	permissionRepo := repository.NewPermissionRepo(db)
	menuRepo := repository.NewMenuRepo(db)

	// 3. 初始化领域服务
	userDomainService := domainuser.NewService(userRepo)
	authDomainService := auth.NewService(tfRepo, patRepo, sessionRepo)
	orderDomainService := domainorder.NewService(orderRepo, paymentRepo, shipmentRepo)
	rbacDomainService := rbac.NewService(roleRepo, permissionRepo, menuRepo)

	// 4. 初始化基础设施服务（端口实现）
	passwordHasher := infraauth.NewPasswordHasher()
	jwtIssuer := infraauth.NewJWTIssuer(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpiry*time.Minute,
		cfg.JWT.RefreshTokenExpiry*time.Hour*24,
	)
	totpGenerator := infraauth.NewTOTPGenerator(cfg.App.Name)
	// emailSender := email.NewSMTPSender(email.Config{
	// 	Host:     cfg.Email.SMTPHost,
	// 	Port:     cfg.Email.SMTPPort,
	// 	Username: cfg.Email.SMTPUsername,
	// 	Password: cfg.Email.SMTPPassword,
	// 	From:     cfg.Email.SMTPFrom,
	// })
	paymentGateway := payment.NewStripeGateway(cfg.Payment.StripeSecretKey)

	// 设置中间件依赖
	middleware.SetTokenValidator(jwtIssuer)
	middleware.SetRoleChecker(middleware.NewRBACRoleChecker(rbacDomainService))

	// 5. 初始化应用服务
	userService := user.NewService(userRepo, userDomainService, passwordHasher)
	authService := appauth.NewService(
		userRepo,
		tfRepo,
		patRepo,
		sessionRepo,
		authDomainService,
		jwtIssuer,
		passwordHasher,
		totpGenerator,
	)
	orderService := order.NewService(
		orderRepo,
		paymentRepo,
		shipmentRepo,
		orderDomainService,
		paymentGateway,
	)
	// RBAC应用服务
	menuService := appmenu.NewService(rbacDomainService, menuRepo)
	roleService := approle.NewService(roleRepo, permissionRepo, rbacDomainService)

	// 6. 初始化HTTP处理器
	userHandler := userhandler.NewHandler(userService)
	authHandler := authhandler.NewHandler(authService)
	orderHandler := orderhandler.NewHandler(orderService)
	menuHandler := rbachandler.NewMenuHandler(menuService)
	roleHandler := rbachandler.NewRoleHandler(roleService)

	// 7. 初始化路由
	router := http.SetupRouter(userHandler, authHandler, orderHandler, menuHandler, roleHandler)

	return &Container{
		Config: cfg,
		Router: router,
		Logger: log,
	}, nil
}
