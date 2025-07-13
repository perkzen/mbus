# Load environment variables from .env
include .env
export

migrationPath=./migrations
dbConnection=$(POSTGRES_URL)

migrate-up:
	@echo "Running goose up..."
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING="$(dbConnection)" goose -dir=$(migrationPath) up

migrate-down:
	@echo "Running goose down..."
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING="$(dbConnection)" goose -dir=$(migrationPath) down

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "❌ Please provide a migration name, e.g. make migrate-create name=create_users_table"; \
		exit 1; \
	fi
	@GOOSE_DRIVER=postgres goose -dir=$(migrationPath) create $(name) sql

seed:
	@go run ./cmd/seed/main.go

scraper:
	@go run ./cmd/scraper/main.go

serve:
	@go run ./cmd/server/main.go

swag:
	@swag init --parseDependency -g cmd/server/main.go




