COMPOSE=docker compose

.PHONY: up down restart build logs ps test health clean

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
