package worker

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	Worker struct {
		config Config
		dbPool *pgxpool.Pool
	}
	Option func(*Worker)
)

func New(config Config, opts ...Option) Worker {
	w := Worker{
		config: config,
	}

	for _, opt := range opts {
		opt(&w)
	}

	return w
}

func WithDBPool(pool *pgxpool.Pool) Option {
	return func(w *Worker) {
		w.dbPool = pool
	}
}
