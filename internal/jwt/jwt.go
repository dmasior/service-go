package jwt

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	secret []byte
}

func New(secret []byte) *JWT {
	return &JWT{secret: secret}
}

func (t *JWT) TokenForSubject(subject string, expiresAt time.Time) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		Subject:   subject,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := jwtToken.SignedString(t.secret)
	if err != nil {
		return "", fmt.Errorf("sign jwt token: %w", err)
	}

	return signedString, nil
}

func (t *JWT) SubjectFromToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return t.secret, nil
	})
	if err != nil {
		slog.Info("parse jwt token", slog.Any("error", err))
		return "", fmt.Errorf("parse jwt token: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid jwt token")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", fmt.Errorf("unexpected claims type")
	}

	return claims.Subject, nil
}
