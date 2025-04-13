package testdb

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/dmasior/service-go/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "embed"
)

//go:embed fixtures.sql
var testdata []byte

func NewPool() (*pgxpool.Pool, func()) {
	ctx := context.Background()

	dbName := "service-go"
	dbUser := "service-go"
	dbPassword := "service-go"

	schema, err := sql.EmbeddedFiles.ReadFile("migrations/1_schema.up.sql")
	if err != nil {
		panic(err)
	}

	schemaTmpPath := filepath.Join(os.TempDir(), "1_schema.sql")
	_ = os.WriteFile(schemaTmpPath, schema, 0o600)
	testDataTmpPath := filepath.Join(os.TempDir(), "testdata.sql")
	_ = os.WriteFile(testDataTmpPath, testdata, 0o600)

	postgresContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithInitScripts(
			schemaTmpPath,
			testDataTmpPath,
		),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		panic(err)
	}

	teardown := func() {
		if err = postgresContainer.Terminate(ctx); err != nil {
			slog.Info("terminate container", slog.Any("err", err))
		}
	}

	pool, err := pgxpool.New(ctx, postgresContainer.MustConnectionString(ctx))
	if err != nil {
		panic(err)
	}

	return pool, teardown
}
