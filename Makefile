LOCAL_BIN := $(CURDIR)/bin
include .env


install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.2.1
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.26.3
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.27.1
	GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@v0.1.7


generate: generate-auth generate-user swagger-static


generate-auth:
	mkdir -p pkg/auth_v1
	protoc -I ./api/proto/auth_v1 \
		--go_out=pkg/auth_v1 --go_opt=paths=source_relative \
		--go-grpc_out=pkg/auth_v1 --go-grpc_opt=paths=source_relative \
		--validate_out=lang=go:pkg/auth_v1 --validate_opt=paths=source_relative \
		./api/proto/auth_v1/auth.proto


generate-user:
	mkdir -p pkg/user_v1 pkg/swagger
	protoc -I ./api/proto/user_v1 \
		-I vendor.protogen \
		--go_out=pkg/user_v1 --go_opt=paths=source_relative \
		--go-grpc_out=pkg/user_v1 --go-grpc_opt=paths=source_relative \
		--validate_out=lang=go:pkg/user_v1 --validate_opt=paths=source_relative \
		--plugin=protoc-gen-validate=bin/protoc-gen-validate \
		--grpc-gateway_out=pkg/user_v1 --grpc-gateway_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway \
		--openapiv2_out=allow_merge=true,merge_file_name=api:pkg/swagger \
		--plugin=protoc-gen-openapiv2=bin/protoc-gen-openapiv2 \
		./api/proto/user_v1/user.proto


swagger-static:
	$(LOCAL_BIN)/statik -src=pkg/swagger/ -include='*.css,*.html,*.js,*.json,*.png'

build:
	GOOS=linux GOARCH=amd64 go build -o auth_linux cmd/main.go

copy-to-server:
	scp auth_linux root@87.228.114.237:

docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t $(SELCLOUD_REGISTRY)/test-server:v0.0.1 .
	docker login -u $(SELCLOUD_USERNAME) -p $(SELCLOUD_TOKEN) $(SELCLOUD_REGISTRY)
	docker push $(SELCLOUD_REGISTRY)/test-server:v0.0.1

docker-up:
	docker compose up -d

docker-down:
	docker compose down

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

vendor-proto:
		@if [ ! -d vendor.protogen/validate ]; then \
			mkdir -p vendor.protogen/validate &&\
			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
			mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
			rm -rf vendor.protogen/protoc-gen-validate ;\
		fi
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi
		@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
                	mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
                	git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
                	mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
                	rm -rf vendor.protogen/openapiv2 ;\
        fi