package apiserver

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"net/http"
	"time"

	"github.com/dmasior/service-go/internal/apiserver/apigen"
	"github.com/dmasior/service-go/internal/auth"
	"github.com/dmasior/service-go/internal/database/dbgen"
	"github.com/dmasior/service-go/internal/jsonresponder"
)

func (a *API) GetTask(w http.ResponseWriter, r *http.Request, ID string) {
	user := auth.MustFromContext(r.Context())
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	queries := dbgen.New(a.dbPool)
	task, err := queries.GetTask(ctx, ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			jsonresponder.NotFound(w)
			return
		}
		slog.ErrorContext(ctx, "get task", slog.Any("error", err))
		jsonresponder.InternalServerError(w)
		return
	}

	if user.ID != task.CreatedBy {
		jsonresponder.NotFound(w)
		return
	}

	resp := apigen.TaskResponse{
		Task: apigen.TaskModel{
			ID:              task.ID,
			Type:            task.Type,
			Payload:         task.Payload,
			Attempts:        task.Attempts,
			Status:          task.Status,
			StatusUpdatedAt: task.StatusUpdatedAt,
			CreatedAt:       task.CreatedAt,
			CreatedBy:       task.CreatedBy,
		},
	}

	jsonresponder.JSON(w, http.StatusOK, resp)
}
