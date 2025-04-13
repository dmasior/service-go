package turnstile

type Config struct {
	SecretKey string `env:"TURNSTILE_SECRET_KEY,required"`
}
