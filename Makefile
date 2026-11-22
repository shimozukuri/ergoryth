include .env
export

export MSYS_NO_PATHCONV := 1

COMPOSE := docker compose -f deploy/docker-compose.yaml

MIGRATE_SERVICE := ergoryth-postgres-migrate
MIGRATIONS_DIR := /migrations
DATABASE_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@ergoryth-postgres:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

.PHONY: proto-lint proto backend-run integration-tests-run unit-tests-run tests-run env-up env-down env-cleanup env-port-forwarder env-port-close migrate-create migrate-up migrate-down migrate-version

proto-lint:
	@buf lint

proto:
	@$(MAKE) proto-lint
	@buf generate

backend-run:
	@cd backend && \
	go run ./cmd/ergoryth

integration-tests-run:
	@$(MAKE) env-port-forwarder
	@trap '$(MAKE) -C "$(CURDIR)" env-port-close' EXIT; \
	cd ./backend/internal/infrastructure/persistence/postgres/pool/pgx && \
	go test ./...

unit-tests-run:
	@cd ./backend/internal/application/service/ && \
	go test ./...

tests-run:
	@$(MAKE) integration-tests-run
	@$(MAKE) unit-tests-run

env-up:
	@$(COMPOSE) up -d ergoryth-postgres

env-down:
	@$(COMPOSE) down ergoryth-postgres

env-cleanup:
	@read -p "Clear all volume environment files? Risk of data loss. [y/N]: " ans; \
	if [ "$$ans" = "y" ] || [ "$$ans" = "Y" ]; then \
	  $(COMPOSE) down -v ergoryth-postgres port-forwarder && \
	  echo "Environment volumes removed successfully"; \
	else \
	  echo "Cleanup is cancel"; \
	fi

env-port-forwarder:
	@$(COMPOSE) up -d port-forwarder

env-port-close:
	@$(COMPOSE) down port-forwarder

migrate-create:
	@test -n "$(name)" || (echo "Usage: make migrate-create name=create_leads" && exit 1)
	@$(COMPOSE) run --rm --no-deps $(MIGRATE_SERVICE) \
		create -ext sql -dir $(MIGRATIONS_DIR) -seq -digits 6 "$(name)"

migrate-up:
	@$(COMPOSE) run --rm $(MIGRATE_SERVICE) \
		-path $(MIGRATIONS_DIR) \
		-database "$(DATABASE_URL)" \
		up

migrate-down:
	@$(COMPOSE) run --rm $(MIGRATE_SERVICE) \
		-path $(MIGRATIONS_DIR) \
		-database "$(DATABASE_URL)" \
		down 1

migrate-version:
	@$(COMPOSE) run --rm $(MIGRATE_SERVICE) \
		-path $(MIGRATIONS_DIR) \
		-database "$(DATABASE_URL)" \
		version

