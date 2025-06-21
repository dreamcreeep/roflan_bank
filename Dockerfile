# Стадия 1: Сборка бинарного файла
FROM golang:1.24-alpine AS builder

# Устанавливаем необходимые пакеты для сборки
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Копируем файлы зависимостей для лучшего кэширования
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Устанавливаем переменные для кросс-компиляции
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
ENV GOFLAGS="-ldflags=-w -ldflags=-s"

# Собираем бинарный файл с оптимизациями
RUN go build -a -installsuffix cgo -o main main.go

# Скачиваем migrate
RUN wget -O migrate.tar.gz https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz && \
    tar -xzf migrate.tar.gz && \
    chmod +x migrate

# Стадия 2: ВРЕМЕННО ДЛЯ ОТЛАДКИ (заменяем scratch на alpine)
FROM alpine:latest

WORKDIR /app

# Копируем приложение и утилиту для миграций
COPY --from=builder /app/main .
COPY --from=builder /app/migrate /usr/local/bin/migrate
COPY db/migration ./migration

EXPOSE 8080

CMD ["/app/main"]