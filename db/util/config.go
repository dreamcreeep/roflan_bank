package util

import (
	"time"
)

type Config struct {
	// Database Configuration
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`

	// Server Configuration
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`
	Environment       string `mapstructure:"ENVIRONMENT"`

	// JWT Configuration
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`

	// Redis Configuration
	RedisAddress  string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`

	// Logging Configuration
	LogLevel  string `mapstructure:"LOG_LEVEL"`
	LogFormat string `mapstructure:"LOG_FORMAT"`

	// Monitoring Configuration
	MetricsEnabled bool   `mapstructure:"METRICS_ENABLED"`
	MetricsPort    string `mapstructure:"METRICS_PORT"`

	// Security Configuration
	CORSAllowedOrigins string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	RateLimitRequests  int    `mapstructure:"RATE_LIMIT_REQUESTS"`
	RateLimitDuration  string `mapstructure:"RATE_LIMIT_DURATION"`

	// External Services
	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPort     int    `mapstructure:"SMTP_PORT"`
	SMTPUsername string `mapstructure:"SMTP_USERNAME"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`

	// Feature Flags
	EnableEmailVerification bool `mapstructure:"ENABLE_EMAIL_VERIFICATION"`
	EnableTwoFactorAuth     bool `mapstructure:"ENABLE_TWO_FACTOR_AUTH"`
	EnableOAuth             bool `mapstructure:"ENABLE_OAUTH"`
}
