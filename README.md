# 🏦 Rofan Bank API

Современное банковское API, построенное на Go с использованием Gin, PostgreSQL и Docker.

## 🚀 Быстрый старт

### Требования
- Docker & Docker Compose
- Go 1.24+ (для локальной разработки)
- Make (опционально)

### Запуск в продакшене
```bash
# Клонирование репозитория
git clone <repository-url>
cd roflan_bank

# Создание .env файла
cp env.example .env
# Отредактируйте .env файл с вашими настройками

# Запуск продакшн окружения
make prod-up

# Или без Make
docker-compose -f docker-compose.prod.yaml up -d
```

### Запуск для разработки
```bash
# Запуск только PostgreSQL
make dev-postgres

# Запуск полного окружения
make dev-up

# Запуск сервера локально
make server
```

## 📁 Структура проекта

```
roflan_bank/
├── api/                 # HTTP API handlers
├── db/                  # Database layer
│   ├── migration/       # SQL migrations
│   ├── sqlc/           # Generated SQL code
│   └── util/           # Database utilities
├── token/              # JWT/PASETO token management
├── Dockerfile          # Multi-stage Docker build
├── docker-compose.yaml # Development environment
├── docker-compose.prod.yaml # Production environment
└── Makefile           # Build and deployment commands
```

## 🔧 Конфигурация

### Переменные окружения

Создайте файл `.env` на основе `env.example` (если его нет). Этот файл используется для конфигурации как локального окружения, так и контейнеров.

Пример содержания `.env`:
```env
# Database
POSTGRES_USER=root
POSTGRES_PASSWORD=your_secure_password
POSTGRES_DB=simple_bank

# Server
SERVER_ADDRESS=0.0.0.0:8080
ENVIRONMENT=production

# JWT
TOKEN_SYMMETRIC_KEY=your_32_character_key
ACCESS_TOKEN_DURATION=15m
REFRESH_TOKEN_DURATION=24h

# Redis (для продакшена)
REDIS_ADDRESS=redis:6379
REDIS_PASSWORD=your_redis_password
```

## 🐳 Docker команды

### Основные команды
```bash
# Сборка образа
make build-image

# Запуск продакшн окружения
make prod-up

# Просмотр логов
make prod-logs

# Остановка
make prod-down

# Очистка
make clean
```

### Безопасность
```bash
# Сканирование уязвимостей
make security-scan
```

## 📊 API Endpoints

### Публичные endpoints
- `GET /health` - Health check
- `POST /users` - Регистрация пользователя
- `POST /users/login` - Авторизация

### Защищенные endpoints (требуют JWT токен)
- `POST /accounts` - Создание счета
- `GET /accounts/:id` - Получение счета
- `GET /accounts` - Список счетов
- `POST /transfers` - Перевод средств
- `POST /tokens/renew_access` - Обновление токена

## 🔒 Безопасность

### Реализованные меры безопасности:
- ✅ JWT/PASETO токены для аутентификации
- ✅ Хеширование паролей (bcrypt)
- ✅ Валидация входных данных
- ✅ CORS настройки
- ✅ Rate limiting (в продакшене)
- ✅ Непривилегированные контейнеры
- ✅ Минимальные Docker образы (scratch)

### Рекомендации для продакшена:
1. Используйте сильные пароли для всех сервисов
2. Настройте SSL/TLS сертификаты
3. Включите мониторинг и логирование
4. Настройте backup стратегию для PostgreSQL
5. Используйте secrets management для чувствительных данных

## 📈 Мониторинг

### Health Checks
- API: `GET /health`
- PostgreSQL: `pg_isready`
- Redis: `redis-cli ping`

### Метрики
```bash
# Статус контейнеров
make ps

# Статистика ресурсов
make stats
```

## 🧪 Тестирование

```bash
# Запуск всех тестов
make test

# Тесты с покрытием
make test-coverage

# Генерация моков
make mock
```

## 🔄 Миграции

```bash
# Применить все миграции
make migrateup

# Откатить все миграции
make migratedown

# Применить одну миграцию
make migrateup1

# Откатить одну миграцию
make migratedown1
```

## 🚀 Развертывание

### Kubernetes
```bash
# Развертывание в кластере
make k8s-deploy

# Удаление из кластера
make k8s-delete
```

### Docker Registry
```bash
# Сборка и публикация образа
make build-image
make push-image
```

## 📝 Логирование

Логи выводятся в stdout/stderr для совместимости с Docker и Kubernetes.

### Уровни логирования:
- `debug` - Детальная отладочная информация
- `info` - Общая информация (по умолчанию)
- `warn` - Предупреждения
- `error` - Ошибки

### Тестирование логирования

#### Полный автоматический тест
Эта команда автоматически запустит все сервисы, выполнит тесты и остановит сервисы после завершения.

```bash
# Запуск полного E2E теста
make test-e2e
```

#### Тестирование уже запущенного приложения
Если вы уже запустили сервисы с помощью `make dev-up`, вы можете выполнить только скрипт тестирования:

```bash
# Запуск тестов логирования
make test-logging
```

### Просмотр логов

```bash
# Просмотр логов в реальном времени
make logs-api

# Фильтрация логов по уровню
make logs-filtered

# Просмотр структурированных JSON логов
make logs-json
```

### Примеры логов

#### HTTP запросы:
```json
{
  "level": "INFO",
  "msg": "HTTP Request",
  "method": "POST",
  "path": "/users",
  "client_ip": "172.18.0.1",
  "user_agent": "curl/7.68.0",
  "status_code": 200,
  "duration_ms": 45.2
}
```

#### Создание пользователя:
```json
{
  "level": "INFO",
  "msg": "Creating new user",
  "username": "john_doe",
  "email": "john@example.com",
  "client_ip": "172.18.0.1"
}
```

#### Ошибки:
```json
{
  "level": "ERROR",
  "msg": "Failed to create user in database",
  "error": "pq: duplicate key value violates unique constraint"
}
```

## 🔧 Устранение неполадок

### Частые проблемы:

1. **PostgreSQL не запускается**
   ```bash
   # Проверьте логи
   docker-compose logs postgres
   
   # Пересоздайте volume
   docker-compose down -v
   docker-compose up -d
   ```

2. **API не может подключиться к БД**
   ```bash
   # Проверьте переменные окружения
   docker-compose exec api env | grep DB
   
   # Проверьте сеть
   docker network ls
   ```

3. **Миграции не применяются**
   ```bash
   # Запустите миграции вручную
   docker-compose exec migrate migrate -path /app/migration -database "postgresql://root:secret@postgres:5432/simple_bank?sslmode=disable" up
   ```

## 🤝 Вклад в проект

1. Fork репозитория
2. Создайте feature branch
3. Внесите изменения
4. Добавьте тесты
5. Создайте Pull Request

## 📄 Лицензия

MIT License

## 👥 Авторы

- Основной разработчик: [Ваше имя]
- Контакты: [email@example.com] 