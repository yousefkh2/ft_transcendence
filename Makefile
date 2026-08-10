COMPOSE=docker compose

.PHONY: up down restart build logs ps test health clean migrate

up:
	$(COMPOSE) up -d --build
	@echo "Waiting for DB to be ready..."
	@until docker exec transcendence-db pg_isready -U $${POSTGRES_USER:-transcendence} -d $${POSTGRES_DB:-transcendence} > /dev/null 2>&1; do sleep 1; done
	$(MAKE) migrate

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
	docker exec -i transcendence-db psql -v ON_ERROR_STOP=1 -U $${POSTGRES_USER:-transcendence} -d $${POSTGRES_DB:-transcendence} -c 'CREATE EXTENSION IF NOT EXISTS pgcrypto;'
	for f in backend/sql/migrations/*.sql; do \
		echo "Applying $$f"; \
		docker exec -i transcendence-db psql -v ON_ERROR_STOP=1 -1 -U $${POSTGRES_USER:-transcendence} -d $${POSTGRES_DB:-transcendence} < $$f; \
	done
