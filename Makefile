APP_SERVICE_NAME=app
APP_CONTAINER_NAME=mecanica_xpto
DB_SERVICE_NAME=db
DB_CONTAINER_NAME=db
APP_BINARY_PATH=/app/mecanica-xpto-api

init:
	docker compose up -d --build

	@echo "⏳ Aguardando banco ficar pronto..."
	docker compose exec $(DB_SERVICE_NAME) sh -c 'until pg_isready -U $$POSTGRES_USER -d $$POSTGRES_DB; do sleep 1; done'

	@echo "⏳ Aguardando container $(APP_CONTAINER_NAME) estar rodando..."
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
