-include .env
export

LOCAL_BIN=$(CURDIR)/bin
GOOSE=$(LOCAL_BIN)/goose
MIGR_DIR?=$(MIGRATION_DIR)

install-deps:
	@mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

local-migration-status:
	$(GOOSE) -dir $(MIGR_DIR) postgres "host=$(PG_HOST) port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable" status -v

local-migration-up:
	$(GOOSE) -dir $(MIGR_DIR) postgres "host=$(PG_HOST) port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable" up -v

local-migration-down:
	$(GOOSE) -dir $(MIGR_DIR) postgres "host=$(PG_HOST) port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable" down -v
