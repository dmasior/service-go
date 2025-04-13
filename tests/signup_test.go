package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dmasior/service-go/internal/apiserver"
	"github.com/dmasior/service-go/internal/apiserver/apigen"
	"github.com/dmasior/service-go/internal/mailing"
	"github.com/dmasior/service-go/tests/testdb"
)

func TestSignUp(t *testing.T) {
	// Arrange
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	dbPool, teardown := testdb.NewPool()
	defer teardown()

	mailSpy := mailing.NewSpy()
	mux := apiserver.New(
		apiServerDefaultConfig,
		apiserver.WithDBPool(dbPool),
		apiserver.WithTurnstile(alwaysPassTurnstile),
		apiserver.WithJWT(jwtService),
		apiserver.WithMailer(mailSpy),
	).Mux()

	body := apigen.SignUpRequest{
		Email:    "test@example.com",
		Password: "password",
		Captcha:  "test-captcha",
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("could not marshal payload: %v", err)
	}

	req := httptest.NewRequestWithContext(ctx, "POST", "/v1/signup", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	// Act
	mux.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		t.Logf("response: %s", w.Body.String())
	}
}
