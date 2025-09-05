FROM golang:1.24.5-alpine3.22 AS builder
WORKDIR /app

кеш модулей
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /out/user-server ./cmd/grpc_server

FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /out/user-server /usr/local/bin/user-server

EXPOSE 8081 50051
ENTRYPOINT ["/usr/local/bin/user-server"]
