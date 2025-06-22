package worker

import (
	"fmt"
	"log/slog"
	"os"
)

// SlogAsynqLogger - это адаптер для использования slog в качестве логгера asynq.
type SlogAsynqLogger struct {
	logger *slog.Logger
}

func NewSlogAsynqLogger(logger *slog.Logger) *SlogAsynqLogger {
	return &SlogAsynqLogger{logger: logger}
}

func (s *SlogAsynqLogger) Debug(args ...interface{}) {
	s.logger.Debug(fmt.Sprint(args...))
}

func (s *SlogAsynqLogger) Info(args ...interface{}) {
	s.logger.Info(fmt.Sprint(args...))
}

func (s *SlogAsynqLogger) Warn(args ...interface{}) {
	s.logger.Warn(fmt.Sprint(args...))
}

func (s *SlogAsynqLogger) Error(args ...interface{}) {
	s.logger.Error(fmt.Sprint(args...))
}

func (s *SlogAsynqLogger) Fatal(args ...interface{}) {
	s.logger.Error(fmt.Sprint(args...))
	os.Exit(1)
}
