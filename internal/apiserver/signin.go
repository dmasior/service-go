package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/mail"
	"time"

	"github.com/dmasior/service-go/internal/apiserver/apigen"
	"github.com/dmasior/service-go/internal/database/dbgen"
	"github.com/dmasior/service-go/internal/jsonresponder"

	"github.com/jackc/pgx/v5"
)

func (a *API) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Decode request
	req := apigen.SignInRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonresponder.Errors(w, http.StatusBadRequest, map[string][]string{
			"body": {"invalid JSON"},
		})
		return
	}

	// Validate request
	if errs := a.validateSignIn(ctx, &req); errs != nil {
		jsonresponder.Errors(w, http.StatusBadRequest, errs)
		return
	}

	queries := dbgen.New(a.dbPool)

	user, err := queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			jsonresponder.NotFound(w)
			return
		}
		slog.ErrorContext(ctx, "get user", slog.Any("err", err))
		jsonresponder.InternalServerError(w)
		return
	}

	if ok, err := a.hasher.VerifyPassword(req.Password, user.PasswordHash); err != nil {
		slog.ErrorContext(ctx, "verify password", slog.Any("err", err))
		jsonresponder.InternalServerError(w)
		return
	} else if !ok {
		jsonresponder.Errors(w, http.StatusUnauthorized, map[string][]string{
			"password": {"invalid password"},
		})
		return
	}

	expiresAt := time.Now().Add(time.Hour * 4)
	token, err := a.jwt.TokenForSubject(user.Email, expiresAt)
	if err != nil {
		slog.ErrorContext(ctx, "generate token", slog.Any("err", err))
		jsonresponder.InternalServerError(w)
		return
	}

	resp := apigen.SingInResponse{
		Data: apigen.TokenModel{
			Token: token,
		},
	}

	jsonresponder.JSON(w, http.StatusOK, resp)
}

// for now, it's the same as validateSignUp
func (a *API) validateSignIn(ctx context.Context, req *apigen.SignInRequest) map[string][]string {
	errs := jsonresponder.NewErrorContainer()

	if req.Email == "" {
		errs.Add("email", "email is required")
	}

	if req.Password == "" {
		errs.Add("password", "password is required")
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
