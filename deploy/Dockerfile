# Этап 1: Сборщик
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

# Копируем ТОЛЬКО файлы модулей из корня проекта
COPY go.mod go.sum ./

# Загружаем зависимости (этот слой будет кешироваться почти всегда)
RUN go mod download

# Копируем ТОЛЬКО нужные для сборки директории
# Docker будет кешировать этот слой, пока вы не измените код в этих папках
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/user-server ./cmd/grpc_server

# Этап 2: Финальный образ
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /out/user-server /usr/local/bin/user-server

EXPOSE 8081 50051
ENTRYPOINT ["/usr/local/bin/user-server"]