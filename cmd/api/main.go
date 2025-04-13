package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/dmasior/service-go/internal/apiserver"
	"github.com/dmasior/service-go/internal/database"
	"github.com/dmasior/service-go/internal/hashing"
	"github.com/dmasior/service-go/internal/jwt"
	"github.com/dmasior/service-go/internal/logging"
	"github.com/dmasior/service-go/internal/mailing"
	"github.com/dmasior/service-go/internal/turnstile"

	"github.com/caarlos0/env/v11"

	_ "github.com/joho/godotenv/autoload"
)

//go:generate go tool oapi-codegen -config=./../../oapi.yaml ../../api.yaml

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer done()

	logging.SetupFromEnv()

	if err := run(ctx); err != nil {
		done()
		slog.ErrorContext(ctx, err.Error())
		os.Exit(1)
	}

	slog.InfoContext(ctx, "bye")
}

func run(ctx context.Context) error {
	// Set GOMAXPROCS
	slog.InfoContext(ctx, "run", slog.Int("GOMAXPROCS", runtime.GOMAXPROCS(0)))

	// API config
	cfg := &struct {
		Server    apiserver.Config
		CORS      apiserver.CORSOptions
		Logging   logging.Config
		Database  database.Config
		Mailing   mailing.Config
		JWT       jwt.Config
		Turnstile turnstile.Config
	}{}
	if err := env.Parse(cfg); err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	// Setup database
	pool, err := database.NewPool(ctx, cfg.Database.DB, cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port)
	if err != nil {
		return fmt.Errorf("setup database: %w", err)
	}

	// Apply migrations
	if err = database.Migrate(pool); err != nil {
		return fmt.Errorf("database migrate: %w", err)
	}

	// Start server
	api := apiserver.New(
		cfg.Server,
		apiserver.WithDBPool(pool),
		apiserver.WithTurnstile(turnstile.NewService(cfg.Turnstile.SecretKey)),
		apiserver.WithMailer(mailing.SetupFromEnv()),
		apiserver.WithHasher(hashing.NewArgon2()),
		apiserver.WithJWT(jwt.New([]byte(cfg.JWT.SecretKey))),
		apiserver.WithCORSOptions(cfg.CORS),
		apiserver.WithSupportEmail(cfg.Mailing.SupportEmail),
	)

	return api.Start(ctx)
}
