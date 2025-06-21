package util

import (
	"context"
	"log/slog"
	"os"
	"time"
)

// NewLogger создает новый логгер slog на основе конфигурации.
func NewLogger(env, format, levelStr string) *slog.Logger {
	var level slog.Level
	switch levelStr {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var handler slog.Handler
	handlerOpts := &slog.HandlerOptions{
		Level: level,
	}

	// Для разработки используем красивый текстовый формат, для продакшена - JSON.
	if env == "development" && format == "text" {
		handler = slog.NewTextHandler(os.Stdout, handlerOpts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, handlerOpts)
	}

	logger := slog.New(handler)
	return logger
}

// LogRequest логирует HTTP запросы с использованием slog.
func LogRequest(logger *slog.Logger, method, path, clientIP, userAgent string, statusCode int, duration time.Duration) {
	logger.Info("HTTP Request",
		slog.String("method", method),
		slog.String("path", path),
		slog.String("client_ip", clientIP),
		slog.String("user_agent", userAgent),
		slog.Int("status_code", statusCode),
		slog.Duration("duration", duration),
	)
}

// LogError логирует ошибки с использованием slog.
func LogError(ctx context.Context, logger *slog.Logger, message string, err error) {
	logger.ErrorContext(ctx, message, slog.Any("error", err))
}

// LogInfo логирует информационные сообщения с использованием slog.
func LogInfo(ctx context.Context, logger *slog.Logger, message string, args ...any) {
	logger.InfoContext(ctx, message, args...)
}

// LogDebug логирует отладочные сообщения с использованием slog.
func LogDebug(ctx context.Context, logger *slog.Logger, message string, args ...any) {
	logger.DebugContext(ctx, message, args...)
}
