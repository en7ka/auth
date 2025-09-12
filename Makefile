LOCAL_BIN := $(CURDIR)/bin
export PATH := $(LOCAL_BIN):$(PATH)
include .env
export
.PHONY: install-deps generate generate-auth

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate: generate-auth generate-user

generate-auth:
	mkdir -p pkg/auth_v1
	protoc -I ./api/proto/auth_v1 \
		--go_out=pkg/auth_v1 --go_opt=paths=source_relative \
		--go-grpc_out=pkg/auth_v1 --go-grpc_opt=paths=source_relative \
		--validate_out=lang=go:pkg/auth_v1 --validate_opt=paths=source_relative \
		./api/proto/auth_v1/auth.proto



generate-user:
	mkdir -p pkg/user_v1
	protoc -I ./api/proto/user_v1 \
		--go_out=pkg/user_v1 --go_opt=paths=source_relative \
		--go-grpc_out=pkg/user_v1 --go-grpc_opt=paths=source_relative \
		--validate_out=lang=go:pkg/user_v1 --validate_opt=paths=source_relative \
		./api/proto/user_v1/user.proto


build:
	GOOS=linux GOARCH=amd64 go build -o auth_linux cmd/main.go

copy-to-server:
	scp auth_linux root@87.228.114.237:

docker-build-and-push:
	docker buildx build -f deploy/Dockerfile --no-cache --platform linux/amd64 -t $(SELCLOUD_REGISTRY)/test-server:v0.0.1 .
	docker login -u $(SELCLOUD_USERNAME) -p $(SELCLOUD_TOKEN) $(SELCLOUD_REGISTRY)
	docker push $(SELCLOUD_REGISTRY)/test-server:v0.0.1

docker-up:
	docker compose -f ./deploy/docker-compose.yaml up --build -d

docker-down:
	docker compose -f ./deploy/docker-compose.yaml down

test:
	go clean -testcache
	go test ./... -covermode count -coverpkg=github.com/en7ka/auth/internal/service/auth/...,github.com/en7ka/auth/internal/api/auth/... -count 5


test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=github.com/en7ka/auth/internal/service/auth/...,github.com/en7ka/auth/internal/api/auth/... -count 5
	rm -rf coverage
	mkdir -p coverage
	grep -v 'mocks\|config' coverage.tmp.out > coverage/coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html;
	go tool cover -func=./coverage/coverage.out | grep "total";