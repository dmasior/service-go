# https://docs.sqlc.dev/en/stable/overview/install.html
.PHONY: generate
generate:
	@go generate ./...
	@sqlc generate -f sqlc.yaml

.PHONY: start-deps
start-deps:
	@docker compose up -d postgres

.PHONY: stop-deps
stop-deps:
	@docker compose down

.PHONY: run-api
run-api:
	@go run ./cmd/api

.PHONY: run-worker
run-worker:
	@go run ./cmd/worker

.PHONY: lint
lint:
	@golangci-lint run --config .golangci.yaml

.PHONY: lint-fix
lint-fix:
	@golangci-lint run --fix --config .golangci.yaml

.PHONY: test
test:
	@go test -v ./...
