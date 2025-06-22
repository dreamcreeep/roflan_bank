package util

import (
	"time"
)

type Config struct {
	// Database Configuration
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`

	// Redis Configuration
	RedisAddress string `mapstructure:"REDIS_ADDRESS"`

	// Server Configuration
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`
	Environment       string `mapstructure:"ENVIRONMENT"`

	// JWT Configuration
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`

	// Logging Configuration
	LogLevel  string `mapstructure:"LOG_LEVEL"`
	LogFormat string `mapstructure:"LOG_FORMAT"`
}
