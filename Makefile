.DEFAULT_GOAL := help

# ==============================================================================
# Docker-окружение
# ==============================================================================

up: ## Запустить все сервисы в Docker
	docker-compose up -d

down: ## Остановить все сервисы и удалить тома
	docker-compose down -v

build: ## Пересобрать Docker образы
	docker-compose build --no-cache

logs: ## Показать логи всех сервисов
	docker-compose logs -f

server: ## Запустить сервер
	go run main.go

# ==============================================================================
# Миграции базы данных (выполняются в Docker)
# ==============================================================================

migrateup: ## Применить все доступные миграции
	docker-compose run --rm migrate up

migratedown: ## Откатить последнюю примененную миграцию
	docker-compose run --rm migrate down

# ==============================================================================
# Утилиты для разработки
# ==============================================================================

sqlc: ## Сгенерировать Go код из SQL запросов
	sqlc generate

test: ## Запустить Go тесты
	go test -v -cover ./...

# ==============================================================================
# Справка
# ==============================================================================
proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	proto/*.proto

evans:
	evans --host localhost --port 9090 -r repl

.PHONY: help up down build logs migrateup migratedown sqlc test proto server evans

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'