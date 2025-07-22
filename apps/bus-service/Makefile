# Load environment variables from .env
include .env
export

.DEFAULT_GOAL := help

migrationPath=./migrations
dbConnection=$(POSTGRES_URL)

.PHONY: help migrate-up migrate-down migrate-create seed truncate scraper serve swag

help:
	@echo ""
	@echo "📦 Mbus Makefile Commands"
	@echo "-------------------------"
	@echo "make migrate-up                  Run all up migrations"
	@echo "make migrate-down                Revert the last migration"
	@echo "make migrate-create name=NAME    Create a new migration with given name"
	@echo "make seed                        Seed the database with initial data"
	@echo "make truncate                    Truncate all database tables"
	@echo "make scraper                     Run the Marprom scraper"
	@echo "make serve                       Run the Go backend server"
	@echo "make swag                        Generate Swagger documentation"
	@echo ""

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
	@echo "Seeding database..."
	@go run ./cmd/seed/main.go

truncate:
	@echo "Truncating tables..."
	@go run ./cmd/truncate-tables/main.go

scraper:
	@echo "Running scraper..."
	@go run ./cmd/scraper/main.go

serve:
	@echo "Starting server..."
	@go run ./cmd/server/main.go

swag:
	@echo "Generating Swagger docs..."
	@swag init --parseDependency -g cmd/server/main.go
