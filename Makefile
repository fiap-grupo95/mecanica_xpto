APP_SERVICE_NAME=app
APP_CONTAINER_NAME=mecanica_xpto
DB_SERVICE_NAME=db
DB_CONTAINER_NAME=db
APP_BINARY_PATH=/app/mecanica-xpto-api

.PHONY: init up down logs swag-generate-docker swag-run-docker test coverage coverage-html

init:
	cp .env-example .env
	docker compose up -d --build

	@echo "Aguardando banco ficar pronto..."
	docker compose exec $(DB_SERVICE_NAME) sh -c 'until pg_isready -U $$POSTGRES_USER -d $$POSTGRES_DB; do sleep 1; done'

	@echo "Aguardando container $(APP_CONTAINER_NAME) estar rodando..."
	@while [ -z "$$(docker compose ps -q $(APP_SERVICE_NAME))" ] || \
		[ "$$(docker inspect -f '{{.State.Running}}' $$(docker compose ps -q $(APP_SERVICE_NAME)))" != "true" ]; do \
		echo "Aguardando container..."; sleep 2; \
	done

	@echo "Container $(APP_CONTAINER_NAME) está rodando!"

	@echo "Aguardando app iniciar (sleep 10s dentro do container)..."
	docker compose exec $(APP_SERVICE_NAME) sh -c 'sleep 10'

	@echo "Rodando migrate..."
	docker compose exec $(APP_SERVICE_NAME) $(APP_BINARY_PATH) migrate

	@echo "Rodando seed..."
	docker compose exec $(APP_SERVICE_NAME) $(APP_BINARY_PATH) seed

up:
	docker compose up -d --build

down:
	docker compose down

logs:
	docker compose logs -f $(APP_SERVICE_NAME)

# Gera documentação Swagger dentro do container Go (opcional, caso queira rodar isolado)
swag-generate-docker:
	@echo "Gerando documentação Swagger dentro do container Go..."
	docker run --rm -v $(PWD):/app -w /app golang:1.24.4 \
		sh -c "go install github.com/swaggo/swag/cmd/swag@latest && \
		       /go/bin/swag init -g internal/infrastructure/http/routes/routes.go \
		       --output ./docs --parseDependency --parseInternal"

dev-up:
	cp .env-example .env
	docker compose up -d dev

test: dev-up
	docker compose exec dev go test ./... -v

coverage: dev-up
	docker compose exec dev go test ./... -coverprofile=coverage.out
	docker compose exec dev go tool cover -func=coverage.out

coverage-html: dev-up
	docker compose exec dev go test ./... -coverprofile=coverage.out
	docker cp $$(docker compose ps -q dev):/app/coverage.out .
	go tool cover -html=coverage.out
