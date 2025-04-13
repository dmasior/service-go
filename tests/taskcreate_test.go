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
	"github.com/dmasior/service-go/internal/domain"
	"github.com/dmasior/service-go/tests/testdb"
)

func TestTaskCreate(t *testing.T) {
	// Arrange
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	dbPool, teardown := testdb.NewPool()
	defer teardown()

	mux := apiserver.New(
		apiServerDefaultConfig,
		apiserver.WithDBPool(dbPool),
		apiserver.WithJWT(jwtService),
	).Mux()

	taskPayload := "{}"
	body := apigen.CreateTaskRequest{
		Type:    domain.TaskTypeFirst.String(),
		Payload: &taskPayload,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("could not marshal payload: %v", err)
	}

	jwtToken, err := jwtService.TokenForSubject(testdb.UserJohn.Email, time.Now().Add(1*time.Hour))
	if err != nil {
		t.Fatalf("could not create JWT token: %v", err)
	}

	req := httptest.NewRequestWithContext(ctx, "POST", "/v1/tasks", bytes.NewReader(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	w := httptest.NewRecorder()

	// Act
	mux.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		t.Logf("response: %s", w.Body.String())
	}

	resp := apigen.TaskResponse{}
	if err = json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if resp.Task.ID == "" {
		t.Errorf("expected task ID, got empty")
	}
	if resp.Task.Type != domain.TaskTypeFirst.String() {
		t.Errorf("expected task type %s, got %s", domain.TaskTypeFirst.String(), resp.Task.Type)
	}
	if resp.Task.Payload == nil {
		t.Errorf("expected task payload, got nil")
	}
	if *resp.Task.Payload != taskPayload {
		t.Errorf("expected task payload %s, got %s", taskPayload, *resp.Task.Payload)
	}
}

func TestTaskCreateRequireAuth(t *testing.T) {
	// Arrange
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	dbPool, teardown := testdb.NewPool()
	defer teardown()

	mux := apiserver.New(
		apiServerDefaultConfig,
		apiserver.WithDBPool(dbPool),
		apiserver.WithJWT(jwtService),
	).Mux()

	body := apigen.CreateTaskRequest{
		Type: domain.TaskTypeFirst.String(),
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("could not marshal payload: %v", err)
	}

	req := httptest.NewRequestWithContext(ctx, "POST", "/v1/tasks", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	// Act
	mux.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
		t.Logf("response: %s", w.Body.String())
	}
}
