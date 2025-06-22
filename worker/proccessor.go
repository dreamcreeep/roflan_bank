package worker

import (
	"context"
	"log/slog"

	db "github.com/dreamcreeep/roflan_bank/db/sqlc"
	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	logger *slog.Logger
	store  db.Store
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, logger *slog.Logger, store db.Store) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				logger.Error("process task failed",
					slog.Any("error", err),
					slog.String("type", task.Type()),
					slog.String("payload", string(task.Payload())))
			}),
			// Используем наш кастомный логгер
			Logger: NewSlogAsynqLogger(logger),
		},
	)

	return &RedisTaskProcessor{
		server: server,
		store:  store,
		logger: logger,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
}
