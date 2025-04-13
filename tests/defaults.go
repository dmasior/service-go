package tests

import (
	"github.com/dmasior/service-go/internal/apiserver"
	"github.com/dmasior/service-go/internal/jwt"
	"github.com/dmasior/service-go/internal/turnstile"
)

var (
	apiServerDefaultConfig = apiserver.Config{
		Port: "8123",
	}

	alwaysPassTurnstile = turnstile.NewService("1x0000000000000000000000000000000AA")

	jwtService = jwt.New([]byte("test-secret"))
)
