module github.com/en7ka/auth

go 1.24            // версия языка, без лишних нулей

toolchain go1.24.3 // можно оставить; если мешает ‒ просто удалите строку

require (
	// рантайм валидатора (нужен сгенерированному коду)
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect

	// HTTP-gateway + OpenAPI
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.1 // indirect

	// рантайм для gRPC-кода
	google.golang.org/grpc v1.73.0 // indirect

	// рантайм для protobuf-сообщений
	google.golang.org/protobuf v1.36.6 // indirect
)
