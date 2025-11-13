package bootstrap

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-ddd-skeleton/internal/adapters/http"
	"github.com/yourusername/go-ddd-skeleton/internal/adapters/http/handler"
	"github.com/yourusername/go-ddd-skeleton/internal/adapters/http/middleware"
	appauth "github.com/yourusername/go-ddd-skeleton/internal/application/auth"
	"github.com/yourusername/go-ddd-skeleton/internal/application/order"
	"github.com/yourusername/go-ddd-skeleton/internal/application/user"
	"github.com/yourusername/go-ddd-skeleton/internal/config"
	"github.com/yourusername/go-ddd-skeleton/internal/domain/auth"
	domainorder "github.com/yourusername/go-ddd-skeleton/internal/domain/order"
	domainuser "github.com/yourusername/go-ddd-skeleton/internal/domain/user"
	infraauth "github.com/yourusername/go-ddd-skeleton/internal/infrastructure/auth"
	"github.com/yourusername/go-ddd-skeleton/internal/infrastructure/cache"
	"github.com/yourusername/go-ddd-skeleton/internal/infrastructure/logger"
	"github.com/yourusername/go-ddd-skeleton/internal/infrastructure/payment"
	"github.com/yourusername/go-ddd-skeleton/internal/infrastructure/persistence"
	"github.com/yourusername/go-ddd-skeleton/internal/infrastructure/persistence/repository"
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

	// 3. 初始化领域服务
	userDomainService := domainuser.NewService(userRepo)
	authDomainService := auth.NewService(tfRepo, patRepo, sessionRepo)
	orderDomainService := domainorder.NewService(orderRepo, paymentRepo, shipmentRepo)

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

	// 6. 初始化HTTP处理器
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)
	orderHandler := handler.NewOrderHandler(orderService)

	// 7. 初始化路由
	router := http.SetupRouter(userHandler, authHandler, orderHandler)

	return &Container{
		Config: cfg,
		Router: router,
		Logger: log,
	}, nil
}
