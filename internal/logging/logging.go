package logging

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
)

func SetupFromEnv() {
	addSource, err := strconv.ParseBool(os.Getenv("LOG_ADD_SOURCE"))
	if err != nil {
		panic("cannot parse LOG_ADD_SOURCE env var")
	}

	cfg := &Config{
		Level:     os.Getenv("LOG_LEVEL"),
		Format:    os.Getenv("LOG_FORMAT"),
		AddSource: addSource,
	}

	logger := newFromConfig(cfg)
	slog.SetDefault(logger)
}

func newFromConfig(cfg *Config) *slog.Logger {
	var level slog.Level
	level.UnmarshalText([]byte(cfg.Level)) //nolint:errcheck

	var handler slog.Handler

	logFormatLower := strings.ToLower(cfg.Format)

	if logFormatLower == "json" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: cfg.AddSource,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: cfg.AddSource,
		})
	}

	logger := slog.New(handler)

	return logger
}
