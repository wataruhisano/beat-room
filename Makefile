-include .env
export

MIGRATE := go run -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@v4.20.1

.PHONY: help setup db-up db-down migrate api web check check-api check-web

help: ## Show available commands
	@grep -hE '^[a-z-]+:.*## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "} {printf "  %-12s %s\n", $$1, $$2}'

setup: ## Create env files and install web dependencies
	test -f .env || cp .env.example .env
	test -f apps/web/.env.local || cp apps/web/.env.example apps/web/.env.local
	cd apps/web && pnpm install

db-up: ## Start PostgreSQL
	docker compose up -d --wait

db-down: ## Stop PostgreSQL (data is kept)
	docker compose down

migrate: ## Apply database migrations
	$(MIGRATE) -path apps/api/migrations -database "$(DATABASE_URL)" up

api: ## Run the API on :8080
	cd apps/api && go run .

web: ## Run the web app on :3000
	cd apps/web && pnpm dev

check: check-api check-web  ## Run all checks

check-api: ## Format, vet, lint, and test the API
	cd apps/api && test -z "$$(gofmt -l .)"
	cd apps/api && go vet ./...
	cd apps/api && golangci-lint run ./...
	cd apps/api && go test ./...

check-web: ## Format, lint, typecheck, and test the web app
	cd apps/web && pnpm format:check
	cd apps/web && pnpm lint
	cd apps/web && pnpm typecheck
	cd apps/web && pnpm test
