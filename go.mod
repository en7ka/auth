module github.com/en7ka/auth

go 1.24

toolchain go1.24.3

require (
	github.com/brianvoe/gofakeit v3.18.0+incompatible
	google.golang.org/grpc v1.73.0 // рантайм для gRPC-кода
	google.golang.org/protobuf v1.36.6 // рантайм для protobuf-сообщений
	github.com/envoyproxy/protoc-gen-validate v1.2.1
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.1
	google.golang.org/grpc v1.73.0
	google.golang.org/protobuf v1.36.6
)

require (
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250603155806-513f23925822 // indirect
)
