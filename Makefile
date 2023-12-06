ifeq ($(POSTGRES_SETUP_TEST),)
	POSTGRES_SETUP_TEST := user=postgres password=test dbname=test host=localhost port=5432 sslmode=disable
endif
export DB_NAME=test
export DB_PASSWORD=test
export DB_HOST=localhost
export DB_USER=postgres
export DB_PORT=5432

POSTGRES_CONTAINER_NAME = postgres_test
DOCKER_COMPOSE_FILE = docker-compose.yaml
INTERNAL_PKG_PATH=$(CURDIR)/internal/repository
MIGRATION_FOLDER=$(INTERNAL_PKG_PATH)/migrations

.PHONY: start-test-environment
start-test-environment:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

.PHONY: stop-test-environment
stop-test-environment:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

.PHONY: integration-test
integration-test: start-test-environment
	go test -tags integration ./...

.PHONY: integration-test-cover
integration-test-cover:
	go test -tags integration ./... -cover

.PHONY: unit-test
unit-test:
	go test ./... -cover

.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql

.PHONY: test-migration-up
test-migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

.PHONY: test-migration-down
test-migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down

.PHONY: clean-test-data
clean-test-data:
	docker exec $(POSTGRES_CONTAINER_NAME) psql -U postgres -d test -c "TRUNCATE table_name RESTART IDENTITY;"

.PHONY: run-tests
run-tests: unit-test integration-test stop-test-environment

.PHONY: all
all: run-tests
