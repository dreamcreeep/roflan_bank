package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/hibiken/asynq"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	logger *slog.Logger,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Error("failed to marshal task payload", slog.Any("error", err))
		return err
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)

	info, err := distributor.client.EnqueueContext(ctx, task)

	if err != nil {
		logger.Error("failed to enqueue task", slog.Any("error", err))
		return err
	}

	logger.Info("enqueued task", slog.String("type", task.Type()),
		slog.String("queue", info.Queue),
		slog.Int("max_retry", info.MaxRetry))

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	logger := processor.logger

	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		logger.Error("failed to unmarshal payload", slog.Any("error", err))
		return err
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("user not found", slog.String("username", payload.Username))
			return nil
		}
		logger.Error("failed to get user", slog.Any("error", err))
		return err
	}

	//TODO: send email to user

	logger.Info("processed task",
		slog.String("type", task.Type()),
		slog.String("username", user.Username),
		slog.String("email", user.Email),
	)

	return nil
}
