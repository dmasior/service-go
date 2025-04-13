package apiserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/dmasior/service-go/internal/apiserver/apigen"
	"github.com/dmasior/service-go/internal/auth"
	"github.com/dmasior/service-go/internal/database/dbgen"
	"github.com/dmasior/service-go/internal/hashing"
	"github.com/dmasior/service-go/internal/jwt"
	"github.com/dmasior/service-go/internal/mailing"
	"github.com/dmasior/service-go/internal/turnstile"

	"github.com/go-chi/chi/v5"
	chimdl "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	API struct {
		config       Config
		dbPool       *pgxpool.Pool
		mailer       mailing.Mailer
		hasher       *hashing.Argon2
		turnstile    *turnstile.Service
		jwt          *jwt.JWT
		cors         CORSOptions
		supportEmail string
	}
	Option func(*API)
)

func New(cfg Config, opts ...Option) *API {
	api := &API{config: cfg}

	for _, opt := range opts {
		opt(api)
	}

	return api
}

func WithTurnstile(turnstile *turnstile.Service) Option {
	return func(api *API) {
		api.turnstile = turnstile
	}
}

func WithDBPool(pool *pgxpool.Pool) Option {
	return func(api *API) {
		api.dbPool = pool
	}
}

func WithMailer(mailer mailing.Mailer) Option {
	return func(api *API) {
		api.mailer = mailer
	}
}

func WithHasher(hasher *hashing.Argon2) Option {
	return func(api *API) {
		api.hasher = hasher
	}
}

func WithJWT(jwt *jwt.JWT) Option {
	return func(api *API) {
		api.jwt = jwt
	}
}

func WithCORSOptions(cors CORSOptions) Option {
	return func(api *API) {
		api.cors = cors
	}
}

func WithSupportEmail(email string) Option {
	return func(api *API) {
		api.supportEmail = email
	}
}

// Start the HTTP server on the provided port.
func (a *API) Start(ctx context.Context) error {
	// Set up the listener on the specified port
	addr := ":" + a.config.Port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listener on %s: %w", addr, err)
	}

	// Create a new HTTP server with the specified address and handler
	server := &http.Server{
		Addr:              listener.Addr().String(),
		Handler:           a.Mux(),
		ReadHeaderTimeout: a.config.ReadHeaderTimeout,
		ReadTimeout:       a.config.ReadTimeout,
		WriteTimeout:      a.config.WriteTimeout,
		IdleTimeout:       a.config.IdleTimeout,
	}

	return a.startHTTP(ctx, server, listener)
}

func (a *API) startHTTP(ctx context.Context, srv *http.Server, listener net.Listener) error {
	errCh := make(chan error, 1)
	doneCh := make(chan struct{})

	// Start the server in a goroutine
	go func() {
		defer close(doneCh)

		slog.InfoContext(ctx, "server is starting", "ip", listener.Addr().String())
		defer slog.InfoContext(ctx, "server is stopped")

		if err := srv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			select {
			case errCh <- err:
			default:
			}
		}
	}()

	// Wait for the server error or context cancellation
	select {
	case err := <-errCh:
		return fmt.Errorf("serve: %w", err)
	case <-ctx.Done():
		slog.DebugContext(ctx, "context is done")
	}

	// Shutdown the server gracefully
	shutdownCtx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()

	slog.DebugContext(ctx, "server is shutting down")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}
	close(errCh)

	// Wait for the server to finish
	<-doneCh

	return nil
}

func (a *API) Mux() *chi.Mux {
	mux := chi.NewMux()

	mux.Use(chimdl.Recoverer)
	mux.Use(chimdl.RealIP)
	mux.Use(chimdl.RequestID)
	mux.Use(cors.Handler(a.cors.ToChiOptions()))
	mux.Use(auth.Middleware(dbgen.New(a.dbPool), a.jwt))

	apigen.HandlerFromMux(a, mux)

	return mux
}
