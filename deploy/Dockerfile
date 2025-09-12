# Этап 1: Сборщик
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./cmd/grpc_server/ ./cmd/grpc_server/
COPY ./internal/ ./internal/
COPY ./pkg/ ./pkg/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/user-server ./cmd/grpc_server

FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /out/user-server /usr/local/bin/user-server

EXPOSE 8081 50051
ENTRYPOINT ["/usr/local/bin/user-server"]