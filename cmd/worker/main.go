package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/dmasior/service-go/internal/database"
	"github.com/dmasior/service-go/internal/logging"
	"github.com/dmasior/service-go/internal/worker"

	"github.com/caarlos0/env/v11"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer done()

	logging.SetupFromEnv()

	if err := run(ctx); err != nil {
		slog.ErrorContext(ctx, "exit on err", slog.Any("err", err))
		os.Exit(1)
	}

	slog.InfoContext(ctx, "bye")
}

func run(ctx context.Context) error {
	// GOMAXPROCS
	slog.InfoContext(ctx, "startup", slog.Int("GOMAXPROCS", runtime.GOMAXPROCS(0)))

	// Worker config
	cfg := &struct {
		Worker   worker.Config
		Database database.Config
	}{}
	if err := env.Parse(cfg); err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	// Setup database
	pool, err := database.NewPool(ctx, cfg.Database.DB, cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port)
	if err != nil {
		return fmt.Errorf("setup database: %w", err)
	}

	if err = database.Migrate(pool); err != nil {
		return fmt.Errorf("database migrate: %w", err)
	}

	// Start worker
	w := worker.New(
		cfg.Worker,
		worker.WithDBPool(pool),
	)

	wg := &sync.WaitGroup{}

	for i := range cfg.Worker.Count {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := fmt.Sprintf("W%d", i)

			if err = w.Run(ctx, id); err != nil {
				slog.ErrorContext(ctx, err.Error())
			}
		}()
	}

	wg.Wait()

	return nil
}
