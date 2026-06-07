include .env
export

DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
LOCAL_DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@host.docker.internal:$(DB_PORT)/$(DB_NAME)?sslmode=disable

##@ Migrations

.PHONY: migrate-create
migrate-create: ## Create new migration (make migrate-create name=add_users)
	@if [ -z "$(name)" ]; then \
		echo "Error: name is required"; \
		echo "Usage: make migrate-create name=your_migration_name"; \
		exit 1; \
	fi
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: migrate-up
migrate-up: ## Apply all migrations (requires payment-service container running)
	docker-compose exec payment-service migrate -path migrations -database "$(DB_URL)" up

.PHONY: migrate-up-local
migrate-up-local: ## Apply all migrations (local dev — postgres in Docker)
	docker run --rm -v "$(CURDIR)/migrations:/migrations" migrate/migrate \
		-path=/migrations -database "$(LOCAL_DB_URL)" up

.PHONY: migrate-down
migrate-down: ## Rollback last migration (requires payment-service container running)
	docker-compose exec payment-service migrate -path migrations -database "$(DB_URL)" down 1

.PHONY: migrate-down-local
migrate-down-local: ## Rollback last migration (local dev)
	docker run --rm -v "$(CURDIR)/migrations:/migrations" migrate/migrate \
		-path=/migrations -database "$(LOCAL_DB_URL)" down 1

.PHONY: migrate-status
migrate-status: ## Show current migration version
	@docker-compose exec payment-service migrate -path migrations -database "$(DB_URL)" version 2>/dev/null || echo "No migrations applied"

.PHONY: migrate-force
migrate-force: ## Force version (make migrate-force version=5)
	@if [ -z "$(version)" ]; then \
		echo "Error: version is required"; \
		echo "Usage: make migrate-force version=NUMBER"; \
		exit 1; \
	fi
	docker-compose exec payment-service migrate -path migrations -database "$(DB_URL)" force $(version)

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'