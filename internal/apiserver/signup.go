package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/mail"
	"time"

	"github.com/dmasior/service-go/internal/apiserver/apigen"
	"github.com/dmasior/service-go/internal/database/dbgen"
	"github.com/dmasior/service-go/internal/idgen"
	"github.com/dmasior/service-go/internal/jsonresponder"

	"github.com/jackc/pgx/v5"
)

func (a *API) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Decode request
	req := apigen.SignUpRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonresponder.Errors(w, http.StatusBadRequest, map[string][]string{
			"body": {"invalid JSON"},
		})
		return
	}

	// Validate request
	if errs := a.validateSignUp(ctx, &req); errs != nil {
		jsonresponder.Errors(w, http.StatusBadRequest, errs)
		return
	}

	// Setup transaction
	tx, err := a.dbPool.Begin(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "begin tx", slog.Any("err", err))
	}
	defer tx.Rollback(ctx) //nolint:errcheck
	queries := dbgen.New(tx)

	// Check if user exists
	_, err = queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			slog.ErrorContext(ctx, "get user by email", slog.Any("err", err))
			jsonresponder.InternalServerError(w)
			return
		}
	} else {
		jsonresponder.Errors(w, http.StatusConflict, map[string][]string{
			"email": {"email already exists"},
		})
		return
	}

	// Create user
	now := time.Now()

	passwordHash, err := a.hasher.HashPassword(req.Password)
	if err != nil {
		slog.ErrorContext(ctx, "could not hash password", slog.Any("err", err))
		jsonresponder.InternalServerError(w)
		return
	}

	user, err := queries.CreateUser(ctx, dbgen.CreateUserParams{
		ID:           idgen.NewUserID(),
		Email:        req.Email,
		PasswordHash: passwordHash,
		LastLoginAt:  now,
		CreatedAt:    now,
	})
	if err != nil {
		slog.ErrorContext(ctx, "could not create user", slog.Any("err", err))
		jsonresponder.InternalServerError(w)
		return
	}

	if err = tx.Commit(ctx); err != nil {
		slog.ErrorContext(ctx, "could not commit transaction", slog.Any("err", err))
		jsonresponder.InternalServerError(w)
		return
	}

	// Send welcome email
	mailBody := fmt.Sprintf("Thank you for signing up, %s!\n\n If you have any questions or received an unexpected email, please contact us at %s.", user.Email, a.supportEmail)
	if err := a.mailer.Send(a.supportEmail, user.Email, "Welcome to service-go!", mailBody); err != nil {
		slog.ErrorContext(ctx, "could not send email", slog.Any("err", err))
	}

	// Respond created
	jsonresponder.Created(w)
}

func (a *API) validateSignUp(ctx context.Context, req *apigen.SignUpRequest) map[string][]string {
	errs := jsonresponder.NewErrorContainer()

	if req.Email == "" {
		errs.Add("email", "email is required")
	}

	if req.Password == "" {
		errs.Add("password", "password is required")
	}

	// it's checking bytes, not characters but it's good enough for now
	if len(req.Password) < 8 {
		errs.Add("password", "password must be at least 8 characters long")
	}

	if req.Captcha == "" {
		errs.Add("captcha", "captcha is required")
	}

	// validate email
	email, err := mail.ParseAddress(req.Email)
	if err != nil {
		errs.Add("email", "invalid email")
	} else {
		req.Email = email.Address
	}

	// if there are any validation errors, return them
	if errs.NotEmpty() {
		return errs
	}

	if !a.turnstile.CheckToken(ctx, req.Captcha) {
		errs.Add("captcha", "invalid captcha")
	}

	if errs.NotEmpty() {
		return errs
	}

	return nil
}
