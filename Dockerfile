# Стадия 1: Сборка бинарного файла
FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

# Устанавливаем переменные для кросс-компиляции
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

RUN go build -o main main.go
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz

# Стадия 2: Минимальный контейнер для запуска
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrate /usr/local/bin/migrate
COPY app.env .
COPY db/migration ./migration

RUN apk add --no-cache curl
RUN chmod +x /usr/local/bin/migrate
EXPOSE 8080

CMD ["/app/main"]