package apiserver

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/dmasior/service-go/internal/apiserver/apigen"
	"github.com/dmasior/service-go/internal/auth"
	"github.com/dmasior/service-go/internal/database/dbgen"
	"github.com/dmasior/service-go/internal/domain"
	"github.com/dmasior/service-go/internal/idgen"
	"github.com/dmasior/service-go/internal/jsonresponder"
)

func (a *API) CreateTask(w http.ResponseWriter, r *http.Request) {
	user := auth.MustFromContext(r.Context())
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	// Decode request
	req := apigen.CreateTaskRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonresponder.Errors(w, http.StatusBadRequest, map[string][]string{
			"body": {"invalid JSON"},
		})
		return
	}

	// Validate request
	if errs := a.validateCreateTask(&req); errs != nil {
		jsonresponder.Errors(w, http.StatusBadRequest, errs)
		return
	}

	now := time.Now().UTC()
	queries := dbgen.New(a.dbPool)

	task, err := queries.CreateTask(ctx, dbgen.CreateTaskParams{
		ID:              idgen.NewTaskID(),
		CreatedBy:       user.ID,
		Type:            req.Type,
		Payload:         req.Payload,
		Status:          domain.TaskStatusCreated.String(),
		Attempts:        0,
		StatusUpdatedAt: now,
		CreatedAt:       now,
	})
	if err != nil {
		slog.ErrorContext(ctx, "create task", slog.String("type", req.Type), slog.String("payload", *req.Payload), slog.Any("error", err))
		jsonresponder.InternalServerError(w)
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

	jsonresponder.JSON(w, http.StatusCreated, resp)
}

// for now, it's the same as validateSignUp
func (a *API) validateCreateTask(req *apigen.CreateTaskRequest) map[string][]string {
	errs := jsonresponder.NewErrorContainer()

	typ := domain.TaskType(req.Type)
	if !typ.IsValid() {
		errs.Add("type", "invalid task type")
	}

	if req.Payload != nil && len(*req.Payload) == 0 {
		errs.Add("payload", "empty payload")
	}

	if req.Payload != nil && len(*req.Payload) > 1000 {
		errs.Add("payload", "payload too long")
	}

	// if there are any validation errors, return them
	if errs.NotEmpty() {
		return errs
	}

	return nil
}
