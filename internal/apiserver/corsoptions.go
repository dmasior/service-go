package apiserver

import (
	"net/http"

	"github.com/go-chi/cors"
)

type CORSOptions struct {
	AllowedOrigins     []string `env:"CORS_ALLOWED_ORIGINS" envDefault:"*"`
	AllowOriginFunc    func(r *http.Request, origin string) bool
	AllowedMethods     []string `env:"CORS_ALLOWED_METHODS" envDefault:"GET,POST,PUT,DELETE,OPTIONS"`
	AllowedHeaders     []string `env:"CORS_ALLOWED_HEADERS" envDefault:"Accept,Authorization,Content-Type,X-CSRF-Token,X-Requested-With"`
	ExposedHeaders     []string `env:"CORS_EXPOSED_HEADERS" envDefault:""`
	AllowCredentials   bool     `env:"CORS_ALLOW_CREDENTIALS" envDefault:"false"`
	MaxAge             int      `env:"CORS_MAX_AGE" envDefault:"0"`
	OptionsPassthrough bool     `env:"CORS_OPTIONS_PASSTHROUGH" envDefault:"false"`
	Debug              bool     `env:"CORS_DEBUG" envDefault:"false"`
}

func (c *CORSOptions) ToChiOptions() cors.Options {
	return cors.Options{
		AllowedOrigins:   c.AllowedOrigins,
		AllowOriginFunc:  c.AllowOriginFunc,
		AllowedMethods:   c.AllowedMethods,
		AllowedHeaders:   c.AllowedHeaders,
		ExposedHeaders:   c.ExposedHeaders,
		AllowCredentials: c.AllowCredentials,
		MaxAge:           c.MaxAge,
		Debug:            c.Debug,
	}
}
