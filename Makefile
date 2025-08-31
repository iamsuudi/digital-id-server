include .env.mk

db_url = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

DB_EXEC = docker compose -f $(COMPOSE_FILE) exec -T db
SERVER_EXEC = docker compose -f $(COMPOSE_FILE) exec -T server

generate:
	$(SERVER_EXEC) sqlc generate

create_migration:
	@read -p "Enter migration name: " name; \
	$(SERVER_EXEC) migrate create -ext sql -dir database/migrations -seq $$name

migrate_up:
	$(SERVER_EXEC) migrate -path database/migrations -database "$(db_url)" up

migrate_down:
	$(SERVER_EXEC) migrate -path database/migrations -database "$(db_url)" down

migrate_force:
	$(SERVER_EXEC) migrate -path database/migrations -database "$(db_url)" force $(version)

migrate_version:
	$(SERVER_EXEC) migrate -path database/migrations -database "$(db_url)" version

schema_dump:
	$(DB_EXEC) pg_dump -U $(DB_USER) -d $(DB_NAME) --schema-only > database/schema.sql

create_db:
	$(DB_EXEC) psql -U $(DB_USER) -c "CREATE DATABASE $(DB_NAME);"

reset_db:
	$(DB_EXEC) psql -U $(DB_USER) -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

reload_schema: reset_db
	$(DB_EXEC) sh -c 'for f in /schemas/*.sql; do \
        psql -U $(DB_USER) \
            -d $(DB_NAME) \
            -v ON_ERROR_STOP=1 \
            -f "$$f"; \
    done'

clear_migrations:
	rm database/migrations/*.sql

view_queries:
	@echo "Current queries:"
	@find database/queries -type f -name "*.sql" -exec sh -c 'echo "\n--- {} ---"; cat {}' \;

list_query_names:
	@echo "Query names:"
	@find database/queries -type f -name "*.sql" -exec grep -E '^-- name:' {} \; | sed 's/-- name: //'

validate:
	@echo "Validating migration files..."
	@find database/migrations -type f -name "*.sql" -exec sqlc vet -f sqlc.yaml {} \;
	@echo "Validating query files..."
	@find database/queries -type f -name "*.sql" -exec sqlc vet -f sqlc.yaml {} \;

clear_db_data:
	$(DB_EXEC) psql -U $(DB_USER) -d $(DB_NAME) -f database/scripts/clear_data.sql

seed_db:
	$(SERVER_EXEC) go run ./cmd/seed/ -file=cmd/seed/residents.json

run_server:
	$(SERVER_EXEC) go run server/cmd/server/main.go
