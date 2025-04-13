package auth

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/dmasior/service-go/internal/database/dbgen"
	"github.com/dmasior/service-go/internal/domain"
	"github.com/dmasior/service-go/internal/jsonresponder"
	"github.com/dmasior/service-go/internal/jwt"

	"github.com/jackc/pgx/v5"
)

type (
	Auth struct {
		jwt     *jwt.JWT
		queries *dbgen.Queries
	}
	userKey struct{}
)

var openRoutes = map[string]struct{}{
	"/v1/signin": {},
	"/v1/signup": {},
}

func New(queries *dbgen.Queries, jwt *jwt.JWT) *Auth {
	return &Auth{queries: queries, jwt: jwt}
}

func (a *Auth) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Check if the request is for an open route
		if _, ok := openRoutes[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}

		// Require authentication for all other routes
		token, ok := getBearerToken(r)
		if !ok {
			jsonresponder.Unauthorized(w)
			return
		}

		email, err := a.jwt.SubjectFromToken(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userFromDB, err := a.queries.GetUserByEmail(ctx, email)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				jsonresponder.Unauthorized(w)
				return
			}

			slog.ErrorContext(r.Context(), "get user by email", slog.Any("err", err), slog.String("email", email))
			jsonresponder.InternalServerError(w)
			return
		}

		user := domain.User{
			ID:          userFromDB.ID,
			Email:       userFromDB.Email,
			LastLoginAt: userFromDB.LastLoginAt,
			CreatedAt:   userFromDB.CreatedAt,
		}

		ctx = NewContext(ctx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Middleware(queries *dbgen.Queries, jwt *jwt.JWT) func(next http.Handler) http.Handler {
	a := New(queries, jwt)
	return a.Handler
}

func getBearerToken(r *http.Request) (string, bool) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", false
	}

	return parseBearerTokenAuth(header)
}

func parseBearerTokenAuth(header string) (string, bool) {
	const prefix = "Bearer "
	if len(header) < len(prefix) || header[:len(prefix)] != prefix {
		return "", false
	}

	return header[len(prefix):], true
}
