-include ../.env
export

LOCAL_BIN=$(CURDIR)/bin
GOOSE=$(LOCAL_BIN)/goose
MIGR_DIR=../db/migrations

PG_HOST_LOCAL=localhost
PG_PORT_LOCAL=5433

install-deps:
	@mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

local-migration-status:
	$(GOOSE) -dir $(MIGR_DIR) postgres "host=$(PG_HOST_LOCAL) port=$(PG_PORT_LOCAL) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable" status -v

local-migration-up:
	$(GOOSE) -dir $(MIGR_DIR) postgres "host=$(PG_HOST_LOCAL) port=$(PG_PORT_LOCAL) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable" up -v

local-migration-down:
	$(GOOSE) -dir $(MIGR_DIR) postgres "host=$(PG_HOST_LOCAL) port=$(PG_PORT_LOCAL) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable" down -v

grpc-load-test:
	ghz \
		--proto api/user_v1/user.proto \
		--call api.user_v1.UserAPI/Get \
		--data '{"id": 1}' \
		--rps 100 \
		--total 3000 \
		--insecure \
		localhost:50051

grpc-error-load-test:
	ghz \
		--proto api/user_v1/user.proto \
		--call api.user_v1.UserAPI/Get \
		--data '{"id": 0}' \
		--rps 100 \
		--total 3000 \
		--insecure \
		localhost:50051