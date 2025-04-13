FROM golang:1.24-alpine AS build

WORKDIR /

COPY . .

RUN go build -o /api ./cmd/api && \
    go build -o /worker ./cmd/worker

FROM alpine:3 AS prod

COPY --from=build /api /api
COPY --from=build /worker /worker
