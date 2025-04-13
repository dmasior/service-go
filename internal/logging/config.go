package logging

type Config struct {
	Level     string `env:"LOG_LEVEL" envDefault:"info"`
	Format    string `env:"LOG_FORMAT" envDefault:"text"`
	AddSource bool   `env:"LOG_ADD_SOURCE" envDefault:"false"`
}
