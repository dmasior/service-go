package database

import (
	"errors"
	"fmt"
	"time"

	sgsql "github.com/dmasior/service-go/sql"

	"github.com/golang-migrate/migrate/v4"
	mpgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	// Register postgres driver
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

func Migrate(pool *pgxpool.Pool) error {
	sourceDrv, err := iofs.New(sgsql.EmbeddedFiles, "migrations")
	if err != nil {
		return fmt.Errorf("could not create migrations source: %w", err)
	}

	dbDrv, err := mpgx.WithInstance(stdlib.OpenDBFromPool(pool), &mpgx.Config{
		StatementTimeout: time.Second * 30,
	})
	if err != nil {
		return fmt.Errorf("could not create pgx instance: %w", err)
	}

	mr, err := migrate.NewWithInstance("iofs", sourceDrv, "pgx", dbDrv)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err = mr.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("migrate up: %w", err)
		}
	}

	return dbDrv.Close()
}
