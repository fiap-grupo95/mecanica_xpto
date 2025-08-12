APP_SERVICE_NAME=app
APP_CONTAINER_NAME=mecanica_xpto
DB_SERVICE_NAME=db
DB_CONTAINER_NAME=db
APP_BINARY_PATH=/app/mecanica-xpto-api
GO_BIN=$(shell go env GOPATH)/bin

.PHONY: init up down logs swag-install swag-generate swag-run

init: swag-install swag-generate
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

# Instala swag globalmente via Go (requer Go instalado no host)
swag-install:
	@echo "Instalando swag (se já instalado, ignorar)..."
	go install github.com/swaggo/swag/cmd/swag@latest

# Gera documentação Swagger a partir dos comentários no código
swag-generate:
	@echo "Gerando documentação Swagger..."
	@if [ ! -f "$(GO_BIN)/swag" ]; then \
		echo "swag não encontrado em $(GO_BIN), por favor rode 'make swag-install' primeiro" && exit 1; \
	fi
	$(GO_BIN)/swag init -g internal/infrastructure/http/routes/routes.go --output ./docs --parseDependency --parseInternal

# Gera documentação Swagger e executa o servidor localmente (fora do Docker)
swag-run: swag-generate
	@echo "Rodando a aplicação localmente..."
	go run cmd/api/main.go
