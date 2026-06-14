COMPOSE=docker compose

.PHONY: up down restart build logs ps test health clean migrate

up:
	$(COMPOSE) up -d --build

down:
	$(COMPOSE) down

restart: down up

build:
	$(COMPOSE) build

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps

test:
	go test ./...

health:
	curl -fsS http://localhost:$${BACKEND_PORT:-8080}/health
	curl -fsS http://localhost:$${BACKEND_PORT:-8080}/health/db

clean:
	$(COMPOSE) down --remove-orphans

migrate:
	docker exec -i transcendence-db psql -U $${POSTGRES_USER:-transcendence} -d $${POSTGRES_DB:-transcendence} < backend/sql/migrations/001_initial_schema.sql
