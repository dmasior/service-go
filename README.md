# Service-Go

This is an example Go service boilerplate. It's a web server and a worker processing background tasks. It's not based on a framework, but a few production-ready and easy replaceable libraries.

It's simple and enables fast development using the OpenAPI spec as the source of truth. Typical development flow:
1. Write the api spec
2. Run `make generate` to generate the server (and database) code
3. Implement the handlers

## Features

- [x] Basic structure
  - [x] Config values loaded from .env
  - [x] Logging setup [slog](https://pkg.go.dev/log/slog)
- [x] API HTTP server
  - [x] Router [chi](https://pkg.go.dev/github.com/go-chi/chi/v5), paths generated from [openapi v3 spec](./api.yaml) using [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)
  - [x] Server auth middleware
- [x] Background worker with simple task queue
- [x] Database setup
  - [x] Postgresql, but can be easily changed to MySQL or SQLite thanks to [sqlc](https://docs.sqlc.dev/en/stable/index.html)
  - [x] Migrations
- [x] Example user signup and signin flow
  - [x] Password hashing
  - [x] JWT-based auth
- [x] Example worker task creation and processing
  - [x] Task queue
  - [x] Task processing
- [x] Integration tests
  - [x] Signin, signup
  - [ ] Task creation and processing
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

After starting the dependencies, run the app CMDs using your IDE for easiest debugging experience.
You can also use the command line:

```sh
make run-api
```

```sh
make run-worker
```

Signup request example:
```sh
curl --location 'http://localhost:8080/v1/signup' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "user@example.com",
    "password": "password",
    "captcha": "test-captcha"
}'
```

After signing up, you can sign in:
```sh
curl --location 'http://localhost:8080/v1/signin' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "user@example.com",
    "password": "password",
    "captcha": "test-captcha"
}'
```

Open [openapi spec](./api.yaml) in https://editor.swagger.io to see the API docs.

## Running tests
You can run the tests using the command line:

```sh
make test
```
