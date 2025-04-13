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
	"github.com/dmasior/service-go/tests/testdb"
)

func TestSignIn(t *testing.T) {
	// Arrange
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	dbPool, teardown := testdb.NewPool()
	defer teardown()

	mux := apiserver.New(
		apiServerDefaultConfig,
		apiserver.WithDBPool(dbPool),
		apiserver.WithTurnstile(alwaysPassTurnstile),
		apiserver.WithJWT(jwtService),
	).Mux()

	body := apigen.SignInRequest{
		Email:    testdb.UserJohn.Email,
		Password: testdb.UserJohn.Password,
		Captcha:  "test-captcha",
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("could not marshal payload: %v", err)
	}

	req := httptest.NewRequestWithContext(ctx, "POST", "/v1/signin", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	// Act
	mux.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		t.Logf("response: %s", w.Body.String())
	}

	var response apigen.SingInResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if response.Data.Token == "" {
		t.Errorf("expected token, got empty")
	}
}
