package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Email    EmailConfig
	Payment  PaymentConfig
	App      AppConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string
	Port int
	Env  string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret               string
	AccessTokenExpiry    time.Duration
	RefreshTokenExpiry   time.Duration
}

// EmailConfig 邮件配置
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
}

// PaymentConfig 支付配置
type PaymentConfig struct {
	StripeSecretKey      string
	StripePublishableKey string
}

// AppConfig 应用配置
type AppConfig struct {
	Name     string
	Env      string
	Debug    bool
	LogLevel string
}

// Load 加载配置
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// 环境变量支持
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config

	// Server
	cfg.Server.Host = viper.GetString("server.host")
	cfg.Server.Port = viper.GetInt("server.port")
	cfg.Server.Env = viper.GetString("server.env")

	// Database
	cfg.Database.Host = viper.GetString("database.host")
	cfg.Database.Port = viper.GetInt("database.port")
	cfg.Database.User = viper.GetString("database.user")
	cfg.Database.Password = viper.GetString("database.password")
	cfg.Database.DBName = viper.GetString("database.dbname")
	cfg.Database.SSLMode = viper.GetString("database.sslmode")

	// Redis
	cfg.Redis.Host = viper.GetString("redis.host")
	cfg.Redis.Port = viper.GetInt("redis.port")
	cfg.Redis.Password = viper.GetString("redis.password")
	cfg.Redis.DB = viper.GetInt("redis.db")

	// JWT
	cfg.JWT.Secret = viper.GetString("jwt.secret")
	cfg.JWT.AccessTokenExpiry = viper.GetDuration("jwt.access_token_expiry")
	cfg.JWT.RefreshTokenExpiry = viper.GetDuration("jwt.refresh_token_expiry")

	// Email
	cfg.Email.SMTPHost = viper.GetString("email.smtp_host")
	cfg.Email.SMTPPort = viper.GetInt("email.smtp_port")
	cfg.Email.SMTPUsername = viper.GetString("email.smtp_username")
	cfg.Email.SMTPPassword = viper.GetString("email.smtp_password")
	cfg.Email.SMTPFrom = viper.GetString("email.smtp_from")

	// Payment
	cfg.Payment.StripeSecretKey = viper.GetString("payment.stripe_secret_key")
	cfg.Payment.StripePublishableKey = viper.GetString("payment.stripe_publishable_key")

	// App
	cfg.App.Name = viper.GetString("app.name")
	cfg.App.Env = viper.GetString("app.env")
	cfg.App.Debug = viper.GetBool("app.debug")
	cfg.App.LogLevel = viper.GetString("app.log_level")

	return &cfg, nil
}

