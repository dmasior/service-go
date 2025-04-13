package worker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/dmasior/service-go/internal/database/dbgen"
	"github.com/dmasior/service-go/internal/domain"

	"github.com/jackc/pgx/v5"
)

func (w *Worker) Run(ctx context.Context, id string) error {
	slog.InfoContext(ctx, "worker run", slog.String("ID", id))

	queries := dbgen.New(w.dbPool)

	for {
		select {
		case <-ctx.Done():
			slog.InfoContext(ctx, "worker done", slog.String("ID", id))
			return nil
		default:
			now := time.Now()
			task, err := queries.PickTask(ctx, w.config.TaskMaxAttempt)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					time.Sleep(1 * time.Second)
					continue
				}
				return fmt.Errorf("pick task: %w", err)
			}

			slog.InfoContext(ctx, "processing task", slog.Any("taskID", task.ID))

			// Use a new context
			processCtx := context.Background()
			if err = w.process(processCtx, task); err != nil {
				slog.ErrorContext(ctx, "process task", slog.Any("err", err), slog.Any("taskID", task.ID))
				if err = queries.UpdateTaskStatus(ctx, dbgen.UpdateTaskStatusParams{
					ID:     task.ID,
					Status: domain.TaskStatusFailed.String(),
				}); err != nil {
					slog.ErrorContext(ctx, "update task status", slog.Any("err", err), slog.Any("taskID", task.ID))
				}
				continue
			}

			if err = queries.UpdateTaskStatus(ctx, dbgen.UpdateTaskStatusParams{
				ID:     task.ID,
				Status: domain.TaskStatusSuccess.String(),
			}); err != nil {
				slog.ErrorContext(ctx, "update task status", slog.Any("err", err), slog.Any("taskID", task.ID))
			}
			slog.InfoContext(ctx, "task processing finished", slog.Any("taskID", task.ID), slog.Any("duration", time.Since(now).String()))
		}
	}
}

func (w *Worker) process(_ context.Context, _ dbgen.Task) error {
	// task processing logic
	return nil
}
