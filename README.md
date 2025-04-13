# Service-Go

This is a Go web service starter-kit consisting of an HTTP server and a worker processing background tasks.
It's not based on a framework but a few easy replaceable libraries.

Service-Go promotes an API-first development approach using the OpenAPI specification as the source of truth. Development workflow:
1. Update the API spec
2. Run `make generate` to generate server stubs
3. Implement business logic in the generated handlers, and write tests

## Features

- [x] Basic structure
  - [x] Config values loaded from .env
  - [x] Logging setup [slog](https://pkg.go.dev/log/slog)
- [x] API HTTP server
  - [x] Router [chi](https://pkg.go.dev/github.com/go-chi/chi/v5), paths generated from [openapi v3 spec](./api.yaml) using [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)
  - [x] Server auth middleware
- [x] Background worker with simple task queue
- [x] Database setup
  - [x] PostgreSQL with [sqlc](https://docs.sqlc.dev/en/stable/index.html) for type-safe SQL. Can be replaced with SQLite, MySQL
  - [x] Migrations
- [x] Example user signup and signin flow
  - [x] Password hashing, JWT-based auth
- [x] Worker task creation via API and background processing
- [x] Integration tests
  - [x] Signin, signup, task creation
  - [ ] Task processing
  - [x] Separate db per test
  - [x] Fixtures
- [ ] Metrics and observability
- [ ] CI/CD pipeline

## Getting started

After cloning the repository, copy `.env.example` to `.env`

```sh
cp .env.example .env
```

Start the dependencies

```sh
make start-deps
```

Run the CMDs locally using your IDE for the best development and debug experience. But, you can also use:

```sh
make run-api
```

```sh
make run-worker
```

### Request examples

Signup:
```sh
curl --location 'http://localhost:8080/v1/signup' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "user@example.com",
    "password": "password",
    "captcha": "test-captcha"
}'
```

SignIn:
```sh
curl --location 'http://localhost:8080/v1/signin' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "user@example.com",
    "password": "password",
    "captcha": "test-captcha"
}'
```

Open [api spec](./api.yaml) in https://editor.swagger.io to see the API docs.

## Running tests

```sh
make test
```
