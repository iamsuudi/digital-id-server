include .env.mk

db_url=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

generate:
	sqlc generate

create_migration:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir db/migrations -seq $$name

migrate_up:
	migrate -path db/migrations -database "$(db_url)" up

migrate_down:
	migrate -path db/migrations -database "$(db_url)" down

migrate_force:
	migrate -path db/migrations -database "$(db_url)" force $(version)

migrate_version:
	migrate -path db/migrations -database "$(db_url)" version

schema_dump:
	pg_dump -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -d $(DB_NAME) --schema-only > db/schema.sql

reset_db:
	psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -c "DROP DATABASE IF EXISTS $(DB_NAME);"
	psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -c "CREATE DATABASE $(DB_NAME);"

clear_migrations:
	rm db/migrations/*.sql


