package worker

type Config struct {
	Count          int32 `env:"WORKER_COUNT" envDefault:"1"`
	TaskMaxAttempt int32 `env:"WORKER_TASK_MAX_ATTEMPT" envDefault:"2"`
}
