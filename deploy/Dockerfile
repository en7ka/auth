docker file  FROM golang:1.24.5-alpine3.22 AS builder

COPY . /github.com/en7ka/auth/source
WORKDIR /github.com/en7ka/auth/source

RUN go mod download

RUN go build -o ./bin/user-server ./cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/en7ka/auth/source/bin/user-server .

CMD ["./user-server"]