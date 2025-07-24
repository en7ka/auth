module github.com/en7ka/auth

go 1.23.0

toolchain go1.24.3

require (
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // рантайм валидатора
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.1 // HTTP-gateway + OpenAPI
	google.golang.org/grpc v1.73.0 // рантайм для gRPC-кода
	google.golang.org/protobuf v1.36.6 // рантайм для protobuf-сообщений
)
